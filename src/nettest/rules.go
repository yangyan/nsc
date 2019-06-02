package nettest

type InstallInfo struct {
	Software string `json:"software"`
	IP       string `json:"ip"`
}

type DependInfo struct {
	SoftwareFrom string `json:"from"`
	SoftwareTo   string `json:"to"`
	Port         int    `json:"port"`
}

type Rules struct {
	InstallInfos []InstallInfo
	DependInfos  []DependInfo
}

var TestRules Rules
var TestRulesJson string
