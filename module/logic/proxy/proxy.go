package proxy

import (
	"regexp"
	"strings"
	"useful-tools/helper/str"
	"useful-tools/module/logic/common"
	"useful-tools/pkg/proxy"
)

type Proxy struct {
	paramsNumber     int
	executing        bool
	proxyInfoRegRule string
}

func New() *Proxy {
	return &Proxy{
		paramsNumber:     5,
		proxyInfoRegRule: `^(.+)://(.*):(.*)@(.+):(.+)$`,
	}
}

func (p *Proxy) IsExecuting() bool {
	return p.executing
}

func (p *Proxy) SetExecuting() {
	p.executing = true
}

func (p *Proxy) ResetExecuting() {
	p.executing = false
}

func (p *Proxy) NormalCheckProxy(ip, port, username, password, proxyType string) (content string, err error) {
	if p.IsExecuting() {
		err = common.ExecutingError
		return
	}

	p.SetExecuting()
	defer p.ResetExecuting()
	return proxy.SendHttpRequestByProxy(str.TrimStringSpace(proxyType, username, password, ip, port)...)
}

func (p *Proxy) ConvenientCheckProxy(convenientModeContent string) (content string, err error) {
	if p.IsExecuting() {
		err = common.ExecutingError
		return
	}

	p.SetExecuting()
	defer p.ResetExecuting()
	return proxy.SendHttpRequestByProxy(p.parserConvenientModeContent(convenientModeContent)...)
}

func (p *Proxy) parserConvenientModeContent(content string) (proxyInfo []string) {
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
