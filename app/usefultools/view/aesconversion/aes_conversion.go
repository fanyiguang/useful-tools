package aesconversion

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
)

var (
	logics       = controller.NewAesConversion()
	view         *widget.Entry
	latestParams string
	keyLoadOnce  sync.Once
)

func Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftScreen(mode), rightScreen(w))
}

func leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("aes conversion mode: %d", mode)
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)
	keyLoadOnce.Do(func() {
		logics.SetAesKey(fyne.CurrentApp().Preferences().String("aes-key"))
		logics.SetAesIv(fyne.CurrentApp().Preferences().String("aes-iv"))
	})
	return container.NewBorder(top, bottom, left, right, From())
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
		view.PlaceHolder = "解析结果"
	}
	view.OnChanged = func(s string) {
		logrus.Infof("aes conversion result: %s", s)
		logics.SetViewText(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("aes conversion view check clear: %s", view.Text)
			view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("aes conversion view check copy: %s", view.Text)
			w.Clipboard().SetContent(strings.TrimSpace(view.Text))
		},
	})
	return container.NewBorder(top, bottom, left, right, container.NewBorder(nil, box, nil, nil, view))
}

func From() fyne.CanvasObject {
	conversionSelect := widget.NewSelect(logics.ConversionList(), func(s string) {})
	if logics.ConversionType() != "" {
		conversionSelect.SetSelected(logics.ConversionType())
	} else {
		conversionSelect.SetSelected("解密")
	}

	aesKey := widget.NewEntry()
	aesKey.SetPlaceHolder("AES KEY")
	aesKey.SetText(logics.AesKey())
	aesKey.OnChanged = func(s string) {
		logrus.Infof("aes conversion aes key: %s", s)
		logics.SetAesKey(s)
		aesKey.SetText(s)
	}

	aesIV := widget.NewEntry()
	aesIV.SetPlaceHolder("AES IV")
	aesIV.SetText(logics.AesIv())
	aesIV.OnChanged = func(s string) {
		logrus.Infof("aes conversion aesIV: %s", s)
		logics.SetAesIv(s)
		aesIV.SetText(s)
	}

	data := widget.NewMultiLineEntry()
	data.Wrapping = fyne.TextWrapWord
	data.SetPlaceHolder("参数")
	data.SetText(logics.Data())
	data.OnChanged = func(s string) {
		logrus.Infof("aes conversion data: %s", s)
		logics.SetData(s)
		data.SetText(s)
	}
	data.SetMinRowsVisible(10)
	data.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("aes conversion data: %s", s)
		if s == "" {
			return errors.New("参数不可为空")
		}
		return nil
	})

	form := &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: "转换类型:", Widget: conversionSelect, HintText: "必选"},
			{Text: "AES KEY:", Widget: aesKey, HintText: "必填"},
			{Text: "AES IV:", Widget: aesIV, HintText: "必填"},
			{Text: "参数:", Widget: data, HintText: "必填"},
		},
		OnCancel: func() {
			logrus.Infof("aes conversion page cancelled")
			aesKey.SetText("")
			aesIV.SetText("")
		},
		OnSubmit: func() {
			latestParams = fmt.Sprintf("%s%s%s%s", conversionSelect.Selected, aesKey.Text, aesIV.Text, data.Text)
			go func() {
				logrus.Infof("aes conversion page submitted")
				key := aesKey.Text
				iv := aesIV.Text
				aesData := data.Text
				mode := conversionSelect.Selected
				response, err := logics.DoConversion(mode, key, iv, aesData)
				if latestParams == fmt.Sprintf("%s%s%s%s", conversionSelect.Selected, aesKey.Text, aesIV.Text, data.Text) { //防止重复提交
					if err != nil {
						logrus.Errorf("aes conversion error: %v", err)
						view.SetText(err.Error())
					} else {
						logrus.Infof("aes conversion result: %v", response)
						view.SetText(response)
					}
				}
			}()

			if fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceSaveAesKey) {
				fyne.CurrentApp().Preferences().SetString("aes-key", aesKey.Text)
				fyne.CurrentApp().Preferences().SetString("aes-iv", aesIV.Text)
			}
		},
		SubmitText: conversionSelect.Selected,
		CancelText: "清空",
		ButtonLayout: func(cancel *widget.Button, submit *widget.Button) *fyne.Container {
			return container.NewGridWithColumns(2, cancel, submit)
		},
		ContentLayout: func(input *fyne.Container, button *fyne.Container) *fyne.Container {
			return container.NewBorder(nil, button, nil, nil, input)
		},
	}
	conversionSelect.OnChanged = func(s string) {
		logrus.Infof("aes conversion mode: %s", s)
		logics.SetConversionType(s)
		form.SubmitText = conversionSelect.Selected
		form.Refresh()
	}
	return form
}
