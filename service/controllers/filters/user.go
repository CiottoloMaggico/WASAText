package filters

import "github.com/ciottolomaggico/wasatext/service/api/filter"

type UserFilterMap struct {
	Uuid     string `filter:"in=uuid,out=user_uuid"`
	Username string `filter:"in=username,out=user_username"`
	Photo    string `filter:"in=photo,out=user_photo"`
}

func NewUserFilter() (filter.SqlFilter, error) {
	return filter.NewSqlFilter(UserFilterMap{})
}
