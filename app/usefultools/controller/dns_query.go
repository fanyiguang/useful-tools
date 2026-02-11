package controller

import (
	"encoding/json"
	sysNet "net"
	"strings"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/helper/net"
	"useful-tools/helper/str"
)

var _ adapter.Controller = (*DnsQuery)(nil)

type DnsQuery struct {
	Base
	host         string
	domain       string
	viewText     string
	preModeInput string
}

func NewDnsQuery() *DnsQuery {
	return &DnsQuery{}

}

func (d *DnsQuery) Host() string {
	return d.host
}

func (d *DnsQuery) SetHost(host string) {
	d.host = host
}

func (d *DnsQuery) Domain() string {
	return d.domain
}

func (d *DnsQuery) SetDomain(domain string) {
	d.domain = domain
}

func (d *DnsQuery) ViewText() string {
	return d.viewText
}

func (d *DnsQuery) SetViewText(viewText string) {
	d.viewText = viewText
}

func (d *DnsQuery) PreModeInput() string {
	return d.preModeInput
}

func (d *DnsQuery) PreTemplate() string {
	return ProDnsQueryTemplate()
}

func (d *DnsQuery) SetPreModeInput(preModeInput string) {
	d.preModeInput = preModeInput
}

func (d *DnsQuery) NormalDns(dnsServer, domain string) (ips []string, err error) {
	dnsInfo := buildDnsInfo(str.TrimStringSpace(dnsServer, domain))
	if dnsInfo.Server == "" || i18n.Matches(i18n.KeyDefault, dnsInfo.Server) {
		return sysNet.LookupHost(domain)
	} else {
		return net.SendDnsRequest(dnsInfo)
	}
}

func (d *DnsQuery) ProQuery(content string) (ips []string, err error) {
	dnsInfo, err := d.parserConvenientModeContent(strings.TrimSpace(content))
	if err != nil {
		return nil, err
	}
	return net.SendDnsRequest(dnsInfo)
}

func (d *DnsQuery) parserConvenientModeContent(data string) (dnsInfo net.DnsInfo, err error) {
	err = json.Unmarshal([]byte(data), &dnsInfo)
	return dnsInfo, err
}

func (d *DnsQuery) ClearCache() {
	d.host = ""
	d.domain = ""
	d.viewText = ""
	d.preModeInput = ""
}
