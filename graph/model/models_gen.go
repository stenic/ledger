// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Application struct {
	Name string `json:"name"`
}

type AuthPayload struct {
	Token string `json:"token"`
}

type DateVersionCount struct {
	Timstamp string `json:"timstamp"`
	Count    int    `json:"count"`
}

type Environment struct {
	Name string `json:"name"`
}

type Location struct {
	Name string `json:"name"`
}

type NewVersion struct {
	Application string  `json:"application"`
	Environment string  `json:"environment"`
	Location    *string `json:"location"`
	Version     string  `json:"version"`
}

type Version struct {
	ID          string       `json:"id"`
	Application *Application `json:"application"`
	Environment *Environment `json:"environment"`
	Location    *Location    `json:"location"`
	Version     string       `json:"version"`
	Timestamp   string       `json:"timestamp"`
}

type VersionOrderByInput struct {
	Application *Sort `json:"application"`
	Environment *Sort `json:"environment"`
	Location    *Sort `json:"location"`
	Timestamp   *Sort `json:"timestamp"`
}

type Sort string

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
)

var AllSort = []Sort{
	SortAsc,
	SortDesc,
}

func (e Sort) IsValid() bool {
	switch e {
	case SortAsc, SortDesc:
		return true
	}
	return false
}

func (e Sort) String() string {
	return string(e)
}

func (e *Sort) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Sort(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Sort", str)
	}
	return nil
}

func (e Sort) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
