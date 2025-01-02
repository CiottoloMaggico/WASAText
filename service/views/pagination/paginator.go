package pagination

import (
	"errors"
	"reflect"
)

type Paginator interface {
	HasNext() bool
	HasPrevious() bool
	NextPageUrl() (string, error)
	PreviousPageUrl() (string, error)
	pageUrl(pageNum int) string
}

type PaginationParams struct {
	Page       int `validate:"min=1"`
	Size       int `validate:"min=1,max=20"`
	Filter     string
	CurrentUrl string
}

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
	page := MakePage(ps.Page, ps.Size, totalEntries, ps.CurrentUrl)
	return newPaginatedView(page, content)
}
