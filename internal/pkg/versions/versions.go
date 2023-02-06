package versions

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/pkg/query"
	"github.com/stenic/ledger/internal/storage"
)

type Version struct {
	ID          int64  `json:"id"`
	Application string `json:"application"`
	Environment string `json:"environment"`
	Location    string `json:"location"`
	Timestamp   string `json:"timestamp"`
	Version     string `json:"version"`
}

type VersionOrderBy struct {
	Application query.Sort
	Environment query.Sort
	Location    query.Sort
	Timestamp   query.Sort
}

func (vop VersionOrderBy) GetParts() query.QueryParts {
	qp := query.QueryParts{}
	if vop.Timestamp != query.SortUndefined {
		qp.Order = append(qp.Order, fmt.Sprintf("timestamp %s", vop.Timestamp))
	}
	if vop.Application != query.SortUndefined {
		qp.Order = append(qp.Order, fmt.Sprintf("application %s", vop.Application))
	}
	if vop.Environment != query.SortUndefined {
		qp.Order = append(qp.Order, fmt.Sprintf("environment %s", vop.Environment))
	}
	if vop.Location != query.SortUndefined {
		qp.Order = append(qp.Order, fmt.Sprintf("location %s", vop.Location))
	}

	return qp
}

func (v Version) Save() int64 {
	stmt, err := storage.Db.Prepare("INSERT INTO versions (location, environment, application, version) VALUES (?, ?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(v.Location, v.Environment, v.Application, v.Version)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	v.ID = id

	return id
}

func GetLast(location, environment, application string) *Version {
	stmtFind, err := storage.Db.Prepare("SELECT id, application, environment, location, version, timestamp FROM versions WHERE location = ? AND environment = ? AND application = ? ORDER BY id DESC LIMIT 1")
	if row := stmtFind.QueryRow(location, environment, application); err == nil && row != nil {
		version := Version{}
		err := row.Scan(
			&version.ID,
			&version.Application,
			&version.Environment,
			&version.Location,
			&version.Version,
			&version.Timestamp,
		)
		if err == nil {
			return &version
		}
	}

	return nil
}

type CountResult struct {
	Timestamp string
	Count     int
}

func CountByDay(filter *query.VersionFilter) []CountResult {
	where, args := query.GetWhereParts(filter)
	rows, err := storage.Db.Query(fmt.Sprintf("SELECT count(id), DATE(`timestamp`) AS date_formatted FROM versions WHERE %s GROUP BY date_formatted", where), args...)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	var items []CountResult
	for rows.Next() {
		var item CountResult
		err := rows.Scan(&item.Count, &item.Timestamp)
		item.Timestamp = item.Timestamp[0:10]
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	return items
}

func CountTotal(filter *query.VersionFilter) *int {
	where, args := query.GetWhereParts(filter)
	if row := storage.Db.QueryRow(fmt.Sprintf("SELECT count(id) FROM versions WHERE %s", where), args...); row != nil {
		var count int
		if err := row.Scan(&count); err != nil {
			logrus.WithField("repo", "versions.CountTotal").Error(err)
		}
		return &count
	}

	return nil
}

func GetAll(opts ...query.QueryOpts) []Version {
	q := query.AddQueryParts("select id, application, environment, location, version, timestamp from versions", opts...)
	versions, err := runQuery(q)
	if err != nil {
		log.Fatal(err)
	}
	return versions
}

func GetAllLast(days *int) []Version {
	var filter string
	if days != nil {
		switch storage.EngineType {
		case "mysql":
			filter = fmt.Sprintf("WHERE timestamp > NOW() - INTERVAL %d DAY", *days)
		case "sqlite":
			filter = fmt.Sprintf("where timestamp > date('now', '-%d days')", *days)
		}
	}
	q := query.AddQueryParts(fmt.Sprintf(
		"select %s from versions where id in (select max(id) from versions %s group by location, environment, application) order by timestamp desc",
		fields,
		filter,
	))
	versions, err := runQuery(q)
	if err != nil {
		log.Fatal(err)
	}
	return versions
}

func runQuery(query string) ([]Version, error) {
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

	return mapVersions(rows)
}

const fields = "id, application, environment, location, version, timestamp"

func mapVersions(rows *sql.Rows) ([]Version, error) {
	var versions []Version
	for rows.Next() {
		var version Version
		err := rows.Scan(
			&version.ID,
			&version.Application,
			&version.Environment,
			&version.Location,
			&version.Version,
			&version.Timestamp,
		)
		if err != nil {
			return versions, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}
