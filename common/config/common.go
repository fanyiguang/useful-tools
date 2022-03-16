package config

var (
	regIpRule = `(((\d\b|[1-9]\d\b|1\d\d\b|2[0-4]\d|25[0-5]).){3}(\d\b|[1-9]\d\b|1\d\d\b|2[0-4]\d|25[0-5]))`
)

func GetRegIpRule() string {
	return regIpRule
}
