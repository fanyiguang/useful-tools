package dnsquery

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
	"strings"
	"time"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
)

var (
	logics       = controller.NewDnsQuery()
	view         *widget.Entry
	scroll       *container.Scroll
	latestParams string
)

func Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftScreen(mode), rightScreen(w))
}

func leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("dns query mode: %d", mode)
	switch mode {
	case constant.ViewModePro:
		return proView()
	default:
		right := viewWidget.MakeCellSize(10, 10)
		left := viewWidget.MakeCellSize(10, 10)
		top := viewWidget.MakeCellSize(10, 10)
		bottom := viewWidget.MakeCellSize(10, 10)
		return container.NewBorder(top, bottom, left, right, From())
	}
}

func proView() fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	multi := widget.NewMultiLineEntryEx(nil, nil, nil, logics.FormatJson)
	multi.Wrapping = fyne.TextWrapWord
	if logics.PreModeInput() != "" {
		multi.SetText(logics.PreModeInput())
	} else {
		multi.SetText(logics.PreTemplate())
	}
	multi.OnChanged = func(s string) {
		logrus.Infof("dns query: %s", s)
		logics.SetPreModeInput(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("pre clear dns query: %s", multi.Text)
			multi.SetText("")
		},
	}, &widget.Button{
		Text:       "检测",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			latestParams = multi.Text
			go func() {
				text := multi.Text
				logrus.Infof("pre check dns query: %s", text)
				response, err := logics.ProQuery(text)
				if latestParams == multi.Text {
					if err != nil {
						logrus.Errorf("pro dns query error: %v", err)
						view.SetText(logFormat(gjson.Get(text, "server").String(), gjson.Get(text, "domain").String(), err.Error()) + view.Text)
					} else {
						logrus.Infof("pro dns query result: %v", response)
						view.SetText(logFormat(gjson.Get(text, "server").String(), gjson.Get(text, "domain").String(), strings.Join(response, " ")) + text)
					}
				}
			}()
		},
	})
	border := container.NewBorder(nil, box, nil, nil, multi)
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
		logrus.Infof("dns query result: %s", s)
		logics.SetViewText(s)
	}
	//view.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("dns query view check clear: %s", view.Text)
			view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("dns query view check copy: %s", view.Text)
			w.Clipboard().SetContent(strings.TrimSpace(view.Text))
		},
	})
	scroll = container.NewVScroll(view)
	border := container.NewBorder(nil, box, nil, nil, scroll)
	return container.NewBorder(top, bottom, left, right, border)
}

func From() fyne.CanvasObject {
	host := widget.NewEntry()
	host.SetPlaceHolder("DNS地址")
	if logics.Host() == "" {
		host.SetText("默认")
	} else {
		host.SetText(logics.Host())
	}
	host.OnChanged = func(s string) {
		logrus.Infof("dns query host: %s", s)
		logics.SetHost(s)
		host.SetText(s)
	}

	domain := widget.NewEntry()
	domain.SetPlaceHolder("解析域名")
	domain.SetText(logics.Domain())
	domain.OnChanged = func(s string) {
		logrus.Infof("dns query domain: %s", s)
		logics.SetDomain(s)
		domain.SetText(s)
	}
	domain.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("dns query validation domain: %s", s)
		if s == "" {
			return nil
		}
		if !strings.Contains(s, ".") {
			return errors.New("域名格式错误！")
		}
		return nil
	})

	var form *widget.StyleForm
	form = &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: "DNS地址:", Widget: host, HintText: "必填"},
			{Text: "解析域名:", Widget: domain, HintText: "必填"},
		},
		OnCancel: func() {
			logrus.Infof("dns query page cancelled")
			host.SetText("")
			domain.SetText("")
		},
		OnSubmit: func() {
			latestParams = fmt.Sprintf("%s%s", host.Text, domain.Text)
			go func() {
				logrus.Infof("dns query page submitted")
				host := host.Text
				domain := domain.Text
				response, err := logics.NormalDns(host, domain)
				if latestParams == fmt.Sprintf("%s%s", host, domain) {
					if err != nil {
						logrus.Errorf("dns query error: %v", err)
						view.SetText(logFormat(host, domain, err.Error()) + view.Text)
					} else {
						logrus.Infof("dns query result: %v", response)
						view.SetText(logFormat(host, domain, strings.Join(response, " ")) + view.Text)
					}
				}
			}()
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

func logFormat(server, host, msg string) string {
	return fmt.Sprintf("[%v] %v(%v) => %v\r\n\r\n", time.Now().Format(`06-01-02 15:04:05`), host, server, msg)
}
