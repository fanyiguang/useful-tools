package dns

import (
	"encoding/json"
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
	return &Dns{}
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
	infos := str.TrimStringSpace(dnsServer, domain)
	if infos[0] == "默认" || infos[0] == "" {
		return sysNet.LookupHost(domain)
	} else {
		return net.SendDnsRequest(infos[1], infos[0])
	}
}

func (t *Dns) ConvenientDns(dnsInfo string) (ips []string, err error) {
	dnsInfos := t.parserConvenientModeContent(strings.TrimSpace(dnsInfo))
	if len(dnsInfos) < 2 {
		return sysNet.LookupHost(dnsInfos[0])
	} else {
		return net.SendDnsRequest(dnsInfos[1], dnsInfos[0])
	}
}

func (t *Dns) parserConvenientModeContent(info string) (splitInfo []string) {
	var tempInfo = make(map[string]string)
	err := json.Unmarshal([]byte(info), &tempInfo)
	if err == nil {
		var dnsServer, domain string
		for key, val := range tempInfo {
			if len(val) == 0 {
				continue
			}
			key = strings.ToLower(key)
			if strings.Contains(key, "domain") || strings.Contains(key, "url") {
				domain = val
				continue
			}

			if strings.Contains(key, "dns") || strings.Contains(key, "server") {
				dnsServer = val
				continue
			}

			if domain != "" {
				if dnsServer != "" {
					splitInfo = append(splitInfo, dnsServer, domain)
				} else {
					splitInfo = append(splitInfo, domain)
				}
				return
			}
		}
	}

	splitInfo = strings.Split(info, ":")
	if len(splitInfo) >= 2 {
		return
	}

	splitInfo = strings.Split(info, " ")
	return
}
