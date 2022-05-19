package db

import (
	"database/sql"
)

type AlertRule struct {
	RuleId      int64          `db:"rule_id"`
	AlertName   string         `db:"alert_name"`
	Expression  string         `db:"expression"`
	Duration    string         `db:"duration"`
	AlertLevel  string         `db:"alert_level"`
	AlertType   string         `db:"alert_type"`
	Notice      string         `db:"noitce"`
	State       string         `db:"state"`
	Description string         `db:"description"`
	CreateUID   sql.NullInt64  `db:"create_uid"`
	CreateTime  sql.NullTime   `db:"create_time"`
	UpdateUID   sql.NullInt64  `db:"update_uid"`
	UpdateTime  sql.NullTime   `db:"update_time"`
	TenantCode  sql.NullString `db:"tenant_code"`
	ProjectID   sql.NullInt64  `db:"project_id"`
	SystemID    sql.NullInt64  `db:"system_id"`
}
