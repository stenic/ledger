package environments

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/pkg/query"
	"github.com/stenic/ledger/internal/storage"
)

type Environment struct {
	Name string `json:"name"`
}

func GetAll(filter *query.VersionFilter) []Environment {
	where, args := query.GetWhereParts(filter)
	rows, err := storage.Db.Query(fmt.Sprintf("SELECT DISTINCT environment FROM versions WHERE %s", where), args...)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	environments, err := mapRows(rows)
	if err != nil {
		log.Fatal(err)
	}
	return environments
}

func mapRows(rows *sql.Rows) ([]Environment, error) {
	var result []Environment
	for rows.Next() {
		var item Environment
		err := rows.Scan(
			&item.Name,
		)
		if err != nil {
			return result, err
		}
		result = append(result, item)
	}
	return result, nil
}
