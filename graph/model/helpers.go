package model

import (
	"github.com/stenic/ledger/internal/pkg/query"
)

func (f VersionFilter) GetFilter() query.VersionFilter {
	return query.VersionFilter{
		Application: *f.Application,
		Location:    *f.Location,
		Environment: *f.Environment,
		Day:         *f.Day,
	}
}

func NewFilter(filter *VersionFilter) *query.VersionFilter {
	if filter == nil {
		return nil
	}
	f := filter.GetFilter()
	return &f
}
