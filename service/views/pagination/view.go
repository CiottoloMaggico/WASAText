package pagination

import (
	"errors"
	"reflect"
)

type PaginatedView struct {
	Page    Paginator   `json:"page"`
	Content interface{} `json:"content"`
}

func newPaginatedView(page Paginator, content interface{}) (PaginatedView, error) {
	contentType := reflect.TypeOf(content)
	if kind := contentType.Kind(); kind != reflect.Array && kind != reflect.Slice {
		return PaginatedView{}, errors.New("content must be an array or slice")
	}

	return PaginatedView{
		page,
		content,
	}, nil
}

func ToPaginatedView(ps PaginationParams, totalEntries int, content interface{}) (PaginatedView, error) {
	page := MakePage(ps, totalEntries)
	return newPaginatedView(page, content)
}
