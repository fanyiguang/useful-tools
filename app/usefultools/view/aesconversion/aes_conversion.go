package aesconversion

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
)

var (
	logics       = controller.NewAesConversion()
	view         *widget.Entry
	scroll       *container.Scroll
	latestParams string
)

func Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return viewWidget.MakeCellSize(10, 10)
	//return container.NewHSplit(leftScreen(mode), rightScreen(w))
}

//
//func leftScreen(mode constant.ViewMode) fyne.CanvasObject {
//	logrus.Infof("aes conversion mode: %d", mode)
//	right := viewWidget.MakeCellSize(10, 10)
//	left := viewWidget.MakeCellSize(10, 10)
//	top := viewWidget.MakeCellSize(10, 10)
//	bottom := viewWidget.MakeCellSize(10, 10)
//	return container.NewBorder(top, bottom, left, right, From())
//}
//
//func rightScreen(w fyne.Window) fyne.CanvasObject {
//	right := viewWidget.MakeCellSize(10, 10)
//	left := viewWidget.MakeCellSize(10, 10)
//	top := viewWidget.MakeCellSize(10, 10)
//	bottom := viewWidget.MakeCellSize(10, 10)
//
//	view = widget.NewMultiLineEntry()
//	view.Wrapping = fyne.TextWrapWord
//	view.Scroll = container.ScrollVerticalOnly
//	view.TextStyle = fyne.TextStyle{Bold: true}
//	if logics.ViewText() != "" {
//		view.SetText(logics.ViewText())
//	} else {
//		view.Text = ""
//	}
//	view.OnChanged = func(s string) {
//		logrus.Infof("aes conversion result: %s", s)
//		logics.SetViewText(s)
//	}
//
//	box := container.NewGridWithColumns(2, &widget.Button{
//		Text:       "清空",
//		Icon:       theme.Icon(theme.IconNameContentClear),
//		Importance: widget.MediumImportance,
//		OnTapped: func() {
//			logrus.Infof("aes conversion view check clear: %s", view.Text)
//			view.SetText("")
//		},
//	}, &widget.Button{
//		Text:       "复制",
//		Icon:       theme.Icon(theme.IconNameContentCopy),
//		Importance: widget.MediumImportance,
//		OnTapped: func() {
//			logrus.Infof("aes conversion view check copy: %s", view.Text)
//			w.Clipboard().SetContent(strings.TrimSpace(view.Text))
//		},
//	})
//	border := container.NewBorder(nil, box, nil, nil, view)
//	return container.NewBorder(top, bottom, left, right, border)
//}
//
//func From() fyne.CanvasObject {
//	conversionSelect := widget.NewSelect(logics.ConversionList(), func(s string) {
//		logrus.Infof("port check scheme: %s", s)
//		logics.SetNetwork(s)
//	})
//	if logics.Network() != "" {
//		networkSelect.SetSelected(logics.Network())
//	} else {
//		networkSelect.SetSelected("TCP")
//	}
//
//	host := widget.NewEntry()
//	host.SetPlaceHolder("DNS地址")
//	if logics.Host() == "" {
//		host.SetText("默认")
//	} else {
//		host.SetText(logics.Host())
//	}
//	host.OnChanged = func(s string) {
//		logrus.Infof("aes conversion host: %s", s)
//		logics.SetHost(s)
//		host.SetText(s)
//	}
//
//	domain := widget.NewEntry()
//	domain.SetPlaceHolder("解析域名")
//	domain.SetText(logics.Domain())
//	domain.OnChanged = func(s string) {
//		logrus.Infof("aes conversion domain: %s", s)
//		logics.SetDomain(s)
//		domain.SetText(s)
//	}
//	domain.Validator = validation.NewAllStrings(func(s string) error {
//		logrus.Infof("aes conversion validation domain: %s", s)
//		if s == "" {
//			return nil
//		}
//		if !strings.Contains(s, ".") {
//			return errors.New("域名格式错误！")
//		}
//		return nil
//	})
//
//	var form *widget.StyleForm
//	form = &widget.StyleForm{
//		Items: []*widget.StyleFormItem{
//			{Text: "DNS地址:", Widget: host, HintText: "必填"},
//			{Text: "解析域名:", Widget: domain, HintText: "必填"},
//		},
//		OnCancel: func() {
//			logrus.Infof("aes conversion page cancelled")
//			host.SetText("")
//			domain.SetText("")
//		},
//		OnSubmit: func() {
//			latestParams = fmt.Sprintf("%s%s", host.Text, domain.Text)
//			go func() {
//				logrus.Infof("aes conversion page submitted")
//				host := host.Text
//				domain := domain.Text
//				response, err := logics.NormalDns(host, domain)
//				if latestParams == fmt.Sprintf("%s%s", host, domain) {
//					if err != nil {
//						logrus.Errorf("aes conversion error: %v", err)
//						view.SetText(logFormat(host, domain, err.Error()) + view.Text)
//					} else {
//						logrus.Infof("aes conversion result: %v", response)
//						view.SetText(logFormat(host, domain, strings.Join(response, " ")) + view.Text)
//					}
//				}
//			}()
//		},
//		SubmitText: "检测",
//		CancelText: "清空",
//		ButtonLayout: func(cancel *widget.Button, submit *widget.Button) *fyne.Container {
//			return container.NewGridWithColumns(2, cancel, submit)
//		},
//		ContentLayout: func(input *fyne.Container, button *fyne.Container) *fyne.Container {
//			return container.NewBorder(nil, button, nil, nil, input)
//		},
//	}
//	return form
//}
