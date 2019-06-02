package nettest

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/sirupsen/logrus"
)

func ReadRules(fileName string) {
	xlsx, err := xlsx.OpenFile(fileName)
	if err != nil {
		logrus.Fatal("无法读取规则文件：rules.xlsx")
	}

	// 清空
	TestRules.InstallInfos = TestRules.InstallInfos[0:0]
	TestRules.DependInfos = TestRules.DependInfos[0:0]

	// 读取数据
	for _, row := range xlsx.Sheets[0].Rows {
		TestRules.InstallInfos = append(TestRules.InstallInfos, InstallInfo{row.Cells[0].String(), row.Cells[1].String()})
	}

	for _, row := range xlsx.Sheets[1].Rows {
		port, _ := row.Cells[2].Int()
		TestRules.DependInfos = append(TestRules.DependInfos, DependInfo{row.Cells[0].String(), row.Cells[1].String(), port})
	}
}

func RulesToJson() {
	bytes, err := json.Marshal(TestRules)
	fmt.Printf("%d\n", len(bytes))
	if err != nil {
		logrus.Fatal("规则文件转化成JSON出现错误！")
	}
	TestRulesJson = string(bytes)
}
