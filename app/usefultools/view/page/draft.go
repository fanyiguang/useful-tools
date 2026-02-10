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
	"github.com/pmezard/go-difflib/difflib"
	"fyne.io/fyne/v2/theme"
	"strings"
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
	var leftIdx, rightIdx int
	if d.leftTabs != nil {
		leftIdx = d.leftTabs.SelectedIndex()
	}
	if d.rightTabs != nil {
		rightIdx = d.rightTabs.SelectedIndex()
	}

	split := container.NewHSplit(d.leftScreen(mode), d.rightScreen(mode))

	if d.leftTabs != nil && leftIdx < len(d.leftTabs.Items) {
		d.leftTabs.SelectIndex(leftIdx)
	}
	if d.rightTabs != nil && rightIdx < len(d.rightTabs.Items) {
		d.rightTabs.SelectIndex(rightIdx)
	}

	if fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceTextCompare) {
		d.updateDiffViews()
	}

	return split
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
		d.updateDiffViews()
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
		d.updateDiffViews()
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
		d.updateDiffViews()
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
		d.updateDiffViews()
	}
	return container.NewBorder(top, bottom, left, right, d.rightTabs)
}

func (d *Draft) ClearCache() {
	d.logics.ClearCache()
}

func (d *Draft) updateDiffViews() {
	if !fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceTextCompare) {
		return
	}
	if d.leftTabs == nil || d.rightTabs == nil {
		return
	}

	d.updateTabDiff(d.leftTabs.Selected(), true)
	d.updateTabDiff(d.rightTabs.Selected(), false)

	d.leftTabs.Refresh()
	d.rightTabs.Refresh()
}

func (d *Draft) updateTabDiff(item *container.TabItem, isLeft bool) {
	if item == nil {
		return
	}

	leftItem := d.leftTabs.Selected()
	rightItem := d.rightTabs.Selected()
	if leftItem == nil || rightItem == nil {
		return
	}

	leftDoc, _ := d.logics.GetLeftDoc(leftItem.Text)
	rightDoc, _ := d.logics.GetRightDoc(rightItem.Text)

	diffView := d.createDiffView(leftDoc.Content, rightDoc.Content, isLeft)

	if split, ok := item.Content.(*container.Split); ok {
		split.Trailing = diffView
		split.Refresh()
	} else {
		// First time entering compare mode or switching tabs
		entry := item.Content
		split := container.NewHSplit(entry, diffView)
		item.Content = split
		item.Content.Refresh()
	}
}

func (d *Draft) createDiffView(text1, text2 string, isLeft bool) fyne.CanvasObject {
	diff := difflib.NewMatcher(difflib.SplitLines(text1), difflib.SplitLines(text2))
	opcodes := diff.GetOpCodes()

	var segments []widget.RichTextSegment

	for _, op := range opcodes {
		switch op.Tag {
		case 'e':
			// Equal
			str := strings.Join(difflib.SplitLines(text1)[op.I1:op.I2], "")
			segments = append(segments, &widget.TextSegment{
				Text: str,
				Style: widget.RichTextStyle{
					Inline: true,
				},
			})
		case 'd':
			// Delete
			if isLeft {
				str := strings.Join(difflib.SplitLines(text1)[op.I1:op.I2], "")
				segments = append(segments, &widget.TextSegment{
					Text: str,
					Style: widget.RichTextStyle{
						Inline:    true,
						ColorName: theme.ColorNameError,
					},
				})
			}
		case 'i':
			// Insert
			if !isLeft {
				str := strings.Join(difflib.SplitLines(text2)[op.J1:op.J2], "")
				segments = append(segments, &widget.TextSegment{
					Text: str,
					Style: widget.RichTextStyle{
						Inline:    true,
						ColorName: theme.ColorNameSuccess,
					},
				})
			}
		case 'r':
			// Replace
			if isLeft {
				str := strings.Join(difflib.SplitLines(text1)[op.I1:op.I2], "")
				segments = append(segments, &widget.TextSegment{
					Text: str,
					Style: widget.RichTextStyle{
						Inline:    true,
						ColorName: theme.ColorNameError,
					},
				})
			} else {
				str := strings.Join(difflib.SplitLines(text2)[op.J1:op.J2], "")
				segments = append(segments, &widget.TextSegment{
					Text: str,
					Style: widget.RichTextStyle{
						Inline:    true,
						ColorName: theme.ColorNameSuccess,
					},
				})
			}
		}
	}
	r := widget.NewRichText(segments...)
	r.Wrapping = fyne.TextWrapBreak
	return container.NewVScroll(r)
}
