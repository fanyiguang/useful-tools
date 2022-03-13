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
	}

	RegIpRule = `(((\d\b|[1-9]\d\b|1\d\d\b|2[0-4]\d|25[0-5]).){3}(\d\b|[1-9]\d\b|1\d\d\b|2[0-4]\d|25[0-5]))`
)
