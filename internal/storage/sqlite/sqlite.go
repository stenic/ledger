package sqlite

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/storage/db"
)

type SqliteEngine struct {
	db *sql.DB
}

func NewEngine() db.DatabaseEngine {
	return &SqliteEngine{}
}

func (e *SqliteEngine) InitDB() *sql.DB {
	datadir := "./data"
	if _, err := os.Stat(datadir); os.IsNotExist(err) {
		if err = os.Mkdir(datadir, 0755); err != nil {
			logrus.Fatal(err)
		}
	}

	idb, err := sql.Open("sqlite3", datadir+"/ledger.db")
	if err != nil {
		logrus.Fatal(err)
	}

	if err = idb.Ping(); err != nil {
		logrus.Fatal(err)
	}

	e.db = idb
	return e.db
}
func (e *SqliteEngine) CloseDB() {
	e.db.Close()
}

func (e *SqliteEngine) Migrate() error {
	if err := e.db.Ping(); err != nil {
		return err
	}
	driver, err := sqlite3.WithInstance(e.db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/sqlite",
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}
	return m.Up()
}
