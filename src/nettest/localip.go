package nettest

import (
	"net"
	"fmt"
	"github.com/sirupsen/logrus"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		logrus.Error("无法获取网卡信息")
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("无法获得一个有效地址")
}
