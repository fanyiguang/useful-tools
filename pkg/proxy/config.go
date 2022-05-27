package proxy

const (
	SSH    = "ssh"
	SSL    = "ssl"
	SOCKS5 = "socks5"
	HTTP   = "http"
	HTTPS  = "https"
)

var (
	CheckIpUrls = []string{
		"http://lumtest.com/myip.json",
		"https://checkip.amazonaws.com",
		"http://myipip.net/",
		"https://ipinfo.io/json",
	}
)
