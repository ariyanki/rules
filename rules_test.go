package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {
	config := `[
		{
			"conditions": [
				{
					"data": "[data1]",
					"function_data": "text",
					"function_value": "substring",
					"logic": "and",
					"operator": "equal",
					"value": "abcdefghij,0,6"
				}
			],
			"label": "Sample",
			"operations": [
				{
					"data": "[text]",
					"function": "text",
					"value": "{[data2.dt1]:[data2.dt2]}"
				},
				{
					"data": "[split]",
					"function": "split",
					"value": "[data2.dt1],|"
				}
			]
		}
	]`
	data := `{
		"data1": "abcdef",
		"data2": {
			"dt1": "a|bc",
			"dt2": 5
		}
	}`
	output := `{
		"text":"[text]",
		"split":"[split]"
	}`
	rule := Rule{
		Data:   data,
		Output: output,
		Config: config,
	}
	result, err := rule.Run()

	assertData := `{
		"text":"{a|bc:5}",
		"split":"["a","bc"]"
	}`
	assert.Equal(t, assertData, result)
	assert.Equal(t, nil, err)
}
