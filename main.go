package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/adifahmi/simplebank/api"
	db "github.com/adifahmi/simplebank/db/sqlc"
	"github.com/adifahmi/simplebank/util"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	dbPort := config.DBPort
	if os.Getenv("IS_GA") == "true" {
		dbPort = config.DBPortGa
	}
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?timeout=2s&parseTime=true",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		dbPort,
		config.DBName,
	)
	// root:root@tcp(127.0.0.1:3306)/app?timeout=2s&parseTime=true

	conn, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
