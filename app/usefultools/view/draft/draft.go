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
	logics    = controller.NewDraft()
	view      *widget.Entry
	scroll    *container.Scroll
	rightTabs *container.DocTabs
	leftTabs  *container.DocTabs
)

func Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(leftScreen(mode), rightScreen(mode))
}

func leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(0, 0)
	bottom := viewWidget.MakeCellSize(10, 10)

	var initTabItems []*container.TabItem

	var docs []*container.TabItem
	if len(logics.LeftDocs()) > 0 {
		for _, s := range logics.LeftDocs() {
			docs = append(docs, createLeftCanvasObject(s.Title, s.PlaceHolder, s.Content))
		}
	} else {
		initTabItems = []*container.TabItem{
			createLeftCanvasObject("草稿1", "草稿1", ""),
			createLeftCanvasObject("草稿2", "草稿2", ""),
			createLeftCanvasObject("草稿3", "草稿3", ""),
		}
		docs = append(docs, initTabItems...)
	}

	leftTabs = container.NewDocTabs(
		docs...,
	)
	i := len(initTabItems)
	leftTabs.CreateTab = func() *container.TabItem {
		i++
		return createLeftCanvasObject(fmt.Sprintf("草稿%d", i), fmt.Sprintf("草稿%d", i), "")
	}
	leftTabs.OnSelected = func(item *container.TabItem) {
		logrus.Infof("selected tab: %s", item.Text)
	}
	return container.NewBorder(top, bottom, left, right, leftTabs)
}

func createLeftCanvasObject(title, placeHolder, text string) *container.TabItem {
	leftMultiEntry := widget.NewMultiLineEntryEx(func() {
		if _, idx := logics.FindLeftNextDocsIndex(leftTabs.Selected().Text); idx >= 0 {
			leftTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find prev doc: %s", leftTabs.Selected().Text)
		}
	}, func() {
		if _, idx := logics.FindLeftPrevDocsIndex(leftTabs.Selected().Text); idx >= 0 {
			leftTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find next doc: %s", leftTabs.Selected().Text)
		}
	}, func() {
		leftTabs.Append(leftTabs.CreateTab())
		logrus.Infof("create new left doc: %s", leftTabs.Selected().Text)
		leftTabs.SelectIndex(len(leftTabs.Items) - 1)
	})
	leftMultiEntry.PlaceHolder = placeHolder
	if text != "" {
		leftMultiEntry.SetText(text)
	}
	leftMultiEntry.OnChanged = func(s string) {
		logrus.Infof("left draft doc[%s] content: %s", title, s)
		logics.AddLeftDocs(title, s, placeHolder)
	}
	logics.AddLeftDocs(title, text, placeHolder)
	return container.NewTabItem(title, leftMultiEntry)
}

func createRightCanvasObject(title, placeHolder, text string) *container.TabItem {
	rightMultiEntry := widget.NewMultiLineEntryEx(func() {
		if _, idx := logics.FindRightNextDocsIndex(rightTabs.Selected().Text); idx >= 0 {
			rightTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find prev doc: %s", rightTabs.Selected().Text)
		}
	}, func() {
		if _, idx := logics.FindRightPrevDocsIndex(rightTabs.Selected().Text); idx >= 0 {
			rightTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find next doc: %s", rightTabs.Selected().Text)
		}
	}, func() {
		rightTabs.Append(rightTabs.CreateTab())
		logrus.Infof("create new right doc: %s", rightTabs.Selected().Text)
		rightTabs.SelectIndex(len(rightTabs.Items) - 1)
	})
	rightMultiEntry.PlaceHolder = placeHolder
	if text != "" {
		rightMultiEntry.SetText(text)
	}
	rightMultiEntry.OnChanged = func(s string) {
		logrus.Infof("right draft doc[%s] content: %s", title, s)
		logics.AddRightDocs(title, s, placeHolder)
	}
	logics.AddRightDocs(title, text, placeHolder)
	return container.NewTabItem(title, rightMultiEntry)
}

func rightScreen(mode constant.ViewMode) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(0, 0)
	bottom := viewWidget.MakeCellSize(10, 10)

	var initTabItems []*container.TabItem

	var docs []*container.TabItem
	if len(logics.RightDocs()) > 0 {
		for _, s := range logics.RightDocs() {
			docs = append(docs, createRightCanvasObject(s.Title, s.PlaceHolder, s.Content))
		}
	} else {
		initTabItems = []*container.TabItem{
			createRightCanvasObject("草稿1", "草稿1", ""),
			createRightCanvasObject("草稿2", "草稿2", ""),
			createRightCanvasObject("草稿3", "草稿3", ""),
		}
		docs = append(docs, initTabItems...)
	}

	rightTabs = container.NewDocTabs(
		docs...,
	)
	i := len(initTabItems)
	rightTabs.CreateTab = func() *container.TabItem {
		i++
		return createRightCanvasObject(fmt.Sprintf("草稿%d", i), fmt.Sprintf("草稿%d", i), "")
	}
	rightTabs.OnSelected = func(item *container.TabItem) {
		logrus.Infof("selected tab: %s", item.Text)
	}
	return container.NewBorder(top, bottom, left, right, rightTabs)
}
