package applications

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/pkg/query"
	"github.com/stenic/ledger/internal/storage"
)

type Application struct {
	Name string `json:"name"`
}

func GetAll(filter *query.VersionFilter) []Application {
	where, args := query.GetWhereParts(filter)
	rows, err := storage.Db.Query(fmt.Sprintf("SELECT DISTINCT application FROM versions WHERE %s", where), args...)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	applications, err := mapRows(rows)
	if err != nil {
		log.Fatal(err)
	}
	return applications
}

func mapRows(rows *sql.Rows) ([]Application, error) {
	var result []Application
	for rows.Next() {
		var item Application
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
