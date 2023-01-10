package storage

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"

	"github.com/stenic/ledger/internal/storage/db"
	"github.com/stenic/ledger/internal/storage/mysql"
	"github.com/stenic/ledger/internal/storage/sqlite"
)

var Db *sql.DB
var EngineType string

type Database struct {
	db     *sql.DB
	engine db.DatabaseEngine
}

func (d *Database) InitDB() {
	engineType := "sqlite"
	if os.Getenv("MYSQL_DSN") != "" {
		engineType = "mysql"
	}

	switch engineType {
	case "mysql":
		d.engine = mysql.NewEngine()
	case "sqlite":
		d.engine = sqlite.NewEngine()
		logrus.Warn("Running with sqlite, this should only be used in development")
	}
	logrus.Debugf("Setting up DB with %s", engineType)
	d.db = d.engine.InitDB()

	Db = d.db
	EngineType = engineType
}

func (d *Database) CloseDB() {
	d.engine.CloseDB()
}

func (d *Database) Migrate() error {
	err := d.engine.Migrate()
	switch err {
	case nil:
		logrus.Debug("Migration complete")
	case migrate.ErrNoChange:
		logrus.Debug("No migration needed")
	default:
		logrus.Fatal(err)
		return err
	}
	return nil
}
