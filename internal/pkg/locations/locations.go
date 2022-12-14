package locations

import (
	"database/sql"
	"log"

	"github.com/stenic/ledger/internal/storage"
)

type Location struct {
	Name string `json:"name"`
}

func GetAll() []Location {
	versions, err := runQuery("select distinct location from versions")
	if err != nil {
		log.Fatal(err)
	}
	return versions
}

func runQuery(query string) ([]Location, error) {
	stmt, err := storage.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mapRows(rows)
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
