package view

import (
	"fyne.io/fyne/v2"
	"net/url"
)

func createMenu(a fyne.App, w fyne.Window, setPage func(Page)) *fyne.MainMenu {
	var main *fyne.MainMenu
	checkedFn := func(item *fyne.MenuItem, extendFn func(*fyne.MenuItem)) {
		extendFn(item)
		item.Checked = !item.Checked
		main.Refresh()
	}
	majorItem := fyne.NewMenuItem("专业模式", nil)
	majorItem.Action = func() {
		checkedFn(majorItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetInt(NavStatePreferenceProMode, 0)
			} else {
				a.Preferences().SetInt(NavStatePreferenceProMode, 1)
			}
			if t, ok := Tutorials[fyne.CurrentApp().Preferences().String(NavStatePreferenceCurrentPage)]; ok {
				setPage(t)
			}
		})
	}
	majorItem.Checked = a.Preferences().Int(NavStatePreferenceProMode) == 1
	mode := fyne.NewMenu("模式", majorItem)

	hideBodyItem := fyne.NewMenuItem("隐藏响应体", nil)
	hideBodyItem.Action = func() {
		checkedFn(hideBodyItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetInt(NavStatePreferenceHideBody, 0)
			} else {
				a.Preferences().SetInt(NavStatePreferenceHideBody, 1)
			}
		})
	}
	hideBodyItem.Checked = a.Preferences().Int(NavStatePreferenceHideBody) == 1
	view := fyne.NewMenu("视图", hideBodyItem)

	saveAesItem := fyne.NewMenuItem("保存AES密钥", nil)
	saveAesItem.Action = func() {
		checkedFn(saveAesItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetInt(NavStatePreferenceSaveAesKey, 0)
			} else {
				a.Preferences().SetInt(NavStatePreferenceSaveAesKey, 1)
			}
		})
	}
	saveAesItem.Checked = a.Preferences().Int(NavStatePreferenceSaveAesKey) == 1
	closeUpgradeItem := fyne.NewMenuItem("关闭自动更新", nil)
	closeUpgradeItem.Action = func() {
		checkedFn(closeUpgradeItem, func(item *fyne.MenuItem) {
			if item.Checked {
				a.Preferences().SetInt(NavStatePreferenceCloseUpgrade, 0)
			} else {
				a.Preferences().SetInt(NavStatePreferenceCloseUpgrade, 1)
			}
		})
	}
	closeUpgradeItem.Checked = a.Preferences().Int(NavStatePreferenceCloseUpgrade) == 1
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
		mode, view, action, help,
	)
	return main
}
