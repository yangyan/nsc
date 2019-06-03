package main

import (
	"net/http"
	"nettest"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"io/ioutil"
	"encoding/json"
	"github.com/tealeg/xlsx"
)

func RuleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(nettest.TestRulesJson))
	}
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(taskID))
	}
}

func RemoteIPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(r.RemoteAddr))
	}
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			// 获取结果
			var results nettest.Results
			err = json.Unmarshal(body, &results)
			if err == nil {
				// 获取任务ID
				fmt.Println(results.TaskID)
				// 扔到chan中处理
				resultsCh <- &results
			}
		}
	}
}

func ProcessResults() {
	for {
		// 读取结果
		results := <-resultsCh
		fileName := results.TaskID + ".xlsx"
		// 创建结果文件
		file, err := xlsx.OpenFile(fileName)
		if err != nil {
			logrus.Infof("无法读取结果文件")
			file = xlsx.NewFile()
		}
		_, ok := file.Sheet["结果"]
		if !ok {
			file.AddSheet("结果")
		}

		// 写入结果
		for _, result := range results.AllStatus {
			row := file.Sheet["结果"].AddRow()
			row.AddCell().SetValue(result.SoftwareFrom)
			row.AddCell().SetValue(result.IpFrom, )
			row.AddCell().SetValue(result.SoftwareTo)
			row.AddCell().SetValue(result.IpTo)
			row.AddCell().SetValue(result.Port)
			row.AddCell().SetValue(result.Status)
		}

		file.Save(fileName)
	}
}

func StartServer() {
	http.HandleFunc("/rules", RuleHandler)
	http.HandleFunc("/task", TaskHandler)
	http.HandleFunc("/remoteip", RemoteIPHandler)
	http.HandleFunc("/results", ResultHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

var taskID string
var resultsCh = make(chan *nettest.Results)

func main() {
	// 初始化日志
	logrus.Info("服务启动")

	// 读取规则
	nettest.ReadRules("D:\\yy-important\\OKR\\运维思路\\工程\\src\\nettest\\testdata\\rules.xlsx")
	nettest.RulesToJson()

	// 任务启动
	node, err := snowflake.NewNode(0)
	if err != nil {
		logrus.Fatal("生成任务ID失败")
	}

	// 服务启动
	go StartServer()
	// 开启记录结果线程
	go ProcessResults()

	// 循环获取命令
	for {
		fmt.Printf("1: 开始新检测的任务\n" +
			"0: 退出\n")
		var opt = -1
		fmt.Scanln(&opt)
		switch opt {
		case 1:
			taskID = node.Generate().String()
			logrus.Infof("生成新的任务ID：%s", taskID)
		case 0:
			goto exit
		default:
			continue
		}
	}
exit:
}
