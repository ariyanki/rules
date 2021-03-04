package rules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
)

type (
	// Rule ..
	Rule struct {
		Data   string
		Output string
		Config string
	}
	//Config ..
	Config struct {
		Label      string      `json:"label"`
		Conditions []Condition `json:"conditions"`
		Operations []Operation `json:"operations"`
	}
)

// Run ..
func (r *Rule) Run() (result string, err error) {
	var config []Config
	json.Unmarshal([]byte(r.Config), &config)

	for _, rl := range config {
		pass, err := executeConditions(rl.Conditions, r.Data)
		if pass {
			for _, obj := range rl.Operations {
				exOpr, err := obj.Run(r.Data)
				if err != nil {
					break
				}
				r.Output = strings.ReplaceAll(r.Output, obj.Data, exOpr)
			}

		}
		if err != nil {
			break
		}
	}

	return r.Output, err
}

func executeConditions(conditions []Condition, data string) (result bool, err error) {
	exp := ""
	result = true
	for _, cond := range conditions {
		condRes, err := cond.Run(data)
		if err != nil {
			break
		}
		if len(conditions) > 1 {
			exp = fmt.Sprintf("%s %t %s", exp, condRes, OperatorsLogic[cond.Logic])
		} else {
			result = condRes
		}

	}
	if err != nil {
		return false, err
	}
	if len(conditions) > 1 {
		expression, err := govaluate.NewEvaluableExpression(exp)
		if err != nil {
			return false, err
		}
		resEx, err := expression.Evaluate(nil)
		if err != nil {
			return false, err
		}
		result = resEx.(bool)
	}
	return result, nil
}
