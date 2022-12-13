package client

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/storage"
)

const LocalClientIssuer = "ledger-local-client"

type Client struct {
	ID        string
	Username  string
	LastUsed  time.Time
	CreatedAt time.Time
}

func CreateClient(username string) (*Client, error) {
	client := Client{
		ID:       uuid.New().String(),
		Username: username,
	}

	stmt, err := storage.Db.Prepare("INSERT INTO clients (id, username, last_usage, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		client.ID,
		client.Username,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func FindByUsername(username string) []Client {
	items, err := runQuery("SELECT id, username, last_usage, created_at FROM clients WHERE username = ?", username)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}

func FindByID(id string) []Client {
	items, err := runQuery("SELECT id, username, last_usage, created_at FROM clients WHERE id = ?", id)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}

func UpdateLastUsageByID(id string) error {
	stmt, err := storage.Db.Prepare("UPDATE clients SET last_usage = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		time.Now(),
		id,
	)
	return err
}

func runQuery(query string, args ...any) ([]Client, error) {
	stmt, err := storage.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mapData(rows)
}

func mapData(rows *sql.Rows) ([]Client, error) {
	var data []Client
	for rows.Next() {
		var item Client
		err := rows.Scan(
			&item.ID,
			&item.Username,
			&item.LastUsed,
			&item.CreatedAt,
		)
		if err != nil {
			return data, err
		}
		data = append(data, item)
	}
	return data, nil
}
