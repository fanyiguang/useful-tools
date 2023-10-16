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
	host             string
	port             string
	username         string
	password         string
	typ              string
	urls             string
	viewContent      string
	proxyInfoRegRule string
	paramsNumber     int
}

func New() *Proxy {
	return &Proxy{
		paramsNumber:     5,
		proxyInfoRegRule: `^(.+)://(.*):(.*)@(.+):(.+)$`,
	}
}

func (p *Proxy) ViewContent() string {
	return p.viewContent
}

func (p *Proxy) SetViewContent(viewContent string) {
	p.viewContent = viewContent
}

func (p *Proxy) SetHost(host string) {
	p.host = host
}

func (p *Proxy) SetPort(port string) {
	p.port = port
}

func (p *Proxy) SetUsername(username string) {
	p.username = username
}

func (p *Proxy) SetPassword(password string) {
	p.password = password
}

func (p *Proxy) SetTyp(typ string) {
	p.typ = typ
}

func (p *Proxy) SetUrls(urls string) {
	p.urls = urls
}

func (p *Proxy) Host() string {
	return p.host
}

func (p *Proxy) Port() string {
	return p.port
}

func (p *Proxy) Username() string {
	return p.username
}

func (p *Proxy) Password() string {
	return p.password
}

func (p *Proxy) Typ() string {
	return p.typ
}

func (p *Proxy) Urls() string {
	return p.urls
}

func (p *Proxy) NormalCheckProxy(ip, port, username, password, proxyType, proxyReqUrls string) (content string, err error) {
	return proxy.SendHttpRequestByProxy(str.TrimStringSpace(proxyType, username, password, ip, port, proxyReqUrls)...)
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

		if proxyType != "" && port != "" && ip != "" {
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
