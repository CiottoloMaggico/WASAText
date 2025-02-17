package database

import (
	"github.com/ciottolomaggico/wasatext/service/api/filter"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
)

type QueryParameters struct {
	Offset int
	Limit  int
	Cursor int64
	Filter string
}

func NewQueryParameters(ps pagination.PaginationParams, filter filter.Filter) (QueryParameters, error) {
	filterQuery, err := filter.Evaluate(ps.Filter)
	if err != nil {
		return QueryParameters{}, err
	}

	return QueryParameters{
		(ps.Page - 1) * ps.Size,
		ps.Size,
		ps.Cursor,
		filterQuery,
	}, nil
}
