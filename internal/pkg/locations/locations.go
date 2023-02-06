package locations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/pkg/query"
	"github.com/stenic/ledger/internal/storage"
)

type Location struct {
	Name string `json:"name"`
}

func GetAll(filter *query.VersionFilter) []Location {
	where, args := query.GetWhereParts(filter)
	rows, err := storage.Db.Query(fmt.Sprintf("SELECT DISTINCT location FROM versions WHERE %s", where), args...)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	locations, err := mapRows(rows)
	if err != nil {
		log.Fatal(err)
	}
	return locations
}

func mapRows(rows *sql.Rows) ([]Location, error) {
	var result []Location
	for rows.Next() {
		var item Location
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
