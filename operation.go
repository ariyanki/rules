package rules

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type (
	//Operation ..
	Operation struct {
		Data     string `json:"data"`
		Function string `json:"function"`
		Value    string `json:"value"`
	}
)

// Run ..
func (o *Operation) Run(data string) (result string, err error) {
	o.Data = strings.TrimSpace(o.Data)

	if o.Function == "text" {
		mapValues := getAllPath(o.Value)
		for _, key := range mapValues {
			exMapRes, err := mapping(data, key)
			if err != nil {
				return result, fmt.Errorf("operation %s", err.Error())
			}
			o.Value = strings.ReplaceAll(o.Value, "["+key+"]", exMapRes)
		}
		result = o.Value
	} else {
		args := strings.Split(o.Value, ",")
		var funcArgs []string
		for _, value := range args {
			mapValues := getAllPath(value)
			if len(mapValues) > 0 {
				for _, key := range mapValues {
					exMapRes, err := mapping(data, key)
					if err != nil {
						return result, fmt.Errorf("operation %s", err.Error())
					}
					value = strings.ReplaceAll(value, "["+key+"]", exMapRes)
				}
			}

			funcArgs = append(funcArgs, strings.TrimSpace(value))
		}
		f := Function{
			FuncName: o.Function,
			Args:     funcArgs,
		}
		funcRes, err := f.Run()
		if err != nil {
			return "", fmt.Errorf("operation %s", err.Error())
		}
		result = funcRes
	}

	return result, nil
}

func mapping(json string, path string) (result string, err error) {
	value := gjson.Get(json, path)
	if !value.Exists() {
		return "", fmt.Errorf("value doesn't exist: [%s]", path)
	}
	return value.String(), nil
}
