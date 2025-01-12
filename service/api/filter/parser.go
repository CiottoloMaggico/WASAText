package filter

import (
	"fmt"
	"reflect"
	"strings"
)

func filterTagParser(field reflect.StructField) (map[string]string, error) {
	tag, ok := field.Tag.Lookup("filter")
	if !ok || tag == "" {
		return nil, NewFilterError("invalid or missing tag", fmt.Sprintf("Missing or invalid filter tag on field %s", field.Name))
	} else if tag == "-" {
		return nil, nil
	}
	result := make(map[string]string)
	explodedTag := strings.Split(tag, ",")
	for _, op := range explodedTag {
		explodedOp := strings.SplitN(op, "=", 2)
		result[explodedOp[0]] = explodedOp[1]
	}

	return result, nil
}

func allowedValueRegexParser(field reflect.StructField) (string, error) {
	fieldType := field.Type
	result, allowNull := "", false

	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
		allowNull = true
	}

	switch fieldType.Kind() {
	case reflect.Bool:
		result = `(?:true|false)`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = `-?\d+`
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = `\d+`
	case reflect.Float32, reflect.Float64:
		result = `-?\d+.\d+`
	case reflect.String:
		result = `'(.*?)'`
	default:
		return "", NewFilterError("unsupported field type", fmt.Sprintf("unsupported field type %s", fieldType.Kind()))
	}

	if allowNull {
		return fmt.Sprintf("^(?:%s|null)$", result), nil
	}
	return fmt.Sprintf("^%s$", result), nil
}

func fieldMappingParser(model reflect.Type) (map[string][2]string, error) {
	result := make(map[string][2]string)
	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)
		tag, err := filterTagParser(field)
		if err != nil {
			return nil, err
		} else if tag == nil {
			continue
		}
		in, ok := tag["in"]
		if !ok {
			return nil, NewFilterError("missing 'in' tag", fmt.Sprintf("Missing in operation on: %s", field.Name))
		}
		out, ok := tag["out"]
		if !ok {
			out = in
		}

		allowedValueRegex, err := allowedValueRegexParser(field)
		if err != nil {
			return nil, err
		}
		result[in] = [2]string{out, allowedValueRegex}
	}

	return result, nil
}
