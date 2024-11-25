package api

import (
	"errors"
	"net/url"
	"strconv"
)

type Paginator interface {
	HasNext() bool
	HasPrevious() bool
	NextPageUrl() string
	PreviousPageUrl() string
	pageUrl(pageNum int) string
}

type Page struct {
	Page         int     `json:"page"`
	finalPage    int     `json:"-"`
	NextPage     *string `json:"nextPage"`
	CurrPage     string  `json:"currPage"`
	PreviousPage *string `json:"previousPage"`
}

func (p Page) HasNext() bool {
	if p.Page+1 > p.finalPage {
		return false
	}
	return true
}

func (p Page) HasPrevious() bool {
	if p.Page-1 < 0 {
		return false
	}
	return true
}

func (p Page) pageUrl(pageNum int) string {
	pageUrl, _ := url.Parse(p.CurrPage)
	reqQuery := pageUrl.Query()
	reqQuery.Set("page", strconv.Itoa(pageNum))
	pageUrl.RawQuery = reqQuery.Encode()
	return pageUrl.RequestURI()
}

func (p Page) NextPageUrl() (string, error) {
	if !p.HasNext() {
		return "", errors.New("No next page")
	}
	return p.pageUrl(p.Page + 1), nil
}

func (p Page) PreviousPageUrl() (string, error) {
	if !p.HasPrevious() {
		return "", errors.New("No previous page")
	}
	return p.pageUrl(p.Page - 1), nil
}

func MakePage(pageNum int, finalPage int, currPage string) *Page {
	page := &Page{
		Page:         pageNum,
		finalPage:    finalPage,
		NextPage:     nil,
		CurrPage:     currPage,
		PreviousPage: nil,
	}
	if page.HasNext() {
		nextPageUrl, _ := page.NextPageUrl()
		page.NextPage = &nextPageUrl
	}
	if page.HasPrevious() {
		prevPageUrl, _ := page.PreviousPageUrl()
		page.PreviousPage = &prevPageUrl
	}
	return page
}
