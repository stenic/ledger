package query

import "strings"

type Sort int

const (
	SortUndefined Sort = iota
	SortAsc
	SortDesc
)

func (s Sort) String() string {
	switch s {
	case SortAsc:
		return "asc"
	case SortDesc:
		return "desc"
	}
	return ""
}

func SortFromString(s string) Sort {
	switch strings.ToLower(s) {
	case "asc":
		return SortAsc
	case "desc":
		return SortDesc
	}
	return SortUndefined
}
