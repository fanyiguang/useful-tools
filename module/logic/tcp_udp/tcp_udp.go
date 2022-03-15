package tcp_udp

import (
	sysNet "net"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/pkg/wlog"
)

type TcpUdp struct {
}

func New() *TcpUdp {
	return &TcpUdp{}
}

func (t *TcpUdp) NormalDial(network, iFace, targetIp, targetPort string) (isSuccess bool, err error) {
	wlog.Info(network, iFace, targetIp, targetPort)
	switch network {
	case "TCP":
		if iFace == "随机" {
			return t.Dial(network, targetIp, targetPort)
		} else {
			isSuccess, err = net.DialTcpByIFace(targetIp, targetPort, iFace)
		}
	case "UDP":
		if iFace == "随机" {
			return t.Dial(network, targetIp, targetPort)
		} else {
			isSuccess, err = net.DialUdpByIFace(targetIp, targetPort, iFace)
		}
	}
	return
}

func (t *TcpUdp) Dial(network, ip, port string) (isSuccess bool, err error) {
	dial, err := sysNet.Dial(strings.ToLower(network), sysNet.JoinHostPort(ip, port))
	if err != nil {
		return false, err
	}
	_ = dial.Close()
	isSuccess = true
	return
}

func (t *TcpUdp) ConvenientDial(netInfo string) (isSuccess bool, err error) {
	return
}
