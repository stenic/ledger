package versions

import (
	"database/sql"
	"fmt"
	"log"

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

func GetAll(opts ...query.QueryOpts) []Version {
	q := query.AddQueryParts("select id, application, environment, location, version, timestamp from versions", opts...)
	versions, err := runQuery(q)
	if err != nil {
		log.Fatal(err)
	}
	return versions
}

func GetLast() []Version {
	q := query.AddQueryParts("select id, application, environment, location, version, timestamp from versions where id in (select max(id) from versions group by location, environment, application)")
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
