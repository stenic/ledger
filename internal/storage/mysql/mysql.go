package mysql

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stenic/ledger/internal/storage/db"
)

type MysqlEngine struct {
	db *sql.DB
}

func NewEngine() db.DatabaseEngine {
	return &MysqlEngine{}
}

func (e *MysqlEngine) InitDB() *sql.DB {
	idb, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		logrus.Fatal(err)
	}

	if err = idb.Ping(); err != nil {
		logrus.Fatal(err)
	}

	e.db = idb
	return e.db
}

func (e *MysqlEngine) CloseDB() {
	e.db.Close()
}
func (e *MysqlEngine) Migrate() error {
	if err := e.db.Ping(); err != nil {
		return err
	}
	driver, err := mysql.WithInstance(e.db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/mysql",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}
	return m.Up()
}
