package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperation(t *testing.T) {
	data := `{
		"data1": "abcdef"
	}`

	operation := Operation{
		Data:     "[data_out]",
		Function: "text",
		Value:    "[data1]",
	}
	res, err := operation.Run(data)
	assert.Nil(t, err)
	assert.Equal(t, "abcdef", res)

	// function not text
	operation = Operation{
		Data:     "[data_out]",
		Function: "substring",
		Value:    "[data1],0,3",
	}
	res, err = operation.Run(data)
	assert.Nil(t, err)
	assert.Equal(t, "abc", res)

	// function error
	operation = Operation{
		Data:     "[data_out]",
		Function: "substring",
		Value:    "[data1]",
	}
	res, err = operation.Run(data)
	assert.Equal(t, "substring: second argument can't be empty", err.Error())
	assert.Equal(t, "", res)
}
