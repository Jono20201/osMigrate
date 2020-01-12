package main

import (
	"database/sql"
	"fmt"
)

type ConnectionParams struct {
	hostname string
	port     uint
	username string
	password string
	schema   string
}

func connect(params ConnectionParams) *sql.DB {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/?multiStatements=true", params.username, params.password, params.hostname, params.port)
	pool, err := sql.Open("mysql", connStr)

	if err != nil {
		logError(err)
	}

	if err := pool.Ping(); err != nil {
		logError(err)
	}

	_, err = pool.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %[1]s; USE %[1]s", params.schema))

	if err != nil {
		logError(err)
	}

	return pool
}
