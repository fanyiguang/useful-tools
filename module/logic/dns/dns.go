package dns

import (
	sysNet "net"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/helper/str"
	"useful-tools/module/logic/base"
	"useful-tools/module/logic/common"
)

type Dns struct {
	base.Base
}

func New() *Dns {
	return &Dns{}
}

func (t *Dns) NormalDns(dnsServer, domain string) (ips []string, err error) {
	if t.IsExecuting() {
		err = common.ExecutingError
		return
	}

	t.SetExecuting()
	defer t.ResetExecuting()
	infos := str.TrimStringSpace(dnsServer, domain)
	if infos[0] != "默认" {
		return net.SendDnsRequest(infos[1], infos[0])
	} else {
		return sysNet.LookupHost(domain)
	}
}

func (t *Dns) ConvenientDns(dnsInfo string) (ips []string, err error) {
	if t.IsExecuting() {
		err = common.ExecutingError
		return
	}

	t.SetExecuting()
	defer t.ResetExecuting()
	dnsInfos := strings.Split(dnsInfo, ":")
	if len(dnsInfos) < 2 {
		return sysNet.LookupHost(dnsInfos[0])
	} else {
		return net.SendDnsRequest(dnsInfos[1], dnsInfos[0])
	}
}
