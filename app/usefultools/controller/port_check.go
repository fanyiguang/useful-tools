package controller

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	sysNet "net"
	"strings"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/app/usefultools/model"
	"useful-tools/helper/net"
	"useful-tools/helper/str"
	"useful-tools/helper/sys"
)

var _ adapter.Controller = (*PortCheck)(nil)

type PortCheck struct {
	Base
	network      string
	iFace        string
	host         string
	port         string
	viewText     string
	preModeInput string
}

func (p *PortCheck) Network() string {
	return p.network
}

func (p *PortCheck) SetNetwork(network string) {
	p.network = network
}

func (p *PortCheck) IFace() string {
	return p.iFace
}

func (p *PortCheck) SetIFace(iFace string) {
	p.iFace = iFace
}

func (p *PortCheck) Host() string {
	return p.host
}

func (p *PortCheck) SetHost(host string) {
	p.host = host
}

func (p *PortCheck) Port() string {
	return p.port
}

func (p *PortCheck) SetPort(port string) {
	p.port = port
}

func (p *PortCheck) ViewText() string {
	return p.viewText
}

func (p *PortCheck) SetViewText(viewText string) {
	p.viewText = viewText
}

func (p *PortCheck) PreModeInput() string {
	return p.preModeInput
}

func (p *PortCheck) PreTemplate() string {
	return ProPortCheckTemplate()
}

func (p *PortCheck) SetPreModeInput(preModeInput string) {
	p.preModeInput = preModeInput
}

func NewPortCheck() *PortCheck {
	return &PortCheck{}
}

func (p *PortCheck) NetworkList() []string {
	return []string{"TCP", "UDP"}
}

func (p *PortCheck) GetInterfaceList() []string {
	ips, err := sys.GetInternalIps()
	if err != nil {
		logrus.Warnf("get internal ips error: %v", err)
		return []string{i18n.T(i18n.KeyAuto)}
	}

	return append([]string{i18n.T(i18n.KeyAuto)}, ips...)
}

func (p *PortCheck) NormalDial(network, iFace, targetIp, targetPort string) (isSuccess bool, err error) {
	logrus.Infof("network: %v, iFace: %v, targetIp: %v, targetPort: %v", network, iFace, targetIp, targetPort)
	trimInfo := str.TrimStringSpace(network, iFace, targetIp, targetPort)
	switch trimInfo[0] {
	case "TCP":
		if i18n.Matches(i18n.KeyAuto, trimInfo[1]) {
			isSuccess, err = p.Dial(trimInfo[0], trimInfo[2], trimInfo[3])
		} else {
			isSuccess, err = net.DialTcpByIFace(trimInfo[2], trimInfo[3], trimInfo[1])
		}
	case "UDP":
		if i18n.Matches(i18n.KeyAuto, trimInfo[1]) {
			isSuccess, err = p.Dial(trimInfo[0], trimInfo[2], trimInfo[3])
		} else {
			isSuccess, err = net.DialUdpByIFace(trimInfo[2], trimInfo[3], trimInfo[1])
		}
	}
	return
}

func (p *PortCheck) ProDial(netInfo string) (isSuccess bool, err error) {
	dialInfo, err := p.parserConvenientModeContent(strings.TrimSpace(netInfo))
	if err != nil {
		return false, err
	}
	return p.NormalDial(dialInfo.Network, dialInfo.LocalIp, dialInfo.Host, dialInfo.Port)
}

func (p *PortCheck) parserConvenientModeContent(content string) (dialInfo model.DialInfo, err error) {
	err = json.Unmarshal([]byte(content), &dialInfo)
	return dialInfo, err
}

func (p *PortCheck) Dial(network, ip, port string) (isSuccess bool, err error) {
	dial, err := sysNet.Dial(strings.ToLower(network), sysNet.JoinHostPort(ip, port))
	if err != nil {
		return false, err
	}
	_ = dial.Close()
	isSuccess = true
	return
}

func (p *PortCheck) ClearCache() {
	p.SetNetwork("")
	p.SetIFace("")
	p.SetHost("")
	p.SetPort("")
	p.SetViewText("")
	p.SetPreModeInput("")
}
