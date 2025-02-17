package pagination

import (
	"errors"
	"net/url"
	"strconv"
)

type Page struct {
	Page         int `json:"page"`
	Cursor       int `json:"cursor"`
	finalPage    int
	CurrPage     string
	NextPage     *string `json:"nextPage"`
	PreviousPage *string `json:"previousPage"`
}

func (p Page) HasNext() bool {
	return p.Page+1 <= p.finalPage
}

func (p Page) HasPrevious() bool {
	return p.Page-1 >= 1
}

func (p Page) pageUrl(pageNum int) string {
	pageUrl, _ := url.Parse(p.CurrPage)
	reqQuery := pageUrl.Query()
	reqQuery.Set("page", strconv.Itoa(pageNum))
	reqQuery.Set("cursor", strconv.Itoa(p.Cursor))
	pageUrl.RawQuery = reqQuery.Encode()
	return pageUrl.RequestURI()
}

func (p Page) NextPageUrl() (string, error) {
	if !p.HasNext() {
		return "", errors.New("no next page")
	}
	return p.pageUrl(p.Page + 1), nil
}

func (p Page) PreviousPageUrl() (string, error) {
	if !p.HasPrevious() {
		return "", errors.New("no previous page")
	}
	return p.pageUrl(p.Page - 1), nil
}

func MakePage(ps PaginationParams, totalEntries int) Page {
	finalPage, remainingEntries := totalEntries/ps.Size, totalEntries%ps.Size
	if remainingEntries != 0 {
		finalPage++
	}

	res := Page{
		Page:         ps.Page,
		finalPage:    finalPage,
		Cursor:       int(ps.Cursor),
		NextPage:     nil,
		CurrPage:     ps.CurrentUrl,
		PreviousPage: nil,
	}
	if res.HasNext() {
		nextPageUrl, _ := res.NextPageUrl()
		res.NextPage = &nextPageUrl
	}
	if res.HasPrevious() {
		prevPageUrl, _ := res.PreviousPageUrl()
		res.PreviousPage = &prevPageUrl
	}
	return res
}
