package main

import (
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/db/seeds"
	"github.com/ezerw/wheel/util"
)

func main() {
	logger := util.NewLogger()

	config, err := util.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dBConn, err := db.Connect(config)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer dBConn.Close()

	store := db.NewStore(dBConn)
	for _, seed := range seeds.All() {
		if err = seed.Run(context.Background(), &store); err != nil {
			logger.WithError(err).
				WithField("seeder", seed.Name).
				Error("error running seeder")
			return
		}
	}
}
