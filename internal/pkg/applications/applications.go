package applications

import (
	"database/sql"
	"log"

	"github.com/stenic/ledger/internal/storage"
)

type Application struct {
	Name string `json:"name"`
}

func GetAll() []Application {
	versions, err := runQuery("select distinct application from versions")
	if err != nil {
		log.Fatal(err)
	}
	return versions
}

func runQuery(query string) ([]Application, error) {
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
