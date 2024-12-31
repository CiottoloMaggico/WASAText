package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
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
		default:
			return errors.New("Unsupported field type: " + field.Type.String())
		}
	}

	return nil
}

func ParseAndValidatePaginationParams(url *url.URL) (pagination.PaginationParams, error) {
	query := url.Query()
	res, page, size, _ := pagination.PaginationParams{
		CurrentUrl: url.String(),
	}, query.Get("page"), query.Get("size"), query.Get("filter")
	if page == "" {
		res.Page = 1
	} else {
		tmpPage, err := strconv.Atoi(page)
		if err != nil {
			return pagination.PaginationParams{}, err
		}
		res.Page = tmpPage
	}
	if size == "" {
		res.Size = DefaultPageSize
	} else {
		tmpSize, err := strconv.Atoi(size)
		if err != nil {
			return pagination.PaginationParams{}, err
		}
		res.Size = tmpSize
	}

	if err := validate.Struct(res); err != nil {
		return pagination.PaginationParams{}, err
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

func ParseMultipartRequestBody(body *multipart.Form, res interface{}) error {
	underlyingType, underlyingValue := reflect.TypeOf(res).Elem(), reflect.ValueOf(res).Elem()

	for i := 0; i < underlyingType.NumField(); i++ {
		var fieldValue reflect.Value
		field := underlyingType.Field(i)
		pointedType, multipartFieldName := field.Type, field.Tag.Get("in")

		if field.Type.Kind() == reflect.Ptr {
			pointedType = field.Type.Elem()
		}

		if val, ok := body.Value[multipartFieldName]; ok {
			switch pointedType.Kind() {
			case reflect.String:
				fieldValue = reflect.ValueOf(&val[0])
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result, err := strconv.ParseInt(val[0], 10, 0)
				if err != nil {
					return err
				}
				fieldValue = reflect.ValueOf(&result)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				result, err := strconv.ParseUint(val[0], 10, 0)
				if err != nil {
					return err
				}
				fieldValue = reflect.ValueOf(&result)
			case reflect.Float32, reflect.Float64:
				result, err := strconv.ParseFloat(val[0], 0)
				if err != nil {
					return err
				}
				fieldValue = reflect.ValueOf(&result)
			case reflect.Bool:
				result, err := strconv.ParseBool(val[0])
				if err != nil {
					return err
				}
				fieldValue = reflect.ValueOf(&result)
			default:
				return errors.New("Unsupported field type: " + field.Type.String())
			}

			if field.Type.Kind() != reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
		} else if val, ok := body.File[multipartFieldName]; ok {
			fieldValue = reflect.ValueOf(val[0])
		} else {
			return fmt.Errorf("please provide the field %s", multipartFieldName)
		}

		underlyingValue.Field(i).Set(fieldValue)
	}
	return nil
}

func ParseAndValidateMultipartRequestBody(req *http.Request, res interface{}) error {
	if err := req.ParseMultipartForm(0); err != nil {
	}
	body := req.MultipartForm
	if err := ParseMultipartRequestBody(body, res); err != nil {
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
