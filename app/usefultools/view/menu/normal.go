package menu

import (
	"net/url"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/view/constant"
	"useful-tools/app/usefultools/view/style"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/sirupsen/logrus"
)

var _ adapter.Menu = (*Normal)(nil)

type Normal struct {
}

func NewNormal() *Normal {
	return &Normal{}
}

func (m *Normal) CreateMenu(a fyne.App, w fyne.Window, tutorials map[string]adapter.Page, setPage func(adapter.Page), clearCacheFn func()) *fyne.MainMenu {
	var main *fyne.MainMenu
	checkedFn := func(item *fyne.MenuItem, extendFn func(*fyne.MenuItem)) {
		extendFn(item)
		item.Checked = !item.Checked
		main.Refresh()
	}

	clearCache := fyne.NewMenuItem("清除所有缓存", func() {
		logrus.Infof("clear cache")

		cnf := dialog.NewConfirm("清除所有缓存", "请问您要清理所有缓存吗?", m.confirmClearCacheCallback(clearCacheFn), w)
		cnf.SetDismissText("否")
		cnf.SetConfirmText("是")
		cnf.Show()
	})

	file := fyne.NewMenu("文件", clearCache)

	majorItem := fyne.NewMenuItem("专业模式", nil)
	majorItem.Action = func() {
		checkedFn(majorItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetInt(constant.NavStatePreferenceProMode, 0)
			} else {
				a.Preferences().SetInt(constant.NavStatePreferenceProMode, 1)
			}
			if t, ok := tutorials[fyne.CurrentApp().Preferences().String(constant.NavStatePreferenceCurrentPage)]; ok {
				setPage(t)
			}
		})
	}
	majorItem.Checked = a.Preferences().Int(constant.NavStatePreferenceProMode) == 1
	mode := fyne.NewMenu("模式", majorItem)

	hideBodyItem := fyne.NewMenuItem("隐藏响应体", nil)
	hideBodyItem.Action = func() {
		checkedFn(hideBodyItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetBool(constant.NavStatePreferenceHideBody, false)
			} else {
				a.Preferences().SetBool(constant.NavStatePreferenceHideBody, true)
			}
		})
	}
	hideBodyItem.Checked = a.Preferences().Bool(constant.NavStatePreferenceHideBody)
	styleKey := a.Preferences().String(constant.NavStatePreferenceStyle)
	if styleKey == "" {
		styleKey = constant.StyleDefault
		a.Preferences().SetString(constant.NavStatePreferenceStyle, styleKey)
	}
	style.Apply(a, styleKey)

	defaultStyleItem := fyne.NewMenuItem("默认", nil)
	lowGreenStyleItem := fyne.NewMenuItem("低饱和绿", nil)
	warmLuxuryStyleItem := fyne.NewMenuItem("轻奢暖调", nil)
	neutralMinimalStyleItem := fyne.NewMenuItem("中性极简", nil)
	setStyleChecked := func(key string) {
		defaultStyleItem.Checked = key == constant.StyleDefault
		lowGreenStyleItem.Checked = key == constant.StyleLowSaturationGreen
		warmLuxuryStyleItem.Checked = key == constant.StyleWarmLuxury
		neutralMinimalStyleItem.Checked = key == constant.StyleNeutralMinimal
	}
	applyStyle := func(key string) {
		a.Preferences().SetString(constant.NavStatePreferenceStyle, key)
		style.Apply(a, key)
		setStyleChecked(key)
		main.Refresh()
	}
	defaultStyleItem.Action = func() {
		applyStyle(constant.StyleDefault)
	}
	lowGreenStyleItem.Action = func() {
		applyStyle(constant.StyleLowSaturationGreen)
	}
	warmLuxuryStyleItem.Action = func() {
		applyStyle(constant.StyleWarmLuxury)
	}
	neutralMinimalStyleItem.Action = func() {
		applyStyle(constant.StyleNeutralMinimal)
	}
	setStyleChecked(styleKey)

	styleMenu := fyne.NewMenu("样式", defaultStyleItem, lowGreenStyleItem, warmLuxuryStyleItem, neutralMinimalStyleItem)
	styleItem := fyne.NewMenuItem("样式", nil)
	styleItem.ChildMenu = styleMenu
	view := fyne.NewMenu("视图", hideBodyItem, styleItem)

	saveAesItem := fyne.NewMenuItem("保存AES密钥", nil)
	saveAesItem.Action = func() {
		checkedFn(saveAesItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetBool(constant.NavStatePreferenceSaveAesKey, false)
			} else {
				a.Preferences().SetBool(constant.NavStatePreferenceSaveAesKey, true)
			}
		})
	}
	saveAesItem.Checked = a.Preferences().Bool(constant.NavStatePreferenceSaveAesKey)
	closeUpgradeItem := fyne.NewMenuItem("关闭自动更新", nil)
	closeUpgradeItem.Action = func() {
		checkedFn(closeUpgradeItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetBool(constant.NavStatePreferenceCloseUpgrade, false)
			} else {
				a.Preferences().SetBool(constant.NavStatePreferenceCloseUpgrade, true)
			}
		})
	}
	closeUpgradeItem.Checked = a.Preferences().Bool(constant.NavStatePreferenceCloseUpgrade)
	action := fyne.NewMenu("功能", saveAesItem, closeUpgradeItem)

	feedbackAesItem := fyne.NewMenuItem("反馈", func() {
		u, _ := url.Parse("https://github.com/fanyiguang/useful-tools/issues")
		_ = a.OpenURL(u)
	})
	helpItem := fyne.NewMenuItem("帮助", func() {
		u, _ := url.Parse("https://github.com/fanyiguang/useful-tools")
		_ = a.OpenURL(u)
	})
	help := fyne.NewMenu("帮助", feedbackAesItem, helpItem)

	main = fyne.NewMainMenu(
		file, mode, view, action, help,
	)
	return main
}

func (m *Normal) confirmClearCacheCallback(clearCacheFn func()) func(response bool) {
	return func(response bool) {
		logrus.Infof("clear cache: %v", response)
		if !response {
			return
		}

		clearCacheFn()
		logrus.Infof("clear cache success")
	}
}
