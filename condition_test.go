package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCondition(t *testing.T) {
	data := `{
		"data1": "abcdef"
	}`

	c := Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "equal",
		FunctionValue: "text",
		Value:         "abcdef",
		Logic:         "AND",
	}
	condRes, err := c.Run(data)
	assert.Nil(t, err)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "substring",
		Data:          "[data1]",
		Operator:      "equal",
		FunctionValue: "text",
		Value:         "abcdef",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.Equal(t, "substring: second argument can't be empty", err.Error())
	assert.False(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "equal",
		FunctionValue: "substring",
		Value:         "12345",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.Equal(t, "substring: second argument can't be empty", err.Error())
	assert.False(t, condRes)

	empty := `{
		"data1": ""
	}`

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "empty",
		FunctionValue: "text",
		Value:         "",
		Logic:         "AND",
	}
	condRes, err = c.Run(empty)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "not_empty",
		FunctionValue: "text",
		Value:         "",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "set",
		FunctionValue: "text",
		Value:         "",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data2]",
		Operator:      "not_set",
		FunctionValue: "text",
		Value:         "",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "start_with",
		FunctionValue: "text",
		Value:         "abc",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "end_with",
		FunctionValue: "text",
		Value:         "def",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "in",
		FunctionValue: "text",
		Value:         "abcdef, hijklmn",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "[data1]",
		Operator:      "not_in",
		FunctionValue: "text",
		Value:         "abscdef, hijklmn",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "2",
		Operator:      "num_greater_than",
		FunctionValue: "text",
		Value:         "1",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.True(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "2",
		Operator:      "num_greater_than",
		FunctionValue: "text",
		Value:         "",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.Equal(t, "Unexpected end of expression", err.Error())
	assert.False(t, condRes)

	c = Condition{
		FunctionData:  "text",
		Data:          "2",
		Operator:      "not_found_func",
		FunctionValue: "text",
		Value:         "",
		Logic:         "AND",
	}
	condRes, err = c.Run(data)
	assert.Equal(t, "Operator not found", err.Error())
	assert.False(t, condRes)

}
