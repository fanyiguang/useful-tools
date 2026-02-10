package page

import (
	"strings"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

var _ adapter.Page = (*FormatConversion)(nil)

type FormatConversion struct {
	BasePage
	logics       *controller.FormatConversion
	view         *widget.EntryEx
	inputEntry   *widget.EntryEx
	targetFormat *widget.Select
	formats      []string
}

func NewFormatConversion() *FormatConversion {
	fc := &FormatConversion{
		BasePage: BasePage{
			ID:         constant.PageIDFormatConversion,
			TitleKey:   i18n.KeyPageFormatConversionTitle,
			IntroKey:   i18n.KeyPageFormatConversionIntro,
			SupportWeb: true,
		},
		logics: controller.NewFormatConversion(),
	}
	fc.formats = fc.logics.GetFormats()
	return fc
}

func (f *FormatConversion) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(f.leftScreen(mode), f.rightScreen(w))
}

func (f *FormatConversion) leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("format conversion mode: %d", mode)
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)
	return container.NewBorder(top, bottom, left, right, f.inputSection())
}

func (f *FormatConversion) rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	f.view = widget.NewMultiLineEntryEx(nil, nil, nil, nil)
	f.view.Wrapping = fyne.TextWrapWord
	f.view.Scroll = container.ScrollVerticalOnly
	f.view.TextStyle = fyne.TextStyle{Bold: true}
	if f.logics.ViewText() != "" {
		f.view.SetText(f.logics.ViewText())
	} else {
		f.view.PlaceHolder = i18n.T(i18n.KeyJsonResultPlaceholder) // Reusing existing placeholder or create new one? Used generic result placeholder in keys
	}
	f.view.OnChanged = func(s string) {
		logrus.Infof("format conversion result: %s", s)
		f.logics.SetViewText(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonClear),
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("format conversion view clear")
			f.view.SetText("")
		},
	}, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonCopy),
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("format conversion view copy")
			w.Clipboard().SetContent(strings.TrimSpace(f.view.Text))
		},
	})
	return container.NewBorder(top, bottom, left, right, container.NewBorder(nil, box, nil, nil, f.view))
}

func (f *FormatConversion) inputSection() fyne.CanvasObject {
	// Target Format Selector
	f.formats = f.logics.GetFormats()
	f.targetFormat = widget.NewSelect(f.formats, func(s string) {
		logrus.Infof("format conversion target selected: %s", s)
	})
	
	if f.logics.TargetFormat() != "" {
		f.targetFormat.SetSelected(f.logics.TargetFormat())
	} else {
		// Default to JSON
		f.targetFormat.SetSelected(i18n.T(i18n.KeyFormatJson))
	}

	// Input Entry
	f.inputEntry = widget.NewMultiLineEntryEx(nil, nil, nil, nil)
	f.inputEntry.PlaceHolder = i18n.T(i18n.KeyFormatConversionInputLabel)
	f.inputEntry.Wrapping = fyne.TextWrapWord
	f.inputEntry.Scroll = container.ScrollVerticalOnly
	f.inputEntry.OnChanged = func(s string) {
		f.logics.SetData(s)
	}
	f.inputEntry.SetText(f.logics.Data())
	f.inputEntry.SetMinRowsVisible(15)

	// Process Button
	processBtn := &widget.Button{
		Text:       i18n.T(i18n.KeyButtonConvert),
		Icon:       theme.Icon(theme.IconNameConfirm),
		Importance: widget.HighImportance,
		OnTapped: func() {
			f.processConversion()
		},
	}

	// Clear Button
	clearInputBtn := &widget.Button{
		Text:       i18n.T(i18n.KeyButtonClear),
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			f.inputEntry.SetText("")
		},
	}

	form := &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: i18n.T(i18n.KeyFormatConversionTargetLabel), Widget: f.targetFormat, HintText: i18n.T(i18n.KeyHintSelectRequired)},
			{Text: i18n.T(i18n.KeyFormatConversionInputLabel), Widget: f.inputEntry, HintText: i18n.T(i18n.KeyHintRequired)},
		},
	}

	f.targetFormat.OnChanged = func(s string) {
		logrus.Infof("format conversion target selected: %s", s)
		f.logics.SetTargetFormat(s)
	}

	buttonRow := container.NewGridWithColumns(2, clearInputBtn, processBtn)

	return container.NewBorder(nil, buttonRow, nil, nil, form)
}

func (f *FormatConversion) processConversion() {
	target := f.targetFormat.Selected
	content := f.inputEntry.Text

	result, err := f.logics.Convert(target, content)
	if err != nil {
		f.view.SetText(i18n.T(i18n.KeyErrorPrefix) + err.Error())
		logrus.Errorf("format conversion error: %v", err)
		return
	}

	f.view.SetText(result)
	logrus.Infof("format conversion success")
}

func (f *FormatConversion) ClearCache() {
	f.view.SetText("")
	f.inputEntry.SetText("")
	f.targetFormat.SetSelectedIndex(0)
	f.logics.ClearCache()
}
