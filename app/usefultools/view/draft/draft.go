package draft

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
)

var (
	logics = controller.NewDraft()
	view   *widget.Entry
	scroll *container.Scroll
)

func Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftScreen(mode), rightScreen(w))
}

func leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(0, 0)
	bottom := viewWidget.MakeCellSize(10, 10)

	var docs []*container.TabItem
	if len(logics.LeftDocs()) > 0 {
		for t, s := range logics.LeftDocs() {
			docs = append(docs, createLeftCanvasObject(t, s.PlaceHolder, s.Content))
		}
	} else {
		docs = append(docs, createLeftCanvasObject("草稿1", "草稿1", ""), createLeftCanvasObject("草稿2", "草稿2", ""), createLeftCanvasObject("草稿3", "草稿3", ""))
	}

	tabs := container.NewDocTabs(
		docs...,
	)
	i := 3
	tabs.CreateTab = func() *container.TabItem {
		i++
		return createLeftCanvasObject(fmt.Sprintf("草稿%d", i), fmt.Sprintf("草稿%d", i), "")
	}
	tabs.OnClosed = func(item *container.TabItem) {
		logrus.Infof("closed tab: %s", item.Text)
	}
	return container.NewBorder(top, bottom, left, right, tabs)
}

func createLeftCanvasObject(title, placeHolder, text string) *container.TabItem {
	entry := widget.NewMultiLineEntry()
	entry.PlaceHolder = placeHolder
	if text != "" {
		entry.SetText(text)
	}
	entry.OnChanged = func(s string) {
		logrus.Infof("left draft doc[%s] content: %s", title, s)
		logics.AddLeftDocs(title, s, placeHolder)
	}
	logics.AddLeftDocs(title, text, placeHolder)
	return container.NewTabItem(title, entry)
}

func createRightCanvasObject(title, placeHolder, text string) *container.TabItem {
	entry := widget.NewMultiLineEntry()
	entry.PlaceHolder = placeHolder
	if text != "" {
		entry.SetText(text)
	}
	entry.OnChanged = func(s string) {
		logrus.Infof("right draft doc[%s] content: %s", title, s)
		logics.AddLeftDocs(title, s, placeHolder)
	}
	logics.AddRightDocs(title, text, placeHolder)
	return container.NewTabItem(title, entry)
}

func rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(0, 0)
	bottom := viewWidget.MakeCellSize(10, 10)

	var docs []*container.TabItem
	if len(logics.RightDocs()) > 0 {
		for t, s := range logics.RightDocs() {
			docs = append(docs, createRightCanvasObject(t, s.PlaceHolder, s.Content))
		}
	} else {
		docs = append(docs, createRightCanvasObject("草稿1", "草稿1", ""), createRightCanvasObject("草稿2", "草稿2", ""), createRightCanvasObject("草稿3", "草稿3", ""))
	}

	tabs := container.NewDocTabs(
		docs...,
	)
	i := 3
	tabs.CreateTab = func() *container.TabItem {
		i++
		return createRightCanvasObject(fmt.Sprintf("草稿%d", i), fmt.Sprintf("草稿%d", i), "")
	}
	return container.NewBorder(top, bottom, left, right, tabs)
}
