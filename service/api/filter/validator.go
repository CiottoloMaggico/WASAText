package filter

import (
	"regexp"
	"strings"
)

var (
	operandRegex = regexp.MustCompile(
		`[a-zA-Z_]+\s(?:eq|ne|gt|ge|lt|le)\s(?:'(.*?)'|-?\d+.\d+|-?\d+)|[a-zA-Z]+\s(?:eq|ne)\s(?:null|true|false)`,
	)
	logicOperators = []string{"and", "or", "not"}
)

func validLogicOperator(op string) bool {
	for _, validOp := range logicOperators {
		if op == validOp {
			return true
		}
	}
	return false
}

func validOperatorSequence(prev string, curr string) bool {
	if !validLogicOperator(curr) && curr != "x" {
		return false
	}

	if (curr == "and" || curr == "or") && prev != "x" {
		return false
	} else if (curr == "x" || curr == "not") && prev == "x" {
		return false
	}

	return true
}

func validBaseFormula(formula string) bool {
	cleanedFormula := operandRegex.ReplaceAllString(formula, "x")
	prev, curr := "", ""
	for _, value := range cleanedFormula {
		if value == ' ' {
			if !validOperatorSequence(prev, curr) {
				return false
			}
			prev, curr = curr, ""
			continue
		}

		curr += string(value)
	}

	if curr != "" && (!validOperatorSequence(prev, curr) || validLogicOperator(curr)) {
		return false
	}
	return true
}

func FormulaIsValid(formula string) bool {
	formula = strings.Trim(formula, " ")
	subFormulas, currFormula := []string{}, ""

	for _, value := range formula {
		if value == '(' {
			subFormulas, currFormula = append(subFormulas, currFormula, string(value)), ""
			continue
		} else if value == ')' && len(subFormulas) > 0 {
			currFormula = strings.Trim(currFormula, " ")
			if !validBaseFormula(currFormula) {
				return false
			}
			subFormulas, currFormula = subFormulas[:len(subFormulas)-2], subFormulas[len(subFormulas)-2]+"x"
			continue
		} else if value == ')' {
			return false
		}

		currFormula += string(value)
	}

	if len(subFormulas) > 0 {
		return false
	} else if !validBaseFormula(currFormula) {
		return false
	}

	return true
}
