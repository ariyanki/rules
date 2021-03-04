package rules

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

type (
	//Condition ..
	Condition struct {
		FunctionData  string `json:"function_data"`
		Data          string `json:"data"`
		Operator      string `json:"operator"`
		FunctionValue string `json:"function_value"`
		Value         string `json:"value"`
		Logic         string `json:"logic"`
	}
)

var (
	//Operators for Biller Config Rules
	Operators = map[string]map[string]string{
		"empty":      {"operator": "", "type": "general"},
		"not_empty":  {"operator": "", "type": "general"},
		"start_with": {"operator": "", "type": "general"},
		"end_with":   {"operator": "", "type": "general"},
		"in":         {"operator": "", "type": "general"},
		"not_in":     {"operator": "", "type": "general"},
		"set":        {"operator": "", "type": "general"},
		"not_set":    {"operator": "", "type": "general"},

		"equal":                     {"operator": "==", "type": "general"},
		"not_equal":                 {"operator": "!=", "type": "general"},
		"num_greater_than":          {"operator": ">", "type": "number"},
		"num_greater_than_or_equal": {"operator": ">=", "type": "number"},
		"num_less_than":             {"operator": "<", "type": "number"},
		"num_less_than_or_equal":    {"operator": "<=", "type": "number"},

		"date_before": {"operator": "<", "type": "string"},
		"date_after":  {"operator": ">", "type": "string"},

		"text_contains":     {"operator": "=~", "type": "string"},
		"text_not_contains": {"operator": "!~", "type": "string"},
	}
	//OperatorsLogic for Biller Config Rules
	OperatorsLogic = map[string]string{
		"and": "&&",
		"or":  "||",
	}
)

// Run ..
func (c *Condition) Run(data string) (result bool, err error) {
	oprData := Operation{
		Data:     c.Data,
		Function: c.FunctionData,
		Value:    c.Data,
	}
	c.Data, err = oprData.Run(data)
	if err != nil {
		return false, err
	}
	oprValue := Operation{
		Data:     c.Value,
		Function: c.FunctionValue,
		Value:    c.Value,
	}
	c.Value, err = oprValue.Run(data)
	if err != nil {
		return false, err
	}

	operator, oprtExist := Operators[c.Operator]
	if oprtExist {
		switch c.Operator {
		case "empty":
			return isEmpty(c.Data), nil
		case "not_empty":
			return !isEmpty(c.Data), nil
		case "start_with":
		case "end_with":
		case "in":
			return inArray(c.Data, c.Value), nil
		case "not_in":
			return !inArray(c.Data, c.Value), nil
		case "set":
			return strings.HasPrefix(c.Data, c.Value), nil
		case "not_set":
			return strings.HasSuffix(c.Data, c.Value), nil
		default:
			exp := ""
			if c.Operator == "equal" || c.Operator == "not_equal" || operator["type"] == "string" {
				c.Data = strings.ReplaceAll(c.Data, `'`, `\'`)
				c.Value = strings.ReplaceAll(c.Value, `'`, `\'`)
				exp = fmt.Sprintf("'%s' %s '%s'", c.Data, operator["operator"], c.Value)
			} else {
				exp = fmt.Sprintf("%s %s %s", c.Data, operator["operator"], c.Value)
			}
			expression, err := govaluate.NewEvaluableExpression(exp)
			if err != nil {
				return false, err
			}
			exRes, err := expression.Evaluate(nil)
			if err != nil {
				return false, err
			}
			return exRes.(bool), nil
		}
	} else {
		return false, fmt.Errorf("Operator not found")
	}

	return true, nil
}
