package proxy

import "useful-tools/module/walk_ui/common"

func getProxyType(name string) int {
	for _, info := range proxyType() {
		if info.Name == name {
			return info.Key
		}
	}
	return 1
}

func proxyType() []*common.CompanyItem {
	return []*common.CompanyItem{
		{Key: 1, Name: "SOCKS5"},
		{Key: 2, Name: "SSL"},
		{Key: 3, Name: "SSH"},
		{Key: 4, Name: "HTTP"},
		{Key: 5, Name: "HTTPS"},
		{Key: 6, Name: "SHADOWSOCKS"},
	}
}

func getRequestInfo(_default, new string) string {
	if new == "" {
		return _default
	}
	return new
}
