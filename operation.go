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
	errMsg := ""
	o.Data = strings.TrimSpace(o.Data)
	if o.Data == "" {
		errMsg = "Operations Data is empty"
	}

	if o.Function == "text" {
		mapValues := getAllPath(o.Value)
		for _, key := range mapValues {
			exMapRes := gjson.Get(data, key)
			o.Value = strings.ReplaceAll(o.Value, "["+key+"]", exMapRes.String())
		}
		result = o.Value
	} else {
		args := strings.Split(o.Value, ",")
		var funcArgs []string
		for _, value := range args {
			mapValues := getAllPath(value)
			if len(mapValues) > 0 {
				for _, key := range mapValues {
					exMapRes := gjson.Get(data, key)
					value = strings.ReplaceAll(value, "["+key+"]", exMapRes.String())
				}
			}

			funcArgs = append(funcArgs, strings.TrimSpace(value))
		}
		f := Function{
			FuncName: o.Function,
			Args:     funcArgs,
		}
		funcRes, _ := f.Run()
		result = funcRes
	}

	if errMsg != "" {
		return "", fmt.Errorf(errMsg)
	}

	return result, nil
}
