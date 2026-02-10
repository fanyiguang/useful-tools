package menu

import (
	"net/url"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/i18n"
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

func (m *Normal) CreateMenu(a fyne.App, w fyne.Window, tutorials map[string]adapter.Page, setPage func(adapter.Page), clearCacheFn func(), onLanguageChange func()) *fyne.MainMenu {
	var main *fyne.MainMenu
	checkedFn := func(item *fyne.MenuItem, extendFn func(*fyne.MenuItem)) {
		extendFn(item)
		item.Checked = !item.Checked
		main.Refresh()
	}

	clearCache := fyne.NewMenuItem(i18n.T(i18n.KeyMenuClearCache), func() {
		logrus.Infof("clear cache")

		cnf := dialog.NewConfirm(i18n.T(i18n.KeyDialogClearCacheTitle), i18n.T(i18n.KeyDialogClearCacheMessage), m.confirmClearCacheCallback(clearCacheFn), w)
		cnf.SetDismissText(i18n.T(i18n.KeyDialogNo))
		cnf.SetConfirmText(i18n.T(i18n.KeyDialogYes))
		cnf.Show()
	})

	file := fyne.NewMenu(i18n.T(i18n.KeyMenuFile), clearCache)

	majorItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuProMode), nil)
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
	mode := fyne.NewMenu(i18n.T(i18n.KeyMenuMode), majorItem)

	hideBodyItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuHideBody), nil)
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

	defaultStyleItem := fyne.NewMenuItem(i18n.T(i18n.KeyStyleDefault), nil)
	lowGreenStyleItem := fyne.NewMenuItem(i18n.T(i18n.KeyStyleLowSaturationGreen), nil)
	warmLuxuryStyleItem := fyne.NewMenuItem(i18n.T(i18n.KeyStyleWarmLuxury), nil)
	neutralMinimalStyleItem := fyne.NewMenuItem(i18n.T(i18n.KeyStyleNeutralMinimal), nil)
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

	styleMenu := fyne.NewMenu(i18n.T(i18n.KeyMenuStyle), defaultStyleItem, lowGreenStyleItem, warmLuxuryStyleItem, neutralMinimalStyleItem)
	styleItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuStyle), nil)
	styleItem.ChildMenu = styleMenu
	view := fyne.NewMenu(i18n.T(i18n.KeyMenuView), hideBodyItem, styleItem)

	saveAesItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuSaveAesKey), nil)
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
	closeUpgradeItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuCloseUpgrade), nil)
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
	langZhItem := fyne.NewMenuItem(i18n.T(i18n.KeyLangChinese), nil)
	langEnItem := fyne.NewMenuItem(i18n.T(i18n.KeyLangEnglish), nil)
	setLanguageChecked := func(lang string) {
		langZhItem.Checked = lang == i18n.LangZhCN
		langEnItem.Checked = lang == i18n.LangEnUS
	}
	applyLanguage := func(lang string) {
		a.Preferences().SetString(constant.NavStatePreferenceLanguage, lang)
		i18n.SetLanguage(lang)
		setLanguageChecked(lang)
		if onLanguageChange != nil {
			onLanguageChange()
		}
	}
	langZhItem.Action = func() {
		applyLanguage(i18n.LangZhCN)
	}
	langEnItem.Action = func() {
		applyLanguage(i18n.LangEnUS)
	}
	setLanguageChecked(a.Preferences().StringWithFallback(constant.NavStatePreferenceLanguage, i18n.LangZhCN))
	languageMenu := fyne.NewMenu(i18n.T(i18n.KeyMenuLanguage), langZhItem, langEnItem)
	languageItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuLanguage), nil)
	languageItem.ChildMenu = languageMenu

	action := fyne.NewMenu(i18n.T(i18n.KeyMenuAction), saveAesItem, closeUpgradeItem, languageItem)

	feedbackAesItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuFeedback), func() {
		u, _ := url.Parse("https://github.com/fanyiguang/useful-tools/issues")
		_ = a.OpenURL(u)
	})
	helpItem := fyne.NewMenuItem(i18n.T(i18n.KeyMenuHelpHome), func() {
		u, _ := url.Parse("https://github.com/fanyiguang/useful-tools")
		_ = a.OpenURL(u)
	})
	help := fyne.NewMenu(i18n.T(i18n.KeyMenuHelp), feedbackAesItem, helpItem)

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
