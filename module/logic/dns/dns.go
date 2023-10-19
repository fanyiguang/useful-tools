package dns

import (
	"encoding/json"
	"fmt"
	sysNet "net"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/helper/str"
	"useful-tools/module/logic/base"
)

type Dns struct {
	base.Base
	server      string
	domain      string
	viewContent string
}

func New() *Dns {
	d := new(Dns)
	d.SetProTemplate(proTemplate())
	return d
}

func (t *Dns) Server() string {
	return t.server
}

func (t *Dns) SetServer(server string) {
	t.server = server
}

func (t *Dns) Domain() string {
	return t.domain
}

func (t *Dns) SetDomain(domain string) {
	t.domain = domain
}

func (t *Dns) ViewContent() string {
	return t.viewContent
}

func (t *Dns) SetViewContent(viewContent string) {
	t.viewContent = viewContent
}

func (t *Dns) NormalDns(dnsServer, domain string) (ips []string, err error) {
	dnsInfo := buildDnsInfo(str.TrimStringSpace(dnsServer, domain))
	if dnsInfo.Server == "默认" || dnsInfo.Server == "" {
		return sysNet.LookupHost(domain)
	} else {
		return net.SendDnsRequest(dnsInfo)
	}
}

func (t *Dns) ConvenientDns(content string) (ips []string, err error) {
	dnsInfo, err := t.parserConvenientModeContent(strings.TrimSpace(content))
	if err != nil {
		return nil, err
	}
	return net.SendDnsRequest(dnsInfo)
}

func (t *Dns) parserConvenientModeContent(data string) (dnsInfo net.DnsInfo, err error) {
	err = json.Unmarshal([]byte(data), &dnsInfo)
	if err != nil {
		err = fmt.Errorf("JSON解析失败：%v", err.Error())
	}
	return
}
