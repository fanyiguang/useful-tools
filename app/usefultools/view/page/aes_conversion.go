package page

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
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/preferencekey"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
	"useful-tools/helper/Go"
)

var _ adapter.Page = (*AesConversion)(nil)

type AesConversion struct {
	BasePage
	logics       *controller.AesConversion
	view         *widget.EntryEx
	latestParams string
	aesKeyList   *preferencekey.AesListKey
}

func NewAesConversion() *AesConversion {
	return &AesConversion{
		BasePage: BasePage{
			Title:      "AES转换",
			Intro:      "加密，解密，解密，加密",
			SupportWeb: true,
		},
		logics:     controller.NewAesConversion(),
		aesKeyList: preferencekey.NewAesKeyList(),
	}

}

func (a *AesConversion) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(a.leftScreen(mode), a.rightScreen(w))
}

func (a *AesConversion) leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("aes conversion mode: %d", mode)
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)
	return container.NewBorder(top, bottom, left, right, a.from())
}

func (a *AesConversion) rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	a.view = widget.NewMultiLineEntryEx(nil, nil, nil, a.logics.FormatJson)
	a.view.Wrapping = fyne.TextWrapWord
	a.view.Scroll = container.ScrollVerticalOnly
	a.view.TextStyle = fyne.TextStyle{Bold: true}
	if a.logics.ViewText() != "" {
		a.view.SetText(a.logics.ViewText())
	} else {
		a.view.PlaceHolder = "解析结果"
	}
	a.view.OnChanged = func(s string) {
		logrus.Infof("aes conversion result: %s", s)
		a.logics.SetViewText(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("aes conversion view check clear: %s", a.view.Text)
			a.view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("aes conversion view check copy: %s", a.view.Text)
			w.Clipboard().SetContent(strings.TrimSpace(a.view.Text))
		},
	})
	return container.NewBorder(top, bottom, left, right, container.NewBorder(nil, box, nil, nil, a.view))
}

func (a *AesConversion) from() fyne.CanvasObject {
	var (
		keyNames []string
		keyList  map[string]string
		ivList   map[string]string
	)

	conversionSelect := widget.NewSelect(a.logics.ConversionList(), func(s string) {})
	if a.logics.ConversionType() != "" {
		conversionSelect.SetSelected(a.logics.ConversionType())
	} else {
		conversionSelect.SetSelected("解密")
	}

	aesKey := widget.NewPasswordEntry()
	aesKey.SetPlaceHolder("AES KEY")
	aesKey.SetText(a.logics.AesKey())
	aesKey.OnChanged = func(s string) {
		logrus.Infof("aes conversion aes key: %s", s)
		a.logics.SetAesKey(s)
		//aesKey.SetText(s)
	}

	aesIV := widget.NewPasswordEntry()
	aesIV.SetPlaceHolder("AES IV")
	aesIV.SetText(a.logics.AesIv())
	aesIV.OnChanged = func(s string) {
		logrus.Infof("aes conversion aesIV: %s", s)
		a.logics.SetAesIv(s)
		//aesIV.SetText(s)
	}

	keyNames, keyList, ivList = a.aesKeyList.GetValue()
	keyGroupSelect := widget.NewSelectEntry(keyNames)
	keyGroupSelect.OnChanged = func(s string) {
		logrus.Infof("aes conversion key group: %s", s)
		if keyList[s] != "" {
			aesKey.SetText(keyList[s])
			aesIV.SetText(ivList[s])
		}
	}
	if len(keyNames) > 0 {
		keyGroupSelect.SetText(keyNames[0])
	} else {
		keyGroupSelect.SetPlaceHolder("可以设置密钥名称后保存")
	}

	data := widget.NewMultiLineEntry()
	data.Wrapping = fyne.TextWrapWord
	data.SetPlaceHolder("参数")
	data.SetText(a.logics.Data())
	data.OnChanged = func(s string) {
		logrus.Infof("aes conversion data: %s", s)
		a.logics.SetData(s)
		//data.SetText(s)
	}
	data.SetMinRowsVisible(10)
	data.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("aes conversion data: %s", s)
		if s == "" {
			return errors.New("参数不可为空")
		}
		return nil
	})

	var form *widget.StyleForm
	form = &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: "转换类型:", Widget: conversionSelect, HintText: "必选"},
			{Text: "密钥名称:", Widget: keyGroupSelect, HintText: "选填"},
			{Text: "AES KEY:", Widget: aesKey, HintText: "必填"},
			{Text: "AES IV:", Widget: aesIV, HintText: "必填"},
			{Text: "参数:", Widget: data, HintText: "必填"},
		},
		OnCancel: func() {
			logrus.Infof("aes conversion page cancelled")
			data.SetText("")
		},
		OnSubmit: func() {
			a.latestParams = fmt.Sprintf("%s%s%s%s", conversionSelect.Selected, aesKey.Text, aesIV.Text, data.Text)
			Go.RelativelySafeGo(func() {
				logrus.Infof("aes conversion page submitted")
				key := aesKey.Text
				iv := aesIV.Text
				aesData := data.Text
				mode := conversionSelect.Selected
				response, err := a.logics.DoConversion(mode, key, iv, aesData)
				if a.latestParams == fmt.Sprintf("%s%s%s%s", conversionSelect.Selected, aesKey.Text, aesIV.Text, data.Text) { //防止重复提交
					if err != nil {
						logrus.Errorf("aes conversion error: %v", err)
						a.view.SetText(err.Error())
					} else {
						logrus.Infof("aes conversion result: %v", response)
						a.view.SetText(response)
					}
				}
			})

			if a.aesKeyList.SetValue(keyGroupSelect.Text, aesKey.Text, aesIV.Text) {
				keyNames, keyList, ivList = a.aesKeyList.GetValue()
				keyGroupSelect.SetOptions(keyNames)
				keyGroupSelect.Refresh()
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
		a.logics.SetConversionType(s)
		form.SubmitText = conversionSelect.Selected
		form.Refresh()
	}
	return form
}

func (a *AesConversion) ClearCache() {
	a.logics.ClearCache()
	a.aesKeyList.Clear()
}
