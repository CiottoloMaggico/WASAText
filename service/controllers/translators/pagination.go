package translators

import "github.com/ciottolomaggico/wasatext/service/views/pagination"

func ToPaginatedView(ps pagination.PaginationParams, totalEntries int, content interface{}) (pagination.PaginatedView, error) {
	page := pagination.MakePage(ps.Page, ps.Size, totalEntries, ps.CurrentUrl)
	return pagination.NewPaginatedView(page, content)
}
