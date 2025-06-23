package page

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
	"useful-tools/helper/Go"
	"useful-tools/utils"
)

var _ adapter.Page = (*PortCheck)(nil)

type PortCheck struct {
	BasePage
	logics       *controller.PortCheck
	view         *widget.Entry
	scroll       *container.Scroll
	latestParams string
}

func NewPortCheck() *PortCheck {
	return &PortCheck{
		BasePage: BasePage{
			Title:      "端口检测",
			Intro:      "TCP/UDP端口检测",
			SupportWeb: true,
		},
		logics: controller.NewPortCheck(),
	}
}

func (p *PortCheck) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(p.leftScreen(mode), p.rightScreen(w))
}

func (p *PortCheck) leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("port check mode: %d", mode)
	switch mode {
	case constant.ViewModePro:
		return p.proView()
	default:
		right := viewWidget.MakeCellSize(10, 10)
		left := viewWidget.MakeCellSize(10, 10)
		top := viewWidget.MakeCellSize(10, 10)
		bottom := viewWidget.MakeCellSize(10, 10)
		return container.NewBorder(top, bottom, left, right, p.portCheckFrom())
	}
}

func (p *PortCheck) proView() fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	multi := widget.NewMultiLineEntryEx(nil, nil, nil, p.logics.FormatJson)
	multi.Wrapping = fyne.TextWrapWord
	if p.logics.PreModeInput() != "" {
		multi.SetText(p.logics.PreModeInput())
	} else {
		multi.SetText(p.logics.PreTemplate())
	}
	multi.OnChanged = func(s string) {
		logrus.Infof("port: %s", s)
		p.logics.SetPreModeInput(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("pre clear port: %s", multi.Text)
			multi.SetText("")
		},
	}, &widget.Button{
		Text:       "检测",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			p.latestParams = multi.Text
			go func() {
				text := multi.Text
				logrus.Infof("pre check port: %s", text)
				response, err := p.logics.ProDial(text)
				if p.latestParams == text {
					now := time.Now()
					if err != nil {
						logrus.Errorf("pro port check error: %v", err)
						p.view.SetText(p.logFormat(gjson.Get(text, "network").String(), gjson.Get(text, "local_ip").String(), gjson.Get(text, "host").String(), gjson.Get(text, "port").String(), err.Error(), time.Now().Sub(now).String()) + p.view.Text)
					} else {
						logrus.Infof("pro port check result: %v", response)
						p.view.SetText(p.logFormat(gjson.Get(text, "network").String(), gjson.Get(text, "local_ip").String(), gjson.Get(text, "host").String(), gjson.Get(text, "port").String(), "OK !", time.Now().Sub(now).String()) + p.view.Text)
					}
				}
			}()
		},
	})
	border := container.NewBorder(nil, box, nil, nil, container.NewVScroll(multi))
	return container.NewBorder(top, bottom, left, right, border)
}

func (p *PortCheck) rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	p.view = widget.NewMultiLineEntry()
	p.view.Wrapping = fyne.TextWrapWord
	p.view.Scroll = container.ScrollVerticalOnly
	p.view.TextStyle = fyne.TextStyle{Bold: true}
	if p.logics.ViewText() != "" {
		p.view.SetText(p.logics.ViewText())
	} else {
		p.view.PlaceHolder = "检测结果"
	}
	p.view.OnChanged = func(s string) {
		logrus.Infof("port check result: %s", s)
		p.logics.SetViewText(s)
	}
	//view.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("port check view check clear: %s", p.view.Text)
			p.view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("port check view check copy: %s", p.view.Text)
			w.Clipboard().SetContent(strings.TrimSpace(p.view.Text))
		},
	})
	p.scroll = container.NewVScroll(p.view)
	border := container.NewBorder(nil, box, nil, nil, p.scroll)
	return container.NewBorder(top, bottom, left, right, border)
}

func (p *PortCheck) portCheckFrom() fyne.CanvasObject {
	networkSelect := widget.NewSelect(p.logics.NetworkList(), func(s string) {
		logrus.Infof("port check scheme: %s", s)
		p.logics.SetNetwork(s)
	})
	if p.logics.Network() != "" {
		networkSelect.SetSelected(p.logics.Network())
	} else {
		networkSelect.SetSelected("TCP")
	}

	interfaceSelect := widget.NewSelect(p.logics.GetInterfaceList(), func(s string) {
		logrus.Infof("port check interface: %s", s)
		p.logics.SetIFace(s)
	})
	if p.logics.IFace() != "" {
		interfaceSelect.SetSelected(p.logics.IFace())
	} else {
		interfaceSelect.SetSelected("自动")
	}

	host := widget.NewEntry()
	host.SetPlaceHolder("目标地址")
	host.SetText(p.logics.Host())
	host.OnChanged = func(s string) {
		logrus.Infof("port check host: %s", s)
		p.logics.SetHost(s)
		host.SetText(s)
	}
	host.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("port check validation host: %s", s)
		if s == "" {
			return nil
		}
		if ip := utils.FindIP(strings.TrimSpace(s)); ip != nil {
			return nil
		} else {
			return errors.New("地址格式错误！")
		}
	})

	port := widget.NewEntry()
	port.SetPlaceHolder("目标端口")
	port.SetText(p.logics.Port())
	port.OnChanged = func(s string) {
		logrus.Infof("port check port: %s", s)
		p.logics.SetPort(s)
		port.SetText(s)
	}
	port.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("prot check validation domain: %s", s)
		if s == "" {
			return nil
		}
		iPort, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			logrus.Warnf("strconv.Atoi error: %v", err)
			return errors.New("端口错误！")
		}
		if iPort > 65535 || iPort < 0 {
			return errors.New("端口不在合法范围内！")
		}
		return nil
	})

	var form *widget.StyleForm
	form = &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: "协议类型:", Widget: networkSelect, HintText: "必选"},
			{Text: "本地网卡:", Widget: interfaceSelect, HintText: "必选"},
			{Text: "目标地址:", Widget: host, HintText: "必填"},
			{Text: "目标端口:", Widget: port, HintText: "必填"},
		},
		OnCancel: func() {
			logrus.Infof("port check page cancelled")
			host.SetText("")
			port.SetText("")
		},
		OnSubmit: func() {
			p.latestParams = fmt.Sprintf("%s%s%s%s", networkSelect.Selected, interfaceSelect.Selected, host.Text, port.Text)
			Go.RelativelySafeGo(func() {
				logrus.Infof("port check page submitted")
				selected := networkSelect.Selected
				face := interfaceSelect.Selected
				text := host.Text
				targetPort := port.Text
				now := time.Now()
				response, err := p.logics.NormalDial(selected, face, text, targetPort)
				if p.latestParams == fmt.Sprintf("%s%s%s%s", selected, face, text, targetPort) {
					if err != nil {
						logrus.Errorf("port check error: %v", err)
						p.view.SetText(p.logFormat(selected, face, text, targetPort, err.Error(), time.Now().Sub(now).String()) + p.view.Text)
					} else {
						logrus.Infof("port check result: %v", response)
						p.view.SetText(p.logFormat(selected, face, text, targetPort, "OK !", time.Now().Sub(now).String()) + p.view.Text)
					}
				}
			})
		},
		SubmitText: "检测",
		CancelText: "清空",
		ButtonLayout: func(cancel *widget.Button, submit *widget.Button) *fyne.Container {
			return container.NewGridWithColumns(2, cancel, submit)
		},
		ContentLayout: func(input *fyne.Container, button *fyne.Container) *fyne.Container {
			return container.NewBorder(nil, button, nil, nil, input)
		},
	}
	return form
}

func (p *PortCheck) logFormat(network, i, host, port, msg, t string) string {
	return fmt.Sprintf("[%v] [%v] %v:%v(%v)[%v] => %v\n\n", time.Now().Format(`01-02 15:04:05`), t, host, port, network, i, msg)
}

func (p *PortCheck) ClearCache() {
	p.logics.ClearCache()
}
