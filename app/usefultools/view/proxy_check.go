package view

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"useful-tools/app/usefultools/controller"
	"useful-tools/utils"
)

var (
	logic = controller.NewProxyCheck()
	view  *widget.Entry
)

func proxyCheckScreen(w fyne.Window, mode ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftProxyCheck(mode), rightProxyCheck(w))
}

func leftProxyCheck(mode ViewMode) fyne.CanvasObject {
	mode = ViewMode(fyne.CurrentApp().Preferences().Int(NavStatePreferenceProMode))
	logrus.Infof("proxy check mode: %d", mode)
	switch mode {
	case ViewModePro:
		right := makeCellSize(10, 10)
		left := makeCellSize(10, 10)
		top := makeCellSize(10, 10)
		bottom := makeCellSize(10, 10)

		multi := widget.NewMultiLineEntry()
		multi.Wrapping = fyne.TextWrapWord
		if logic.PreModeInput() != "" {
			multi.SetText(logic.PreModeInput())
		} else {
			multi.SetText(logic.PreTemplate())
		}
		multi.OnChanged = func(s string) {
			logrus.Infof("proxy: %s", s)
			logic.SetPreModeInput(s)
		}

		box := container.NewGridWithColumns(2, &widget.Button{
			Text:       "清空",
			Icon:       theme.Icon(theme.IconNameContentClear),
			Importance: widget.MediumImportance,
			OnTapped: func() {
				logrus.Infof("pre clear proxy: %s", multi.Text)
				multi.SetText("")
			},
		}, &widget.Button{
			Text:       "检测",
			Icon:       theme.Icon(theme.IconNameContentCopy),
			Importance: widget.MediumImportance,
			OnTapped: func() {
				logrus.Infof("pre check proxy: %s", multi.Text)
				response, err := logic.PreCheckProxy(multi.Text, fyne.CurrentApp().Preferences().Bool(NavStatePreferenceHideBody))
				if err != nil {
					logrus.Warnf("pre check proxy failed: %s", err)
					view.SetText(err.Error())
					return
				}
				logrus.Infof("pre check proxy result: %s", response)
				view.SetText(response)
			},
		})
		border := container.NewBorder(nil, box, nil, nil, container.NewVScroll(multi))
		return container.NewBorder(top, bottom, left, right, border)
	default:
		right := makeCellSize(10, 10)
		left := makeCellSize(10, 10)
		top := makeCellSize(10, 10)
		bottom := makeCellSize(10, 10)
		return container.NewBorder(top, bottom, left, right, checkFrom())
	}
}

func rightProxyCheck(w fyne.Window) fyne.CanvasObject {
	right := makeCellSize(10, 10)
	left := makeCellSize(10, 10)
	top := makeCellSize(10, 10)
	bottom := makeCellSize(10, 10)

	view = widget.NewMultiLineEntry()
	view.Wrapping = fyne.TextWrapWord
	if logic.ViewText() != "" {
		view.SetText(logic.ViewText())
	} else {
		view.Text = "检测结果"
	}
	view.TextStyle = fyne.TextStyle{Bold: true}
	view.OnChanged = func(s string) {
		logrus.Infof("proxy check result: %s", s)
		logic.SetViewText(s)
	}
	//view.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("proxy check view check clear: %s", view.Text)
			view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("proxy check view check copy: %s", view.Text)
			w.Clipboard().SetContent(view.Text)
		},
	})
	border := container.NewBorder(nil, box, nil, nil, container.NewVScroll(view))
	return container.NewBorder(top, bottom, left, right, border)
}

func checkFrom() fyne.CanvasObject {
	proxyTypeSelect := widget.NewSelect(logic.SupportProxyTypeList(), func(s string) {
		logrus.Infof("proxy check type: %s", s)
		logic.SetTyp(s)
	})
	if logic.Typ() != "" {
		proxyTypeSelect.SetSelected(logic.Typ())
	} else {
		proxyTypeSelect.SetSelected("SOCKS5")
	}

	host := widget.NewEntry()
	host.SetPlaceHolder("代理地址")
	host.SetText(logic.Host())
	host.OnChanged = func(s string) {
		logrus.Infof("proxy check host: %s", s)
		logic.SetHost(s)
		host.SetText(s)
	}
	host.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("proxy check validation host: %s", s)
		if ip := utils.FindIP(strings.TrimSpace(s)); ip != nil {
			return nil
		} else {
			return errors.New("地址格式错误！")
		}
	})

	port := widget.NewEntry()
	port.SetPlaceHolder("代理端口")
	port.SetText(logic.Port())
	port.OnChanged = func(s string) {
		logrus.Infof("proxy check port: %s", s)
		logic.SetPort(s)
		port.SetText(s)
	}
	port.Validator = validation.NewAllStrings(func(s string) error {
		iPort, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			logrus.Warnf("strconv.Atoi error: %v", err)
			return errors.New("代理端口错误！")
		}
		if iPort > 65535 || iPort < 0 {
			return errors.New("代理端口不在合法范围内！")
		}
		return nil
	})

	username := widget.NewEntry()
	username.SetPlaceHolder("代理账号")
	username.SetText(logic.Username())
	username.OnChanged = func(s string) {
		logrus.Infof("proxy check username: %s", s)
		logic.SetUsername(s)
	}

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("代理密码")
	password.SetText(logic.Password())
	password.OnChanged = func(s string) {
		logrus.Infof("proxy check password: %s", s)
		logic.SetPassword(s)
	}

	urls := widget.NewMultiLineEntry()
	urls.SetPlaceHolder("代理URL")
	urls.SetText(logic.Urls())
	urls.OnChanged = func(s string) {
		logrus.Infof("proxy check urls: %s", s)
		logic.SetUrls(s)
	}

	var form *widget.StyleForm
	form = &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: "代理类型:", Widget: proxyTypeSelect, HintText: "必选"},
			{Text: "代理地址:", Widget: host, HintText: "必填"},
			{Text: "代理端口:", Widget: port, HintText: "必填"},
			{Text: "代理账号:", Widget: username, HintText: "必填"},
			{Text: "代理密码:", Widget: password, HintText: "必填"},
			{Text: "检测地址:", Widget: urls, HintText: "选填"},
		},
		OnCancel: func() {
			logrus.Infof("proxy check page cancelled")
			proxyTypeSelect.SetSelected("SOCKS5")
			host.SetText("")
			port.SetText("")
			username.SetText("")
			password.SetText("")
			urls.SetText("")
		},
		OnSubmit: func() {
			logrus.Infof("proxy check page submitted")
			response, err := logic.NormalCheckProxy(host.Text, port.Text, username.Text, password.Text, proxyTypeSelect.Selected, urls.Text, fyne.CurrentApp().Preferences().Bool(NavStatePreferenceHideBody))
			if err != nil {
				logrus.Errorf("proxy check error: %v", err)
				view.SetText(err.Error())
				return
			}
			logrus.Infof("proxy check result: %s", response)
			view.SetText(response)
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
