package proxy

import (
	"encoding/json"
	"fmt"
	"strings"
	"useful-tools/helper/str"
	"useful-tools/module/logic/base"
	"useful-tools/pkg/proxy"
)

type Proxy struct {
	base.Base
	host        string
	port        string
	username    string
	password    string
	typ         string
	urls        string
	viewContent string
}

func New() *Proxy {
	p := new(Proxy)
	p.SetProTemplate(proTemplate())
	return p
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

func (p *Proxy) NormalCheckProxy(ip, port, username, password, proxyType, proxyReqUrls string, hiddenBody bool) (content string, err error) {
	reqInfo, err := buildRequestInfo(str.TrimStringSpace(proxyType, username, password, ip, port, proxyReqUrls))
	if err != nil {
		return "", err
	}
	reqInfo.HiddenBody = hiddenBody
	return proxy.SendHttpRequestByProxy(reqInfo)
}

func (p *Proxy) ConvenientCheckProxy(convenientModeContent string, hiddenBody bool) (content string, err error) {
	reqInfo, err := p.parserConvenientModeContent(strings.TrimSpace(convenientModeContent))
	if err != nil {
		return "", err
	}
	reqInfo.HiddenBody = hiddenBody
	return proxy.SendHttpRequestByProxy(reqInfo)
}

func (p *Proxy) parserConvenientModeContent(content string) (reqInfo proxy.RequestInfo, err error) {
	err = json.Unmarshal([]byte(content), &reqInfo)
	if err != nil {
		err = fmt.Errorf("JSON解析失败：%v", err.Error())
	}
	return
}
