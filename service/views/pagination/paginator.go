package pagination

type Paginator interface {
	HasNext() bool
	HasPrevious() bool
	NextPageUrl() (string, error)
	PreviousPageUrl() (string, error)
	pageUrl(pageNum int) string
}

type PaginationParams struct {
	Page       int    `validate:"min=1"`
	Size       int    `validate:"min=1,max=20"`
	Cursor     int64  `validate:"min=-1"`
	CurrentUrl string `validate:"required"`
	Filter     string `validate:"omitempty,formula"`
}
