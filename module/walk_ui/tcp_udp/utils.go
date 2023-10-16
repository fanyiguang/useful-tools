package tcp_udp

import (
	"fmt"
	"time"
	"useful-tools/module/walk_ui/common"
)

func network() []*common.CompanyItem {
	return []*common.CompanyItem{
		{Key: 1, Name: "TCP"},
		{Key: 2, Name: "UDP"},
	}
}

func defaultIFaceList() []*common.CompanyItem {
	return []*common.CompanyItem{
		{Key: 1, Name: "随机"},
	}
}

func getNetwork(net string) int {
	for _, info := range network() {
		if info.Name == net {
			return info.Key
		}
	}
	return 1
}

func createIFaceList(ips []string) []*common.CompanyItem {
	list := defaultIFaceList()
	for i, ip := range ips {
		list = append(list, &common.CompanyItem{
			Name: ip,
			Key:  i + 2,
		})
	}
	return list
}

func getInterface(i string, ips []string) int {
	for key, info := range createIFaceList(ips) {
		if info.Name == i {
			return key
		}
	}
	return 0
}

func logFormat(network, i, host, port, msg string) string {
	return fmt.Sprintf("[%v] %v:%v(%v)[%v] => %v\r\n\r\n", time.Now().Format(`01-02 15:04:05`), host, port, network, i, msg)
}
