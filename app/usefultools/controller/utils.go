package controller

import (
	"github.com/pkg/errors"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/pkg/proxy"
)

func buildRequestInfo(proxyInfo []string) (proxy.RequestInfo, error) {
	if len(proxyInfo) != 6 {
		return proxy.RequestInfo{}, errors.New("proxy info error")
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
			Urls:   strings.Split(proxyInfo[5], "\n"),
			Header: make(map[string][]string),
			Body:   "",
		},
		Timeout: 15,
	}, nil
}

func buildDnsInfo(info []string) net.DnsInfo {
	return net.DnsInfo{
		Server:  info[0],
		Domain:  info[1],
		Qtype:   "",
		Timeout: 10,
	}
}
