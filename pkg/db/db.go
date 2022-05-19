package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"time"
)
import (
	_ "RexPromAgent/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLDB struct {
	db *sql.DB
}

func (s MySQLDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		logrus.Panicf("Failed to ping database: %s", err)
		return err
	}

	return nil
}

func (s MySQLDB) Close() {
	err := s.db.Close()
	if err != nil {
		logrus.Panicf("close database error: %v", err)
	}
}

func (s MySQLDB) SaveAlert(data *AlertGroup) error {
	return s.unitOfWork(func(tx *sql.Tx) error {

		r, err := tx.Exec(`
				INSERT INTO AlertGroup (time, receiver, status, externalURL, groupKey)
				VALUES (now(), ?, ?, ?, ?)`, data.Receiver, data.Status, data.ExternalURL, data.GroupKey)
		if err != nil {
			return fmt.Errorf("failed to insert into AlertGroups: %s", err)
		}

		alertGroupID, err := r.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get AlertGroups inserted id: %s", err)
		}
		fmt.Println(alertGroupID)

		for k, v := range data.GroupLabels {
			_, err := tx.Exec(`
				INSERT INTO GroupLabel (alertGroupID, GroupLabel, Value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into GroupLabel: %s", err)
			}
		}
		for k, v := range data.CommonLabels {
			_, err := tx.Exec(`
				INSERT INTO CommonLabel (alertGroupID, Label, Value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into CommonLabel: %s", err)
			}
		}
		for k, v := range data.CommonAnnotations {
			_, err := tx.Exec(`
				INSERT INTO CommonAnnotation (alertGroupID, Annotation, Value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into CommonAnnotation: %s", err)
			}
		}

		for _, alert := range data.Alerts {
			var result sql.Result
			if alert.EndsAt.Before(alert.StartsAt) {
				result, err = tx.Exec(`
				INSERT INTO Alert (alertGroupID, status, startsAt, generatorURL, fingerprint)
				VALUES (?, ?, ?, ?, ?)`,
					alertGroupID, alert.Status, alert.StartsAt, alert.GeneratorURL, alert.Fingerprint)
			} else {
				result, err = tx.Exec(`
				INSERT INTO Alert (alertGroupID, status, startsAt, endsAt, generatorURL, fingerprint)
				VALUES (?, ?, ?, ?, ?, ?)`,
					alertGroupID, alert.Status, alert.StartsAt, alert.EndsAt, alert.GeneratorURL, alert.Fingerprint)
			}
			if err != nil {
				return fmt.Errorf("failed to insert into Alert: %s", err)
			}

			alertID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get Alert inserted id: %s", err)
			}

			for k, v := range alert.Labels {
				_, err := tx.Exec(`
					INSERT INTO AlertLabel (AlertID, Label, Value)
					VALUES (?, ?, ?)`, alertID, k, v)
				if err != nil {
					return fmt.Errorf("failed to insert into AlertLabel: %s", err)
				}
			}
			for k, v := range alert.Annotations {
				_, err := tx.Exec(`
					INSERT INTO AlertAnnotation (AlertID, Annotation, Value)
					VALUES (?, ?, ?)`, alertID, k, v)
				if err != nil {
					return fmt.Errorf("failed to insert into AlertAnnotation: %s", err)
				}
			}
		}

		return nil
	})
}

func (s MySQLDB) unitOfWork(f func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err)
	}

	err = f(tx)

	if err != nil {
		log.Printf("commit data error: %s\n", err)
		e := tx.Rollback()
		if e != nil {
			return fmt.Errorf("failed to rollback transaction (%s) after failing execution: %s", e, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %s", err)
	}
	return nil
}

func ConnectDB() (*MySQLDB, error) {
	// Get the dsn/driver from config
	dsn := viper.GetString("database.dsn")
	driver := viper.GetString("database.driver")
	logrus.Infof("connect the %v database %v", driver, dsn)
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}

	// Set database connect settings.
	conn.SetMaxIdleConns(viper.GetInt("database.maxIdleConns"))
	conn.SetMaxOpenConns(viper.GetInt("database.maxOpenConns"))
	lifttime := viper.GetInt("database.connMaxLifetime")
	conn.SetConnMaxLifetime(time.Duration(int(time.Second) * lifttime))

	// Construct MySQLDB
	database := &MySQLDB{conn}

	// Ping the database connection
	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func (s MySQLDB) FetchAlerts(alerts *[]AlertRule) error {
	return s.unitOfWork(func(tx *sql.Tx) error {
		var alertRule AlertRule

		rows, err2 := tx.Query("select rule_id, alert_name, expression, duration, alert_level, alert_type, noitce, description, create_uid, state, create_time, update_uid, update_time, tenant_code, project_id, system_id from t_alert_rule where state='U' or state is null")

		if err2 != nil {
			return fmt.Errorf("failed to insert into AlertGroups: %s", err2)
		}
		for rows.Next() {
			alertRule = AlertRule{}
			err3 := rows.Scan(&alertRule.RuleId, &alertRule.AlertName, &alertRule.Expression, &alertRule.Duration,
				&alertRule.AlertLevel, &alertRule.AlertType, &alertRule.Notice, &alertRule.Description, &alertRule.CreateUID,
				&alertRule.State, &alertRule.CreateTime, &alertRule.UpdateUID, &alertRule.UpdateTime, &alertRule.TenantCode,
				&alertRule.ProjectID, &alertRule.SystemID)
			if err3 != nil {
				log.Printf("parse the database result error: %s", err3)
			}
			*alerts = append(*alerts, alertRule)
		}
		err4 := rows.Close()
		if err4 != nil {
			log.Printf("close the database result error: %s", err4)
		}
		return nil
	})

}
