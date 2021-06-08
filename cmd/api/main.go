package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/handler"
	"github.com/ezerw/wheel/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dBConn, err := db.Connect(config)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer dBConn.Close()

	store := db.NewStore(dBConn)
	server, err := handler.NewServer(config, store)
	if err != nil {
		log.Panic("cannot create server:", err)
	}

	err = server.Start(config.AppAddress, config.AppPort)
	if err != nil {
		log.Panic("cannot start server:", err)
	}
}
