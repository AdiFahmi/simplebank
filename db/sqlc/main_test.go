package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	dbSource = "app:app@tcp(127.0.0.1:3390)/app?timeout=2s&parseTime=true"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}

	testDB = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}
