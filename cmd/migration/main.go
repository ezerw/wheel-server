package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/util"
)

func main() {
	logger := util.NewLogger()

	if len(os.Args) == 1 {
		log.Fatal("Missing options up or down")
	}
	option := os.Args[1]

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dBConn, err := db.Connect(config)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer dBConn.Close()

	driver, err := mysql.WithInstance(dBConn, &mysql.Config{})
	if err != nil {
		log.Panic("cannot get db instance for migration")
	}
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "wheel", driver)
	if err != nil {
		log.Panic("cannot initialize migrations:", err)
	}

	switch option {
	case "new":
		if len(os.Args) != 3 {
			log.Panic("migration name is missing: e.g: \"new create-user\"")
		}

		err = createMigrationFiles(os.Args[2])
		if err != nil {
			logger.WithError(err).Error("failed to create migration")
			return
		}
		logger.Printf("New migration created: %s\n", os.Args[2])
	case "ver":
		ver, _, err := m.Version()
		if err != nil {
			logger.WithError(err).Error("failed to get migration version")
			return
		}
		fmt.Println("current migration version:", ver)
	case "up":
		if err = m.Up(); err != nil {
			logger.WithError(err).Error("failed to run migrate up")
			return
		}
		fmt.Println("migrations applied")
	case "down":
		if err = m.Down(); err != nil {
			logger.WithError(err).Error("failed to run migrate down")
			return
		}
		fmt.Println("migrations reverted")
	}
}

// createMigrationFiles creates up and down migration files.
func createMigrationFiles(name string) error {
	timestamp := time.Now().Format("20060102150405")

	// migration up
	file1 := fmt.Sprintf("./db/migrations/%s_%s.up.sql", timestamp, name)
	f1, err := os.Create(file1)
	if err != nil {
		return err
	}
	defer f1.Close()

	// migration down
	file2 := fmt.Sprintf("./db/migrations/%s_%s.down.sql", timestamp, name)
	f2, err := os.Create(file2)
	if err != nil {
		return err
	}
	defer f2.Close()

	return nil
}
