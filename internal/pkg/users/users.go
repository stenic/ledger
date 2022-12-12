package users

import (
	"database/sql"

	"github.com/stenic/ledger/internal/storage"
)

type User struct {
	ID       string
	Username string
	Password string
}

func FindByUsername(username string) *User {
	users, err := runQuery("SELECT id, username, password FROM users WHERE username = ?", username)
	if err != nil {
		return nil
	}
	if len(users) != 1 {
		return nil
	}
	return &users[0]
}

func runQuery(query string, args ...string) ([]User, error) {
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

	return mapData(rows)
}

func mapData(rows *sql.Rows) ([]User, error) {
	var data []User
	for rows.Next() {
		var item User
		err := rows.Scan(
			&item.ID,
			&item.Username,
			&item.Password,
		)
		if err != nil {
			return data, err
		}
		data = append(data, item)
	}
	return data, nil
}
