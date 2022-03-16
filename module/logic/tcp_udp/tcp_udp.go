package tcp_udp

import (
	sysNet "net"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/helper/sys"
	"useful-tools/pkg/wlog"

	"github.com/pkg/errors"
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
			isSuccess, err = t.Dial(network, targetIp, targetPort)
		} else {
			isSuccess, err = net.DialTcpByIFace(targetIp, targetPort, iFace)
		}
	case "UDP":
		if iFace == "随机" {
			isSuccess, err = t.Dial(network, targetIp, targetPort)
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
	netInfos := strings.Split(netInfo, ":")
	if len(netInfos) < 2 {
		err = errors.Wrap(errors.New("Parameter exception"), "")
		return
	}

	if len(netInfos) == 2 {
		return t.NormalDial("tcp", "随机", netInfos[0], netInfos[1])
	}

	return t.NormalDial(netInfos[0], "随机", netInfos[1], netInfos[2])
}

func (t *TcpUdp) GetIFaceList() (list []string) {
	ips, err := sys.GetInternalIps()
	if err != nil {
		wlog.Warm("sys.GetInternalIps() failed: %v", err)
		return
	}

	list = ips
	return
}
