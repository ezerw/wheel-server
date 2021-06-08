package db

import (
	"database/sql"
	"fmt"
	"github.com/ezerw/wheel/util"
	"net/url"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

func Connect(config util.Config) (*sql.DB, error) {
	dbDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s&multiStatements=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		url.QueryEscape(config.AppTimezone),
	)
	conn, err := sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
