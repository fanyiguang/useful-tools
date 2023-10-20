package dns

import "useful-tools/helper/net"

func buildDnsInfo(info []string) net.DnsInfo {
	return net.DnsInfo{
		Server:  info[0],
		Domain:  info[1],
		Qtype:   "",
		Timeout: 10,
	}
}

func proTemplate() string {
	return "{\r\n" +
		"    \"server\": \"114.114.114.114\",\r\n" +
		"    \"domain\": \"www.baidu.com\",\r\n" +
		"    \"timeout\": 10,\r\n" +
		"    \"qtype\": \"ipv4\"\r\n" +
		"}"
}
