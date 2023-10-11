package net

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func SendDnsRequest(url, dnsServer string) (ips []string, err error) {
	var msg *dns.Msg
	client := dns.Client{Timeout: 10 * time.Second}
	dnsMsg := dns.Msg{}
	dnsMsg.SetQuestion(fmt.Sprintf("%s.", url), dns.TypeA)
	if strings.Contains(dnsServer, ":") {
		msg, _, err = client.Exchange(&dnsMsg, dnsServer)
	} else {
		msg, _, err = client.Exchange(&dnsMsg, fmt.Sprintf("%s:53", dnsServer))
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
