package nettest

import (
	"fmt"
	"testing"
)

func TestReadRules(t *testing.T) {
	ReadRules("./testdata/rules.xlsx")
}

func TestRulesToJson(t *testing.T) {
	RulesToJson()
	fmt.Printf("%s\n", TestRulesJson)
}
