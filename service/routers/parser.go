package routers

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

var validate = validator.New()

func ParseUrlParams(ps httprouter.Params, res interface{}) error {
	underlyingType, underlyingValue := reflect.TypeOf(res).Elem(), reflect.ValueOf(res).Elem()

	for i := 0; i < underlyingType.NumField(); i++ {
		field := underlyingType.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			underlyingValue.Field(i).SetString(ps.ByName(field.Name))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result, err := strconv.Atoi(ps.ByName(field.Name))
			if err != nil {
				return err
			}
			underlyingValue.Field(i).SetInt(int64(result))
		case reflect.Bool:
			result, err := strconv.ParseBool(ps.ByName(field.Name))
			if err != nil {
				return err
			}
			underlyingValue.Field(i).SetBool(result)
		}
	}

	return nil
}

func ParseAndValidatePaginationParams(query url.Values) (PaginationQueryParams, error) {
	res, page, size := PaginationQueryParams{}, query.Get("page"), query.Get("size")
	if page == "" {
		res.Page = 0
	}
	if size == "" {
		res.Size = DefaultPageSize
	}

	if err := validate.Struct(res); err != nil {
		return PaginationQueryParams{}, err
	}
	return res, nil
}

func ParseAndValidateUrlParams(ps httprouter.Params, res interface{}) error {
	if err := ParseUrlParams(ps, res); err != nil {
		return err
	}
	if err := validate.Struct(res); err != nil {
		return err
	}
	return nil
}

func ParseAndValidateRequestBody(req *http.Request, res interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(res); err != nil {
		return err
	}
	if err := validate.Struct(res); err != nil {
		return err
	}
	return nil
}
