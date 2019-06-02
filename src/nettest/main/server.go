package main

import (
	"net/http"
	"nettest"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/bwmarrin/snowflake"
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

func StartServer() {
	http.HandleFunc("/rules", RuleHandler)
	http.HandleFunc("/task", TaskHandler)
	http.HandleFunc("/remoteip", RemoteIPHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

var taskID string

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

	// 循环获取命令
	for {
		fmt.Printf("1: 开始新检测的任务\n" +
			"0: 退出\n")
		var opt int
		fmt.Scanln(&opt)
		switch opt {
		case 1:
			taskID = node.Generate().String()
			logrus.Infof("生成新的任务ID：%s", taskID)
		case 0:
			goto exit
		}
	}
exit:
}
