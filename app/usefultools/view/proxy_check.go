package view

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"strconv"
	"useful-tools/app/usefultools/controller"
	"useful-tools/utils"
)

var logic = controller.NewProxyCheck()

var proxyType = []string{"SOCKS5", "SSL", "HTTP", "HTTPS", "SS"}

func proxyCheckScreen(w fyne.Window, mode ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftProxyCheck(mode), rightProxyCheck())
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

		box := container.NewGridWithColumns(2, &widget.Button{
			Text:       "清空",
			Icon:       theme.Icon(theme.IconNameContentClear),
			Importance: widget.MediumImportance,
			OnTapped:   func() { fmt.Println("high importance button") },
		}, &widget.Button{
			Text:       "检测",
			Icon:       theme.Icon(theme.IconNameContentCopy),
			Importance: widget.MediumImportance,
			OnTapped:   func() { fmt.Println("high importance button") },
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

func rightProxyCheck() fyne.CanvasObject {
	right := makeCellSize(10, 10)
	left := makeCellSize(10, 10)
	top := makeCellSize(10, 10)
	bottom := makeCellSize(10, 10)

	multi := widget.NewMultiLineEntry()
	multi.Wrapping = fyne.TextWrapWord
	multi.Text = "检测结果"
	multi.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped:   func() { fmt.Println("high importance button") },
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped:   func() { fmt.Println("high importance button") },
	})
	border := container.NewBorder(nil, box, nil, nil, container.NewVScroll(multi))
	return container.NewBorder(top, bottom, left, right, border)
}

func checkFrom() fyne.CanvasObject {
	proxyTypeSelect := widget.NewSelect(proxyType, func(s string) {
	})
	proxyTypeSelect.SetSelected("SOCKS5")

	host := widget.NewEntry()
	host.SetPlaceHolder("代理地址")
	host.Validator = validation.NewAllStrings(func(s string) error {
		if ip := utils.FindIP(s); ip != nil {
			return nil
		} else {
			return errors.New("地址格式错误！")
		}
	})

	port := widget.NewEntry()
	port.SetPlaceHolder("代理端口")
	port.Validator = validation.NewAllStrings(func(s string) error {
		iPort, err := strconv.Atoi(s)
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

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("代理密码")

	urls := widget.NewMultiLineEntry()
	urls.SetPlaceHolder("代理URL")

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
