package utils

import (
	"net"
	"regexp"
)

var ipv4Regx = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

var ipv6Regx = regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|` +
	`([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|` +
	`([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|` +
	`([0-9a-fA-F]{1,4}:){1,4}|` +
	`(:[0-9a-fA-F]{1,4}){1,3}|` +
	`([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|` +
	`([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|` +
	`[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|` +
	`:((:[0-9a-fA-F]{1,4}){1,7}|:)|` +
	`fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|` +
	`::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])` +
	`|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)

func FindIPv4(data string) net.IP {
	ipStr := ipv4Regx.FindString(data)
	return net.ParseIP(ipStr)
}

func FindIPv6(data string) net.IP {
	ipStr := ipv6Regx.FindString(data)
	return net.ParseIP(ipStr)
}

func FindIP(data string) net.IP {
	ipv4 := FindIPv4(data)
	if ipv4 != nil {
		return ipv4
	}

	return FindIPv6(data)
}
