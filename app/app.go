package app

import (
	walkUI "useful-tools/module/walk_ui"

	"github.com/lxn/walk"
)

func App() (err error) {
	_ = walk.Resources.SetRootDirPath(`./resource`)
	walkUI.New().Run()
	return
}
