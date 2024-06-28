package widget

import (
	"fyne.io/fyne/v2"
)

func createMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	var main *fyne.MainMenu
	checkedFn := func(item *fyne.MenuItem) {
		item.Checked = !item.Checked
		main.Refresh()
	}
	majorItem := fyne.NewMenuItem("专业模式", nil)
	majorItem.Action = func() {
		checkedFn(majorItem)
	}
	mode := fyne.NewMenu("模式", majorItem)

	showPassItem := fyne.NewMenuItem("显示密码", nil)
	showPassItem.Action = func() {
		checkedFn(showPassItem)
	}
	hideBodyItem := fyne.NewMenuItem("隐藏响应体", nil)
	hideBodyItem.Action = func() {
		checkedFn(hideBodyItem)
	}
	view := fyne.NewMenu("视图", showPassItem, hideBodyItem)

	saveAesItem := fyne.NewMenuItem("保存AES密钥", nil)
	saveAesItem.Action = func() {
		checkedFn(saveAesItem)
	}
	closeUpgradeItem := fyne.NewMenuItem("关闭自动更新", nil)
	closeUpgradeItem.Action = func() {
		checkedFn(closeUpgradeItem)
	}
	action := fyne.NewMenu("功能", saveAesItem, closeUpgradeItem)

	feedbackAesItem := fyne.NewMenuItem("反馈", nil)
	helpItem := fyne.NewMenuItem("帮助", nil)
	help := fyne.NewMenu("帮助", feedbackAesItem, helpItem)

	main = fyne.NewMainMenu(
		mode, view, action, help,
	)
	return main
}
