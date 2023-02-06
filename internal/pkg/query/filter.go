package query

import "strings"

type VersionFilter struct {
	Application string
	Location    string
	Environment string
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
