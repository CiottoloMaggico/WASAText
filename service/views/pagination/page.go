package pagination

import (
	"errors"
	"net/url"
	"strconv"
)

type Page struct {
	Page         int     `json:"page"`
	finalPage    int     `json:"-"`
	CurrPage     string  `json:"-"`
	NextPage     *string `json:"nextPage"`
	PreviousPage *string `json:"previousPage"`
}

func (p Page) HasNext() bool {
	if p.Page+1 > p.finalPage {
		return false
	}
	return true
}

func (p Page) HasPrevious() bool {
	if p.Page-1 < 1 {
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

func MakePage(page int, pageSize int, totalEntries int, currentUrl string) Page {
	finalPage, remainingEntries := totalEntries/pageSize, totalEntries%pageSize
	if remainingEntries != 0 {
		finalPage++
	}
	res := Page{
		Page:         page,
		finalPage:    finalPage,
		NextPage:     nil,
		CurrPage:     currentUrl,
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
