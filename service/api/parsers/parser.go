package parsers

import (
	"encoding/json"
	"errors"
	api_errors "github.com/ciottolomaggico/wasatext/service/api/api-errors"
	"github.com/ciottolomaggico/wasatext/service/validators"
	"github.com/ciottolomaggico/wasatext/service/views/pagination"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

const DEFAULT_PAGE_SIZE = 20

func renderValidationErrors(errs validator.ValidationErrors) map[string]string {
	res := make(map[string]string)
	for _, fieldErr := range errs {
		res[fieldErr.Field()] = fieldErr.Error()
	}
	return res
}

func ParseAndValidatePaginationParams(url *url.URL) (pagination.PaginationParams, error) {
	query := url.Query()
	res := pagination.PaginationParams{
		Page:       1,
		Size:       DEFAULT_PAGE_SIZE,
		CurrentUrl: url.String(),
		Filter:     query.Get("filter"),
	}

	if page := query.Get("page"); page != "" {
		tmpPage, err := strconv.Atoi(page)
		if err != nil {
			return pagination.PaginationParams{}, api_errors.InvalidUrlParameters()
		}
		res.Page = tmpPage
	}

	if size := query.Get("size"); size != "" {
		tmpSize, err := strconv.Atoi(size)
		if err != nil {
			return pagination.PaginationParams{}, api_errors.InvalidUrlParameters()
		}
		res.Size = tmpSize
	}

	if err := validators.Validate.Struct(res); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return pagination.PaginationParams{},
				api_errors.UnprocessableContent(renderValidationErrors(errs))
		}
		return pagination.PaginationParams{}, err
	}
	return res, nil
}

func ParseUrlParams(ps httprouter.Params, res interface{}) error {
	underlyingType, underlyingValue := reflect.TypeOf(res).Elem(), reflect.ValueOf(res).Elem()

	for i := 0; i < underlyingType.NumField(); i++ {
		var fieldValue reflect.Value
		var result interface{}
		var err error = nil
		field := underlyingType.Field(i)
		urlRawValue := ps.ByName(field.Tag.Get("url"))

		switch field.Type.Kind() {
		case reflect.String:
			result = urlRawValue
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result, err = strconv.ParseInt(urlRawValue, 10, 0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result, err = strconv.ParseUint(urlRawValue, 10, 0)
		case reflect.Float32, reflect.Float64:
			result, err = strconv.ParseFloat(urlRawValue, 0)
		case reflect.Bool:
			result, err = strconv.ParseBool(urlRawValue)
		default:
			return errors.New("Unsupported field type: " + field.Type.String())
		}

		if err != nil {
			return api_errors.InvalidUrlParameters()
		}
		fieldValue = reflect.ValueOf(result)
		underlyingValue.Field(i).Set(fieldValue)
	}

	return nil
}

func ParseMultipartRequestBody(body *multipart.Form, res interface{}) error {
	underlyingType, underlyingValue := reflect.TypeOf(res).Elem(), reflect.ValueOf(res).Elem()

	for i := 0; i < underlyingType.NumField(); i++ {
		field, fieldValue := underlyingType.Field(i), reflect.ValueOf(nil)
		pointedType, multipartFieldName := field.Type, field.Tag.Get("form")

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
					return api_errors.InvalidMultipartBody()
				}
				fieldValue = reflect.ValueOf(&result)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				result, err := strconv.ParseUint(val[0], 10, 0)
				if err != nil {
					return api_errors.InvalidMultipartBody()
				}
				fieldValue = reflect.ValueOf(&result)
			case reflect.Float32, reflect.Float64:
				result, err := strconv.ParseFloat(val[0], 0)
				if err != nil {
					return api_errors.InvalidMultipartBody()
				}
				fieldValue = reflect.ValueOf(&result)
			case reflect.Bool:
				result, err := strconv.ParseBool(val[0])
				if err != nil {
					return api_errors.InvalidMultipartBody()
				}
				fieldValue = reflect.ValueOf(&result)
			default:
				return errors.New("Unsupported field type: " + field.Type.String())
			}
			if field.Type.Kind() != reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			underlyingValue.Field(i).Set(fieldValue)
		} else if val, ok := body.File[multipartFieldName]; ok {
			fieldValue = reflect.ValueOf(val[0])
			underlyingValue.Field(i).Set(fieldValue)
		}

	}
	return nil
}

func ParseAndValidateUrlParams(ps httprouter.Params, res interface{}) error {
	if err := ParseUrlParams(ps, res); err != nil {
		return err
	}
	if err := validators.Validate.Struct(res); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return api_errors.UnprocessableContent(renderValidationErrors(errs))
		}
		return err
	}
	return nil
}

func ParseAndValidateMultipartRequestBody(req *http.Request, res interface{}) error {
	if err := req.ParseMultipartForm(0); err != nil {
		return api_errors.InvalidMultipartBody()
	}
	body := req.MultipartForm
	if err := ParseMultipartRequestBody(body, res); err != nil {
		return err
	}
	if err := validators.Validate.Struct(res); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return api_errors.UnprocessableContent(renderValidationErrors(errs))
		}
		return err
	}
	return nil
}

func ParseAndValidateRequestBody(req *http.Request, res interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(res); err != nil {
		return api_errors.InvalidJson()
	}
	if err := validators.Validate.Struct(res); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return api_errors.UnprocessableContent(renderValidationErrors(errs))
		}
		return err
	}
	return nil
}
