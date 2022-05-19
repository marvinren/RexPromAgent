package db

import "testing"

func TestMysqlConnect(t *testing.T) {
	sqldb, err := ConnectDB()
	if err != nil {
		t.Fatalf("connection fail, %s", err)
	}
	t.Log(sqldb.db)

}
