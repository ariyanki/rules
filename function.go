package rules

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type (
	// Function ..
	Function struct {
		FuncName string
		Args     []string
	}
)

// Run ..
func (f *Function) Run() (string, error) {
	switch f.FuncName {
	case "split":
		if f.Args[0] == "" {
			return "", fmt.Errorf("split: invalid first argument %s", fmt.Sprint(f.Args[0]))
		}
		if f.Args[1] == "" {
			f.Args[1] = ","
		}

		res := strings.Split(f.Args[0], f.Args[1])
		result, _ := json.Marshal(res)

		return string(result), nil
	case "substring":
		if f.Args[0] == "" {
			return "", fmt.Errorf("substring: first argument can't be empty")
		}
		lenArgs := len(f.Args)
		if lenArgs < 2 {
			return "", fmt.Errorf("substring: second argument can't be empty")
		}

		first := []rune(f.Args[0]) //support for special char
		lenFirstArgs := len(first)
		second, _ := strconv.Atoi(f.Args[1])
		third := 0

		if second < 0 || second >= lenFirstArgs {
			return "", fmt.Errorf("substring: second argument index out of range")
		}
		if lenArgs > 2 {
			third, _ = strconv.Atoi(f.Args[2])
			if third < 1 || third > lenFirstArgs {
				return "", fmt.Errorf("substring: third argument index out of range")
			}
			if second == third {
				return "", fmt.Errorf("substring: second and third argument can't be the same")
			}
			if second >= third {
				return "", fmt.Errorf("substring: second argument index out of range")
			}
		}

		var result string
		if third < 1 {
			result = string(first[second:])
		} else {
			result = string(first[second:third])
		}
		return result, nil

	}
	return "", fmt.Errorf("Function doesn't exist: %s", f.FuncName)
}
