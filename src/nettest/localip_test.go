package nettest

import "testing"

func TestGetLocalIP(t *testing.T) {
	localip, err := GetLocalIP()
	if err == nil {
		t.Logf("获取本地地址: %s", localip)
	} else {
		t.Logf("错误值: %s", err.Error())
	}
}
