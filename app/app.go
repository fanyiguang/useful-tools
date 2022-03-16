package app

import (
	walkUI "useful-tools/module/walk_ui"

	"github.com/lxn/walk"
)

func App() (err error) {
	_ = walk.Resources.SetRootDirPath(`D:\study\zixun\go_gui_walk\examples\img`)
	walkUI.New().Run()
	return
}
