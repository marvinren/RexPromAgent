package store

import (
	"RexPromAgent/pkg/db"
)

type Storer struct {
	Conn *db.MySQLDB
}
