package sys

import (
	"net"
	"regexp"
	"useful-tools/common/config"

	"github.com/pkg/errors"
)

func GetInternalIps() (internalIps []string, err error) {
	iFaces, err := net.Interfaces()
	if err != nil {
		err = errors.Wrap(err, "net.Interfaces()")
		return
	}

	for _, face := range iFaces {
		if (face.Flags & net.FlagUp) == 0 { // 过滤未开启的网卡
			continue
		}

		addresses, err := face.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}

			if matchString, err := regexp.MatchString(config.GetRegIpRule(), ip.String()); err == nil && matchString {
				internalIps = append(internalIps, ip.String())
			}
		}
	}
	return
}
