package nettest

import (
	"github.com/ziutek/telnet"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func TestConnection(ip string, port int) bool {
	addr := fmt.Sprintf("%s:%d", ip, port)
	logrus.Infof("开始检测%s", addr)
	conn, err := telnet.DialTimeout("tcp", addr, time.Second)
	logrus.Infof("结束检测%s", addr)
	if err == nil {
		conn.Close()
		return true
	} else {
		return false
	}
}
