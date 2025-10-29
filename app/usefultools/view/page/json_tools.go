package page

import (
	"strings"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
)

var _ adapter.Page = (*JsonTools)(nil)

type JsonTools struct {
	BasePage
	logics     *controller.JsonTools
	view       *widget.EntryEx
	inputEntry *widget.EntryEx
	opType     *widget.Select
	operations []string
}

func NewJsonTools() *JsonTools {
	jsonTools := &JsonTools{
		BasePage: BasePage{
			Title:      "JSON工具",
			Intro:      "JSON格式化、压缩和去除转义工具",
			SupportWeb: true,
		},
		logics: controller.NewJsonTools(),
	}
	jsonTools.operations = jsonTools.logics.GetOperations()
	return jsonTools
}

func (j *JsonTools) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(j.leftScreen(mode), j.rightScreen(w))
}

func (j *JsonTools) leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("json tools mode: %d", mode)
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)
	return container.NewBorder(top, bottom, left, right, j.inputSection())
}

func (j *JsonTools) rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	j.view = widget.NewMultiLineEntryEx(nil, nil, nil, nil)
	j.view.Wrapping = fyne.TextWrapWord
	j.view.Scroll = container.ScrollVerticalOnly
	j.view.TextStyle = fyne.TextStyle{Bold: true}
	if j.logics.ViewText() != "" {
		j.view.SetText(j.logics.ViewText())
	} else {
		j.view.PlaceHolder = "处理结果"
	}
	j.view.OnChanged = func(s string) {
		logrus.Infof("json tools result: %s", s)
		j.logics.SetViewText(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("json tools view clear")
			j.view.SetText("")
		},
	}, &widget.Button{
		Text:       "复制",
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("json tools view copy")
			w.Clipboard().SetContent(strings.TrimSpace(j.view.Text))
		},
	})
	return container.NewBorder(top, bottom, left, right, container.NewBorder(nil, box, nil, nil, j.view))
}

func (j *JsonTools) inputSection() fyne.CanvasObject {
	// 创建操作类型选择器
	j.opType = widget.NewSelect(j.operations, func(s string) {
		logrus.Infof("json tools operation selected: %s", s)
	})
	if j.logics.ConversionType() != "" {
		j.opType.SetSelected(j.logics.ConversionType())
	} else {
		j.opType.SetSelected("压缩")
	}

	// 创建输入框
	j.inputEntry = widget.NewMultiLineEntryEx(nil, nil, nil, nil)
	j.inputEntry.PlaceHolder = "请输入JSON内容"
	j.inputEntry.Wrapping = fyne.TextWrapWord
	j.inputEntry.Scroll = container.ScrollVerticalOnly
	j.inputEntry.OnChanged = func(s string) {
		j.logics.SetData(s)
	}
	j.inputEntry.SetText(j.logics.Data())
	j.inputEntry.SetMinRowsVisible(15)

	// 创建处理按钮
	processBtn := &widget.Button{
		Text:       "处理",
		Icon:       theme.Icon(theme.IconNameConfirm),
		Importance: widget.HighImportance,
		OnTapped: func() {
			j.processJson()
		},
	}

	// 创建清空输入按钮
	clearInputBtn := &widget.Button{
		Text:       "清空",
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			j.inputEntry.SetText("")
		},
	}

	// 使用StyleForm布局，参考AES转换页面的实现方式
	form := &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: "操作类型:", Widget: j.opType, HintText: "必选"},
			{Text: "内容输入:", Widget: j.inputEntry, HintText: "必填"},
		},
	}

	j.opType.OnChanged = func(s string) {
		logrus.Infof("json tools operation selected: %s", s)
		j.logics.SetConversionType(s)
	}

	// 创建按钮区域
	buttonRow := container.NewGridWithColumns(2, clearInputBtn, processBtn)

	// 使用Border布局将按钮放在底部，与右侧页面的按钮对齐
	return container.NewBorder(nil, buttonRow, nil, nil, form)
}

func (j *JsonTools) processJson() {
	opType := j.opType.Selected
	content := j.inputEntry.Text

	result, err := j.logics.ProcessJson(opType, content)
	if err != nil {
		j.view.SetText("错误: " + err.Error())
		logrus.Errorf("json process error: %v", err)
		return
	}

	j.view.SetText(result)
	logrus.Infof("json process success")
}

func (j *JsonTools) ClearCache() {
	j.view.SetText("")
	j.inputEntry.SetText("")
	j.opType.SetSelectedIndex(0)
	j.logics.SetData("")
	j.logics.SetViewText("")
}
