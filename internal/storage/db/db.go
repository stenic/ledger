package db

import "database/sql"

type DatabaseEngine interface {
	InitDB() *sql.DB
	CloseDB()
	Migrate() error
}

type DatabaseOpts struct {
	Type string
	Dsn  string
}
