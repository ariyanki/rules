package rules

import (
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
		funcRes, err := f.Run()
		if err != nil {
			return "", err
		}
		result = funcRes
	}

	return result, nil
}
