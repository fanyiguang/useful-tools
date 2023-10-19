package net

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type DnsInfo struct {
	Server  string `json:"server"`
	Domain  string `json:"domain"`
	Qtype   string `json:"qtype"`
	Timeout int    `json:"timeout"`
}

func SendDnsRequest(info DnsInfo) (ips []string, err error) {
	var msg *dns.Msg
	client := dns.Client{Timeout: time.Duration(info.Timeout) * time.Second}
	dnsMsg := dns.Msg{}
	switch strings.ToLower(info.Qtype) {
	case "ipv6":
		dnsMsg.SetQuestion(fmt.Sprintf("%s.", info.Domain), dns.TypeAAAA)
	case "ipv4":
		dnsMsg.SetQuestion(fmt.Sprintf("%s.", info.Domain), dns.TypeA)
	default:
		dnsMsg.SetQuestion(fmt.Sprintf("%s.", info.Domain), dns.TypeA)
	}
	if strings.Contains(info.Server, ":") {
		msg, _, err = client.Exchange(&dnsMsg, info.Server)
	} else {
		msg, _, err = client.Exchange(&dnsMsg, fmt.Sprintf("%s:53", info.Server))
	}
	if err != nil {
		return
	}
	for _, d := range msg.Answer {
		if d.Header().Rrtype == dns.TypeA {
			a := d.(*dns.A)
			ips = append(ips, a.A.String())
		}
		if d.Header().Rrtype == dns.TypeAAAA {
			a := d.(*dns.AAAA)
			ips = append(ips, a.AAAA.String())
		}
	}
	return
}
