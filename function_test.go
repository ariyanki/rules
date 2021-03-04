package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionNotFound(t *testing.T) {
	funcArgs := []string{"abc|cde", "|"}

	f := Function{
		FuncName: "notfound",
		Args:     funcArgs,
	}
	_, err := f.Run()
	assert.Equal(t, "Function doesn't exist: notfound", err.Error())
}

func TestFunctionSplit(t *testing.T) {
	funcArgs := []string{"abc|cde", "|"}

	f := Function{
		FuncName: "split",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "[\"abc\",\"cde\"]", funcRes)
}

func TestFunctionSubString(t *testing.T) {
	funcArgs := []string{"abcdefgh", "0", "6"}

	f := Function{
		FuncName: "substring",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "abcdef", funcRes)
}

func TestFunctionFormula(t *testing.T) {
	funcArgs := []string{"20", "/", "4"}

	f := Function{
		FuncName: "formula",
		Args:     funcArgs,
	}
	funcRes, _ := f.Run()
	assert.Equal(t, "5", funcRes)
}
