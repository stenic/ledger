package environments

import (
	"database/sql"
	"log"

	"github.com/stenic/ledger/internal/storage"
)

type Environment struct {
	Name string `json:"name"`
}

func GetAll() []Environment {
	versions, err := runQuery("select distinct environment from versions")
	if err != nil {
		log.Fatal(err)
	}
	return versions
}

func runQuery(query string) ([]Environment, error) {
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
