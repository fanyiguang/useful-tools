package base

import "useful-tools/module/logic/app"

var (
	MenuItemLogic *app.MenuItem
)

func init() {
	MenuItemLogic = app.NewMenuItem()
}
