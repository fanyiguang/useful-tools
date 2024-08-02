package portcheck

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
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
	"useful-tools/utils"
)

var (
	logics = controller.NewPortCheck()
	view   *widget.Entry
	scroll *container.Scroll
)

func Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftScreen(mode), rightScreen(w))
}

func leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("proxy check mode: %d", mode)
	switch mode {
	case constant.ViewModePro:
		return proView()
	default:
		right := viewWidget.MakeCellSize(10, 10)
		left := viewWidget.MakeCellSize(10, 10)
		top := viewWidget.MakeCellSize(10, 10)
		bottom := viewWidget.MakeCellSize(10, 10)
		return container.NewBorder(top, bottom, left, right, portCheckFrom())
	}
}

func proView() fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	multi := widget.NewMultiLineEntry()
	multi.Wrapping = fyne.TextWrapWord
	if logics.PreModeInput() != "" {
		multi.SetText(logics.PreModeInput())
	} else {
		multi.SetText(logics.PreTemplate())
	}
	multi.OnChanged = func(s string) {
		logrus.Infof("port: %s", s)
		logics.SetPreModeInput(s)
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
			logrus.Infof("pre check port: %s", multi.Text)
			response, err := logics.ProDial(multi.Text)
			if err != nil {
				logrus.Errorf("pro port check error: %v", err)
				view.SetText(logFormat(gjson.Get(multi.Text, "network").String(), gjson.Get(multi.Text, "local_ip").String(), gjson.Get(multi.Text, "host").String(), gjson.Get(multi.Text, "port").String(), err.Error()) + view.Text)
			} else {
				logrus.Infof("pro port check result: %v", response)
				view.SetText(logFormat(gjson.Get(multi.Text, "network").String(), gjson.Get(multi.Text, "local_ip").String(), gjson.Get(multi.Text, "host").String(), gjson.Get(multi.Text, "port").String(), "OK !") + view.Text)
			}
		},
	})
	border := container.NewBorder(nil, box, nil, nil, container.NewVScroll(multi))
	return container.NewBorder(top, bottom, left, right, border)
}

func rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	view = widget.NewMultiLineEntry()
	view.Wrapping = fyne.TextWrapWord
	view.Scroll = container.ScrollVerticalOnly
	view.TextStyle = fyne.TextStyle{Bold: true}
	if logics.ViewText() != "" {
		view.SetText(logics.ViewText())
	} else {
		view.Text = ""
	}
	view.OnChanged = func(s string) {
		logrus.Infof("port check result: %s", s)
		logics.SetViewText(s)
	}
	//view.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("port check view check clear: %s", view.Text)
			view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("port check view check copy: %s", view.Text)
			w.Clipboard().SetContent(strings.TrimSpace(view.Text))
		},
	})
	scroll = container.NewVScroll(view)
	border := container.NewBorder(nil, box, nil, nil, scroll)
	return container.NewBorder(top, bottom, left, right, border)
}

func portCheckFrom() fyne.CanvasObject {
	networkSelect := widget.NewSelect(logics.NetworkList(), func(s string) {
		logrus.Infof("port check scheme: %s", s)
		logics.SetNetwork(s)
	})
	if logics.Network() != "" {
		networkSelect.SetSelected(logics.Network())
	} else {
		networkSelect.SetSelected("TCP")
	}

	interfaceSelect := widget.NewSelect(logics.GetInterfaceList(), func(s string) {
		logrus.Infof("port check interface: %s", s)
		logics.SetIFace(s)
	})
	if logics.IFace() != "" {
		interfaceSelect.SetSelected(logics.IFace())
	} else {
		interfaceSelect.SetSelected("自动")
	}

	host := widget.NewEntry()
	host.SetPlaceHolder("目标地址")
	host.SetText(logics.Host())
	host.OnChanged = func(s string) {
		logrus.Infof("port check host: %s", s)
		logics.SetHost(s)
		host.SetText(s)
	}
	host.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("port check validation host: %s", s)
		if ip := utils.FindIP(strings.TrimSpace(s)); ip != nil {
			return nil
		} else {
			return errors.New("地址格式错误！")
		}
	})

	port := widget.NewEntry()
	port.SetPlaceHolder("目标端口")
	port.SetText(logics.Port())
	port.OnChanged = func(s string) {
		logrus.Infof("port check port: %s", s)
		logics.SetPort(s)
		port.SetText(s)
	}
	port.Validator = validation.NewAllStrings(func(s string) error {
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
			logrus.Infof("port check page submitted")
			response, err := logics.NormalDial(networkSelect.Selected, interfaceSelect.Selected, host.Text, port.Text)
			if err != nil {
				logrus.Errorf("port check error: %v", err)
				view.SetText(logFormat(networkSelect.Selected, interfaceSelect.Selected, host.Text, port.Text, err.Error()) + view.Text)
			} else {
				logrus.Infof("port check result: %v", response)
				view.SetText(logFormat(networkSelect.Selected, interfaceSelect.Selected, host.Text, port.Text, "OK !") + view.Text)
			}
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

func logFormat(network, i, host, port, msg string) string {
	return fmt.Sprintf("[%v] %v:%v(%v)[%v] => %v\r\n\r\n", time.Now().Format(`01-02 15:04:05`), host, port, network, i, msg)
}
