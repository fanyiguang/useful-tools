package tcp_udp

import (
	"encoding/json"
	sysNet "net"
	"strings"
	"useful-tools/helper/net"
	"useful-tools/helper/str"
	"useful-tools/helper/sys"
	"useful-tools/module/logic/base"
	"useful-tools/pkg/wlog"

	"github.com/pkg/errors"
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
	return &TcpUdp{}
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
		if trimInfo[1] == "随机" {
			isSuccess, err = t.Dial(trimInfo[0], trimInfo[2], trimInfo[3])
		} else {
			isSuccess, err = net.DialTcpByIFace(trimInfo[2], trimInfo[3], trimInfo[1])
		}
	case "UDP":
		if trimInfo[1] == "随机" {
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
	netInfos := t.parserConvenientModeContent(strings.TrimSpace(netInfo))
	if len(netInfos) < 2 {
		err = errors.Wrap(errors.New("t.parserConvenientModeContent failed"), "params too little")
		return
	}

	if len(netInfos) == 2 {
		return t.NormalDial("tcp", "随机", netInfos[0], netInfos[1])
	}

	if len(netInfos) == 3 {
		return t.NormalDial(netInfos[0], "随机", netInfos[1], netInfos[2])
	}

	return t.NormalDial(netInfos[0], netInfos[1], netInfos[2], netInfos[3])
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

func (t *TcpUdp) parserConvenientModeContent(netInfo string) (splitInfo []string) {
	var tempInfo map[string]string
	err := json.Unmarshal([]byte(netInfo), &tempInfo)
	if err == nil {
		var ip, port, network, iFace string
		for key, val := range tempInfo {
			if len(val) == 0 {
				continue
			}

			key = strings.ToLower(key)
			if strings.Contains(key, "network") {
				network = val
				continue
			}

			if strings.Contains(key, "iFace") || strings.Contains(key, "interface") {
				iFace = val
				continue
			}

			if strings.Contains(key, "port") {
				port = val
				continue
			}

			if strings.Contains(key, "ip") || strings.Contains(key, "addr") {
				ip = val
				continue
			}

			if ip != "" && port != "" {
				if network != "" {
					splitInfo = append(splitInfo, network)
				}
				if iFace != "" {
					splitInfo = append(splitInfo, iFace)
				}
				splitInfo = append(splitInfo, ip, port)
				return
			}
		}
	}

	splitInfo = strings.Split(netInfo, ":")
	if len(splitInfo) >= 2 {
		return
	}

	splitInfo = strings.Split(netInfo, " ")
	return
}
