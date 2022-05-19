package store

import (
	"RexPromAgent/db"
)

type Storer struct {
	Conn *db.MySQLDB
}
