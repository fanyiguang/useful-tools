package controller

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strings"
	"useful-tools/helper/str"
	"useful-tools/pkg/proxy"
)

type ProxyCheck struct {
	Base
	typ          string
	host         string
	port         string
	username     string
	password     string
	urls         string
	viewText     string
	preModeInput string
}

func (p *ProxyCheck) ViewText() string {
	return p.viewText
}

func (p *ProxyCheck) ProTemplate() string {
	return ProProxyCheckTemplate()
}

func (p *ProxyCheck) SetViewText(viewText string) {
	p.viewText = viewText
}

func (p *ProxyCheck) PreModeInput() string {
	return p.preModeInput
}

func (p *ProxyCheck) SetPreModeInput(preModeInput string) {
	p.preModeInput = preModeInput
}

func (p *ProxyCheck) Typ() string {
	return p.typ
}

func (p *ProxyCheck) SetTyp(typ string) {
	p.typ = typ
}

func (p *ProxyCheck) Host() string {
	return p.host
}

func (p *ProxyCheck) SetHost(host string) {
	p.host = host
}

func (p *ProxyCheck) Port() string {
	return p.port
}

func (p *ProxyCheck) SetPort(port string) {
	p.port = port
}

func (p *ProxyCheck) Username() string {
	return p.username
}

func (p *ProxyCheck) SetUsername(username string) {
	p.username = username
}

func (p *ProxyCheck) Password() string {
	return p.password
}

func (p *ProxyCheck) SetPassword(password string) {
	p.password = password
}

func (p *ProxyCheck) Urls() string {
	return p.urls
}

func (p *ProxyCheck) SetUrls(urls string) {
	p.urls = urls
}

func NewProxyCheck() *ProxyCheck {
	return &ProxyCheck{}
}

func (p *ProxyCheck) NormalCheckProxy(ip, port, username, password, typ, urls string, hideBody bool) (string, error) {
	reqInfo, err := buildRequestInfo(str.TrimStringSpace(typ, username, password, ip, port, urls))
	if err != nil {
		return "", err
	}
	reqInfo.HiddenBody = hideBody
	return proxy.SendHttpRequestByProxy(reqInfo)
}

func (p *ProxyCheck) ProCheckProxy(content string, hideBody bool) (string, error) {
	reqInfo, err := p.parserConvenientModeContent(strings.TrimSpace(content))
	if err != nil {
		return "", err
	}
	reqInfo.HiddenBody = hideBody
	logrus.Infof("check urls: %s", reqInfo.Request.Urls)
	return proxy.SendHttpRequestByProxy(reqInfo)
}

func (p *ProxyCheck) parserConvenientModeContent(content string) (reqInfo proxy.RequestInfo, err error) {
	err = json.Unmarshal([]byte(content), &reqInfo)
	if err != nil {
		return reqInfo, err
	}
	return
}

func (p *ProxyCheck) SupportProxyTypeList() []string {
	return []string{"SOCKS5", "SSL", "SSH", "HTTP", "HTTPS", "SS"}
}
