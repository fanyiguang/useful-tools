package proxy

import (
	"log"
	"network-tool/module/logic/common"
	"network-tool/pkg/proxy"
	"strings"
)

type Proxy struct {
	//Page *proxy.Page
	executing bool
}

func New() *Proxy {
	return new(Proxy)
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

func (p *Proxy) CheckProxy(ip, port, username, password, proxyType string) (content string, err error) {
	if p.IsExecuting() {
		err = common.ExecutingError
		return
	}

	p.SetExecuting()
	defer p.ResetExecuting()
	log.Println(ip, port, username, password, proxyType)
	return proxy.SendHttpRequestByProxy(ip, port, username, password, strings.ToLower(proxyType))
}
