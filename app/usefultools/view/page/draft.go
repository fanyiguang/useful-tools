package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
)

var _ adapter.Page = (*Draft)(nil)

type Draft struct {
	BasePage
	logics    *controller.Draft
	rightTabs *container.DocTabs
	leftTabs  *container.DocTabs
}

func NewDraft() *Draft {
	return &Draft{
		BasePage: BasePage{
			ID:         constant.PageIDDraft,
			TitleKey:   i18n.KeyPageDraftTitle,
			IntroKey:   i18n.KeyPageDraftIntro,
			SupportWeb: true,
		},
		logics: controller.NewDraft(),
	}
}

func (d *Draft) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(d.leftScreen(mode), d.rightScreen(mode))
}

func (d *Draft) leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(0, 0)
	bottom := viewWidget.MakeCellSize(10, 10)

	var initTabItems []*container.TabItem

	var docs []*container.TabItem
	if len(d.logics.LeftDocs()) > 0 {
		for _, s := range d.logics.LeftDocs() {
			docs = append(docs, d.createLeftCanvasObject(s.Title, s.PlaceHolder, s.Content))
		}
	} else {
		initTabItems = []*container.TabItem{
			d.createLeftCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, 1), i18n.Tf(i18n.KeyDraftTabTitle, 1), ""),
			d.createLeftCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, 2), i18n.Tf(i18n.KeyDraftTabTitle, 2), ""),
			d.createLeftCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, 3), i18n.Tf(i18n.KeyDraftTabTitle, 3), ""),
		}
		docs = append(docs, initTabItems...)
	}

	d.leftTabs = container.NewDocTabs(
		docs...,
	)
	i := len(initTabItems)
	d.leftTabs.CreateTab = func() *container.TabItem {
		i++
		return d.createLeftCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, i), i18n.Tf(i18n.KeyDraftTabTitle, i), "")
	}
	d.leftTabs.OnSelected = func(item *container.TabItem) {
		logrus.Infof("selected tab: %s", item.Text)
	}
	return container.NewBorder(top, bottom, left, right, d.leftTabs)
}

func (d *Draft) createLeftCanvasObject(title, placeHolder, text string) *container.TabItem {
	leftMultiEntry := widget.NewMultiLineEntryEx(func() {
		if _, idx := d.logics.FindLeftNextDocsIndex(d.leftTabs.Selected().Text); idx >= 0 {
			d.leftTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find prev doc: %s", d.leftTabs.Selected().Text)
		}
	}, func() {
		if _, idx := d.logics.FindLeftPrevDocsIndex(d.leftTabs.Selected().Text); idx >= 0 {
			d.leftTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find next doc: %s", d.leftTabs.Selected().Text)
		}
	}, func() {
		d.leftTabs.Append(d.leftTabs.CreateTab())
		logrus.Infof("create new left doc: %s", d.leftTabs.Selected().Text)
		d.leftTabs.SelectIndex(len(d.leftTabs.Items) - 1)
	}, d.logics.FormatJson)

	leftMultiEntry.PlaceHolder = placeHolder
	if text != "" {
		leftMultiEntry.SetText(text)
	}
	leftMultiEntry.OnChanged = func(s string) {
		logrus.Infof("left draft doc[%s] content: %s", title, s)
		d.logics.AddLeftDocs(title, s, placeHolder)
	}
	leftMultiEntry.Wrapping = fyne.TextWrapWord
	d.logics.AddLeftDocs(title, text, placeHolder)
	return container.NewTabItem(title, leftMultiEntry)
}

func (d *Draft) createRightCanvasObject(title, placeHolder, text string) *container.TabItem {
	rightMultiEntry := widget.NewMultiLineEntryEx(func() {
		if _, idx := d.logics.FindRightNextDocsIndex(d.rightTabs.Selected().Text); idx >= 0 {
			d.rightTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find prev doc: %s", d.rightTabs.Selected().Text)
		}
	}, func() {
		if _, idx := d.logics.FindRightPrevDocsIndex(d.rightTabs.Selected().Text); idx >= 0 {
			d.rightTabs.SelectIndex(idx)
		} else {
			logrus.Warnf("can't find next doc: %s", d.rightTabs.Selected().Text)
		}
	}, func() {
		d.rightTabs.Append(d.rightTabs.CreateTab())
		logrus.Infof("create new right doc: %s", d.rightTabs.Selected().Text)
		d.rightTabs.SelectIndex(len(d.rightTabs.Items) - 1)
	}, d.logics.FormatJson)
	rightMultiEntry.PlaceHolder = placeHolder
	if text != "" {
		rightMultiEntry.SetText(text)
	}
	rightMultiEntry.OnChanged = func(s string) {
		logrus.Infof("right draft doc[%s] content: %s", title, s)
		d.logics.AddRightDocs(title, s, placeHolder)
	}
	rightMultiEntry.Wrapping = fyne.TextWrapWord
	d.logics.AddRightDocs(title, text, placeHolder)
	return container.NewTabItem(title, rightMultiEntry)
}

func (d *Draft) rightScreen(mode constant.ViewMode) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(0, 0)
	bottom := viewWidget.MakeCellSize(10, 10)

	var initTabItems []*container.TabItem

	var docs []*container.TabItem
	if len(d.logics.RightDocs()) > 0 {
		for _, s := range d.logics.RightDocs() {
			docs = append(docs, d.createRightCanvasObject(s.Title, s.PlaceHolder, s.Content))
		}
	} else {
		initTabItems = []*container.TabItem{
			d.createRightCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, 1), i18n.Tf(i18n.KeyDraftTabTitle, 1), ""),
			d.createRightCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, 2), i18n.Tf(i18n.KeyDraftTabTitle, 2), ""),
			d.createRightCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, 3), i18n.Tf(i18n.KeyDraftTabTitle, 3), ""),
		}
		docs = append(docs, initTabItems...)
	}

	d.rightTabs = container.NewDocTabs(
		docs...,
	)
	i := len(initTabItems)
	d.rightTabs.CreateTab = func() *container.TabItem {
		i++
		return d.createRightCanvasObject(i18n.Tf(i18n.KeyDraftTabTitle, i), i18n.Tf(i18n.KeyDraftTabTitle, i), "")
	}
	d.rightTabs.OnSelected = func(item *container.TabItem) {
		logrus.Infof("selected tab: %s", item.Text)
	}
	return container.NewBorder(top, bottom, left, right, d.rightTabs)
}

func (d *Draft) ClearCache() {
	d.logics.ClearCache()
}
