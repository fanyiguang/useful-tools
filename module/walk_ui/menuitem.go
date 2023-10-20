package walkUI

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"useful-tools/helper/Go"
	"useful-tools/helper/proc"
	"useful-tools/module/walk_ui/base"
)

var (
	ConvenientModeMenu *walk.Action
	ShowPassMenu       *walk.Action
	HiddenBodyMenu     *walk.Action
	SaveAesKeyMenu     *walk.Action
)

func MenuItems(mw *AppMainWindow) []MenuItem {
	return []MenuItem{
		Menu{
			Text: "模式",
			Items: []MenuItem{
				Action{
					AssignTo: &ConvenientModeMenu,
					Text:     "专业模式",
					//Checked: Bind("openHiddenCB.Visible"),
					Checked: false,
					OnTriggered: func() {
						switch ConvenientModeMenu.Checked() {
						case false:
							_ = ConvenientModeMenu.SetChecked(true)
							base.MenuItemLogic.SetProMode(true)
							go base.MenuItemLogic.SetProModeToFile(1)
						case true:
							_ = ConvenientModeMenu.SetChecked(false)
							base.MenuItemLogic.SetProMode(false)
							go base.MenuItemLogic.SetProModeToFile(0)
						}
						_ = mw.SetCurrentAction(mw.CurrentAction())
					},
				},
			},
		},
		Menu{
			Text: "视图",
			Items: []MenuItem{
				Action{
					AssignTo: &ShowPassMenu,
					Text:     "显示密码",
					//Checked: Bind("openHiddenCB.Visible"),
					Checked: false,
					OnTriggered: func() {
						switch ShowPassMenu.Checked() {
						case false:
							_ = ShowPassMenu.SetChecked(true)
							base.MenuItemLogic.SetShowPass(true)
							go base.MenuItemLogic.SetShowPassToFile(1)
						case true:
							_ = ShowPassMenu.SetChecked(false)
							base.MenuItemLogic.SetShowPass(false)
							go base.MenuItemLogic.SetShowPassToFile(0)
						}
						if mw.CurrentAction().Text() == "代理检测" {
							_ = mw.SetCurrentAction(mw.CurrentAction())
						}
					},
				},
				Action{
					AssignTo: &HiddenBodyMenu,
					Text:     "隐藏响应体",
					//Checked: Bind("openHiddenCB.Visible"),
					Checked: false,
					OnTriggered: func() {
						switch HiddenBodyMenu.Checked() {
						case false:
							_ = HiddenBodyMenu.SetChecked(true)
							base.MenuItemLogic.SetHiddenBody(true)
							go base.MenuItemLogic.SetHiddenBodyToFile(1)
						case true:
							_ = HiddenBodyMenu.SetChecked(false)
							base.MenuItemLogic.SetHiddenBody(false)
							go base.MenuItemLogic.SetHiddenBodyToFile(0)
						}
						//_ = mw.SetCurrentAction(mw.CurrentAction())
					},
				},
			},
		},
		Menu{
			Text: "功能",
			Items: []MenuItem{
				Action{
					AssignTo: &SaveAesKeyMenu,
					Text:     "保存AES密钥",
					//Checked: Bind("openHiddenCB.Visible"),
					Checked: false,
					OnTriggered: func() {
						switch SaveAesKeyMenu.Checked() {
						case false:
							_ = SaveAesKeyMenu.SetChecked(true)
							base.MenuItemLogic.SetSaveAesKey(true)
							go base.MenuItemLogic.SetSaveAesKeyToFile(1)
						case true:
							_ = SaveAesKeyMenu.SetChecked(false)
							base.MenuItemLogic.SetSaveAesKey(false)
							go base.MenuItemLogic.SetSaveAesKeyToFile(0)
						}
					},
				},
			},
		},
		Menu{
			Text: "帮助",
			Items: []MenuItem{
				Action{
					Text: "反馈",
					OnTriggered: func() {
						Go.Go(func() {
							_ = proc.RunProcByWin32Api(`https://github.com/fanyiguang/useful-tools/issues`, true)
						})
					},
				},
				Action{
					Text: "帮助",
					OnTriggered: func() {
						Go.Go(func() {
							_ = proc.RunProcByWin32Api(`https://github.com/fanyiguang/useful-tools`, true)
						})
					},
				},
			},
		},
	}
}

func initMenuItems() {
	if base.MenuItemLogic.ProMode() {
		_ = ConvenientModeMenu.SetChecked(true)
	}
	if base.MenuItemLogic.ShowPass() {
		_ = ShowPassMenu.SetChecked(true)
	}
	if base.MenuItemLogic.HiddenBody() {
		_ = HiddenBodyMenu.SetChecked(true)
	}
	if base.MenuItemLogic.SaveAesKey() {
		_ = SaveAesKeyMenu.SetChecked(true)
	}
}
