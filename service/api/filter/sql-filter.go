package filter

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var sqlOperatorsMapping = map[string]string{
	"eq": "=",
	"ne": "!=",
	"gt": ">",
	"ge": ">=",
	"lt": "<",
	"le": "<=",
}

type SqlFilter struct {
	model reflect.Type
}

func NewSqlFilter(model interface{}) (SqlFilter, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Struct {
		return SqlFilter{}, NewFilterError("invalid model", "model must be a struct")
	}

	for iOp := range sqlOperatorsMapping {
		if ok := apiOperatorsRegex.MatchString(iOp); !ok {
			return SqlFilter{}, NewFilterError("unsupported operation", fmt.Sprintf("%s %s", iOp, " is not supported"))
		}
	}
	return SqlFilter{modelType}, nil
}

func (f SqlFilter) validateFormula(formula string) error {
	if !FormulaIsValid(formula) {
		return NewFilterError("invalid syntax", "the provided formula have an invalid syntax")
	}

	rawOperands := operandRegex.FindAllString(formula, -1)
	fieldMapping, err := fieldMappingParser(f.model)
	if err != nil {
		return err
	}
	for _, rawOperand := range rawOperands {
		explodedOp := strings.SplitN(rawOperand, " ", 3)

		if _, ok := sqlOperatorsMapping[explodedOp[1]]; !ok {
			return NewFilterError("invalid operator", fmt.Sprintf("%s is not supported", explodedOp[1]))
		}

		fieldInfo, ok := fieldMapping[explodedOp[0]]
		if !ok {
			return NewFilterError("invalid field", fmt.Sprintf("%s is not supported", explodedOp[0]))
		}

		if ok, err := regexp.MatchString(fieldInfo[1], explodedOp[2]); !ok || err != nil {
			return NewFilterError("invalid field value", fmt.Sprint("field '", explodedOp[0], "' have an invalid value: ", explodedOp[2]))
		}
	}

	return nil
}

func (f SqlFilter) Evaluate(formula string) (string, error) {
	if formula == "" {
		return "", nil
	}

	if err := f.validateFormula(formula); err != nil {
		return "", err
	}
	fieldMap, err := fieldMappingParser(f.model)
	if err != nil {
		return "", err
	}

	var mappedField, mappedOp, mappedOperand string
	result := formula
	rawOperands := operandRegex.FindAllString(result, -1)
	for _, rawOperand := range rawOperands {
		explodedOperand := strings.SplitN(rawOperand, " ", 3)

		mappedField = fieldMap[explodedOperand[0]][0]
		mappedOperand = fmt.Sprintf("%s %s", explodedOperand[1], explodedOperand[2])
		switch mappedOperand {
		case "eq null":
			mappedOp = "IS"
		case "ne null":
			mappedOp = "IS NOT"
		default:
			mappedOp = sqlOperatorsMapping[explodedOperand[1]]
		}

		mappedOperand = fmt.Sprintf("%s %s %s", mappedField, mappedOp, explodedOperand[2])
		result = strings.Replace(result, rawOperand, mappedOperand, -1)
	}

	return result, nil
}
