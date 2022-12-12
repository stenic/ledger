package query

import (
	"fmt"
	"strings"
)

type QueryParts struct {
	Where []string
	Order []string
	Limit []string
}

type QueryOpts interface {
	GetParts() QueryParts
}

func AddQueryParts(q string, parts ...QueryOpts) string {
	c := QueryParts{}
	for _, p := range parts {
		c.Limit = append(c.Limit, p.GetParts().Limit...)
		c.Order = append(c.Order, p.GetParts().Order...)
		c.Where = append(c.Where, p.GetParts().Where...)
	}

	if len(c.Where) > 0 {
		q = fmt.Sprintf(q+" WHERE %s", strings.Join(c.Where, " AND "))
	}
	if len(c.Order) > 0 {
		q = fmt.Sprintf(q+" ORDER BY %s", strings.Join(c.Order, ", "))
	}
	if len(c.Limit) > 0 {
		q = fmt.Sprintf(q+" LIMIT %s", c.Limit)
	}

	return q
}
