package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/adifahmi/simplebank/util"
	_ "github.com/go-sql-driver/mysql"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}

	var version, isoLevel string

	err = conn.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.QueryRow("SELECT @@transaction_isolation;").Scan(&isoLevel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version, isoLevel)

	testDB = conn
	testQueries = New(conn)

	os.Exit(m.Run())
}
