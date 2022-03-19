package proxy

import (
	"encoding/json"
	"regexp"
	"strings"
	"useful-tools/helper/str"
	"useful-tools/module/logic/base"
	"useful-tools/pkg/proxy"
)

type Proxy struct {
	base.Base
	paramsNumber     int
	proxyInfoRegRule string
}

func New() *Proxy {
	return &Proxy{
		paramsNumber:     5,
		proxyInfoRegRule: `^(.+)://(.*):(.*)@(.+):(.+)$`,
	}
}

func (p *Proxy) NormalCheckProxy(ip, port, username, password, proxyType string) (content string, err error) {
	return proxy.SendHttpRequestByProxy(str.TrimStringSpace(proxyType, username, password, ip, port)...)
}

func (p *Proxy) ConvenientCheckProxy(convenientModeContent string) (content string, err error) {
	return proxy.SendHttpRequestByProxy(p.parserConvenientModeContent(strings.TrimSpace(convenientModeContent))...)
}

func (p *Proxy) parserConvenientModeContent(content string) (proxyInfo []string) {
	var infos = make(map[string]string)
	err := json.Unmarshal([]byte(content), &infos)
	if err == nil {
		var proxyType, username, password, ip, port string
		for key, info := range infos {
			key = strings.ToLower(key)
			if strings.Contains(key, "type") {
				proxyType = info
				continue
			}

			if strings.Contains(key, "name") {
				username = info
				continue
			}

			if strings.Contains(key, "pass") {
				password = info
				continue
			}

			if strings.Contains(key, "ip") || strings.Contains(key, "addr") {
				ip = info
				continue
			}

			if strings.Contains(key, "port") {
				port = info
				continue
			}
		}

		if proxyType != "" && username != "" && password != "" && port != "" && ip != "" {
			proxyInfo = append(proxyInfo, proxyType, username, password, ip, port)
			return
		}
	}

	compile := regexp.MustCompile(p.proxyInfoRegRule)
	subMatch := compile.FindAllStringSubmatch(content, -1)
	if len(subMatch) > 0 {
		return subMatch[0][1:]
	}
	proxyInfo = strings.Split(content, " ")
	if len(proxyInfo) == p.paramsNumber {
		return
	}

	proxyInfo = strings.Split(content, ":")
	return
}
