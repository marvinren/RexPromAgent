package db

import (
	"RexPromAgent/pkg/config"
	"RexPromAgent/pkg/log"
	"fmt"
	"testing"
)

func TestMysqlConnect(t *testing.T) {
	sqldb, err := ConnectDB()
	if err != nil {
		t.Fatalf("connection fail, %s", err)
	}
	t.Log(sqldb.db)

}

func TestReadAlertRuleConfig(t *testing.T) {
	sqldb, err := ConnectDB()
	if err != nil {
		t.Fatalf("connection fail, %s", err)
	}
	alerts := make([]AlertRule, 0)
	err = sqldb.FetchAlerts(&alerts)
	if err != nil {
		t.Errorf("read alert rules data error: %v", err)
	}
	fmt.Println(len(alerts))

}

func TestMain(m *testing.M) {
	log.Initialize()
	config.Initialize()
	m.Run()
}
