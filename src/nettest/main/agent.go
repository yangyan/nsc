package main

import (
	"net/http"
	"io/ioutil"
	"github.com/sirupsen/logrus"
	"time"
	"encoding/json"
	"nettest"
	"fmt"
	"strings"
)

var url string = "http://127.0.0.1:8080"
var currentTaskID string

func GetTaskID() string {
	// 循环获取当前的任务
	resp, err := http.Get(url + "/task")
	if err != nil {
		logrus.Errorf("获取任务信息失败，错误原因: %s", err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("获取任务信息失败，错误原因: %s", err.Error())
		return ""
	}

	return string(body)
}

func GetRules() string {
	// 循环获取当前的任务
	resp, err := http.Get(url + "/rules")
	if err != nil {
		logrus.Errorf("获取规则信息失败，错误原因: %s", err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("获取规则信息失败，错误原因: %s", err.Error())
		return ""
	}

	return string(body)
}

func GetLocalIP() string {
	// 循环获取当前的任务
	resp, err := http.Get(url + "/remoteip")
	if err != nil {
		logrus.Errorf("获取IP信息失败，错误原因: %s", err.Error())
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("获取IP信息失败，错误原因: %s", err.Error())
		return ""
	}

	ipport := string(body)

	parts := strings.Split(ipport, ":")
	if len(parts) != 2 {
		logrus.Errorf("获取IP信息失败，获得数据格式正确：%s", ipport)
		return ""
	}

	return parts[0]
}

func SendResults(results string) {
	// 循环获取当前的任务
	resp, err := http.Post(url+"/results", "text/json", strings.NewReader(results))
	if err != nil {
		logrus.Errorf("上传结果出现错误：%s", err)
	}
	defer resp.Body.Close()
}


func Detect(taskid string, rules *nettest.Rules, localip string) *nettest.Results {
	// 构建软件和IP的对应表
	soft2ips := make(map[string][]string)
	for _, installinfo := range rules.InstallInfos {
		soft2ips[installinfo.Software] = append(soft2ips[installinfo.Software], installinfo.IP)
	}
	soft2depends := make(map[string][]nettest.DependInfo)
	for _, dependinfo := range rules.DependInfos {
		soft2depends[dependinfo.SoftwareFrom] = append(soft2depends[dependinfo.SoftwareFrom], dependinfo)
	}

	results := new(nettest.Results)
	results.TaskID = taskid
	for _, installinfo := range rules.InstallInfos {
		if installinfo.IP == localip {
			softs, ok := soft2depends[installinfo.Software]
			if ok {
				for _, depend := range softs {
					ips, ok := soft2ips[depend.SoftwareTo]
					if ok {
						for _, ipto := range ips {
							results.AllStatus = append(results.AllStatus, &nettest.Result{installinfo.Software, depend.SoftwareTo, localip, ipto, depend.Port, false})
						}
					}
				}
			}
		}
	}

	for _, result := range results.AllStatus {
		ok := nettest.TestConnection(result.IpTo, result.Port)
		result.Status = ok
	}

	return results
}

func main() {
	for {
		newTaskID := GetTaskID()
		logrus.Infof("获取任务ID：%s", newTaskID)

		if newTaskID != "" && newTaskID != currentTaskID {
			currentTaskID = newTaskID
			// 开始新的任务
			logrus.Info("开始信息任务")

			// 获取规则信息
			rulesJson := GetRules()
			var rules nettest.Rules
			if len(rulesJson) > 0 {
				json.Unmarshal([]byte(rulesJson), &rules)
				fmt.Println(rules)
			}

			// 获取IP信息
			localip := GetLocalIP()
			logrus.Infof("获得本机地址： %s", localip)

			// 计算任务
			results := Detect(newTaskID, &rules, localip)
			bytes, _ := json.Marshal(results)
			fmt.Print(string(bytes))

			// 上报结果
			SendResults(string(bytes))
		}

		time.Sleep(time.Second * 5)
	}
}
