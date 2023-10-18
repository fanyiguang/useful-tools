package proxy

import (
	"github.com/pkg/errors"
	"strings"
	"useful-tools/pkg/proxy"
	"useful-tools/pkg/wlog"
)

func proTemplate() string {
	return "{\r\n" +
		"    \"proxy\": {\r\n" +
		"        \"type\": \"\",\r\n" +
		"        \"host\": \"\",\r\n" +
		"        \"port\": \"\",\r\n" +
		"        \"username\": \"\",\r\n" +
		"        \"password\": \"\"\r\n" +
		"    },\r\n" +
		"    \"request\": {\r\n" +
		"        \"method\": \"GET\",\r\n" +
		"        \"urls\": [\r\n" +
		"            \"https://www.baidu.com\"\r\n" +
		"        ],\r\n" +
		"        \"header\": {\r\n" +
		"            \"Host\": [\r\n" +
		"                \"www.baidu.com\"\r\n" +
		"            ],\r\n" +
		"            \"Accept-Encoding\": [\r\n" +
		"                \"zh-CN,zh\",\r\n" +
		"                \"q=0.9,ko\"\r\n" +
		"            ]\r\n" +
		"        },\r\n" +
		"        \"body\": \"\"\r\n" +
		"    },\r\n" +
		"    \"timeout\": 10\r\n" +
		"}"
}

func buildRequestInfo(proxyInfo []string) (proxy.RequestInfo, error) {
	wlog.Info("proxyInfo: %v", proxyInfo)
	if len(proxyInfo) != 6 {
		return proxy.RequestInfo{}, errors.New("proxyInfo params number neq 5")
	}
	if proxyInfo[5] == "" {
		proxyInfo[5] = proxy.CheckIpUrls
	}
	// proxyType, username, password, ip, port, proxyReqUrls
	return proxy.RequestInfo{
		Proxy: struct {
			Type     string `json:"type"`
			Host     string `json:"host"`
			Port     string `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
		}{Type: proxyInfo[0], Host: proxyInfo[3], Port: proxyInfo[4], Username: proxyInfo[1], Password: proxyInfo[2]},
		Request: struct {
			Method string              `json:"method"`
			Urls   []string            `json:"urls"`
			Header map[string][]string `json:"header"`
			Body   string              `json:"body"`
		}{
			Method: "GET",
			Urls:   strings.Split(proxyInfo[5], ";"),
			Header: make(map[string][]string),
			Body:   "",
		},
		Timeout: 15,
	}, nil
}
