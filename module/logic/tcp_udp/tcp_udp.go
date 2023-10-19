package tcp_udp

import (
	"encoding/json"
	"fmt"
	sysNet "net"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/helper/str"
	"useful-tools/helper/sys"
	"useful-tools/module/logic/base"
	"useful-tools/pkg/wlog"
)

type TcpUdp struct {
	base.Base
	network     string
	iFace       string
	host        string
	port        string
	viewContent string
}

func New() *TcpUdp {
	t := new(TcpUdp)
	t.SetProTemplate(proTemplate())
	return t
}

func (t *TcpUdp) ViewContent() string {
	return t.viewContent
}

func (t *TcpUdp) SetViewContent(viewContent string) {
	t.viewContent = viewContent
}

func (t *TcpUdp) Network() string {
	return t.network
}

func (t *TcpUdp) SetNetwork(network string) {
	t.network = network
}

func (t *TcpUdp) IFace() string {
	return t.iFace
}

func (t *TcpUdp) SetIFace(iFace string) {
	t.iFace = iFace
}

func (t *TcpUdp) Host() string {
	return t.host
}

func (t *TcpUdp) SetHost(host string) {
	t.host = host
}

func (t *TcpUdp) Port() string {
	return t.port
}

func (t *TcpUdp) SetPort(port string) {
	t.port = port
}

func (t *TcpUdp) NormalDial(network, iFace, targetIp, targetPort string) (isSuccess bool, err error) {
	wlog.Info(network, iFace, targetIp, targetPort)
	trimInfo := str.TrimStringSpace(network, iFace, targetIp, targetPort)
	switch trimInfo[0] {
	case "TCP":
		if trimInfo[1] == "自动" {
			isSuccess, err = t.Dial(trimInfo[0], trimInfo[2], trimInfo[3])
		} else {
			isSuccess, err = net.DialTcpByIFace(trimInfo[2], trimInfo[3], trimInfo[1])
		}
	case "UDP":
		if trimInfo[1] == "自动" {
			isSuccess, err = t.Dial(trimInfo[0], trimInfo[2], trimInfo[3])
		} else {
			isSuccess, err = net.DialUdpByIFace(trimInfo[2], trimInfo[3], trimInfo[1])
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
	dialInfo, err := t.parserConvenientModeContent(strings.TrimSpace(netInfo))
	if err != nil {
		return false, err
	}
	return t.NormalDial(dialInfo.Network, dialInfo.LocalIp, dialInfo.Host, dialInfo.Port)
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

func (t *TcpUdp) parserConvenientModeContent(content string) (dialInfo DialInfo, err error) {
	err = json.Unmarshal([]byte(content), &dialInfo)
	if err != nil {
		err = fmt.Errorf("JSON解析失败：%v", err.Error())
	}
	return
}
