package rules

import (
	"encoding/json"
	"strings"
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
		conds := Conditions{
			Conditions: rl.Conditions,
		}
		pass, err := conds.Runs(r.Data)
		if err != nil {
			return "", err
		}
		if pass {
			for _, operation := range rl.Operations {
				exOpr, err := operation.Run(r.Data)
				if err != nil {
					return "", err
				}
				r.Output = strings.ReplaceAll(r.Output, operation.Data, exOpr)
			}

		}
	}

	return r.Output, err
}
