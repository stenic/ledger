package query

import (
	"strings"

	"github.com/stenic/ledger/internal/storage"
)

type VersionFilter struct {
	Application string
	Location    string
	Environment string
	Day         string
}

func (f VersionFilter) getWhere() (string, []any) {
	var query []string
	var values []any

	if f.Application != "" {
		query = append(query, "application LIKE ?")
		values = append(values, "%"+f.Application+"%")
	}
	if f.Environment != "" {
		query = append(query, "environment LIKE ?")
		values = append(values, "%"+f.Environment+"%")
	}
	if f.Location != "" {
		query = append(query, "location LIKE ?")
		values = append(values, "%"+f.Location+"%")
	}
	if f.Day != "" {
		switch storage.EngineType {
		case "mysql":
			query = append(query, "DATE(timestamp) = ?")
		case "sqlite":
			query = append(query, "DATE(timestamp) = ?")
		}
		values = append(values, f.Day)
	}

	if len(query) == 0 {
		return "", nil
	}

	return "(" + strings.Join(query, " AND ") + ")", values
}

func GetWhereParts(filter *VersionFilter) (string, []any) {
	where := "1=1"
	args := []any{}
	if filter != nil {
		if w, a := filter.getWhere(); w != "" {
			where = w
			args = a
		}
	}

	return where, args
}
