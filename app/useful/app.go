package useful

import (
	walkUI "useful-tools/module/walk_ui"

	"github.com/lxn/walk"
	_ "useful-tools/common/config"
)

func App() (err error) {
	err = initConfig()
	if err != nil {
		return err
	}

	backGround()

	_ = walk.Resources.SetRootDirPath(`./resource`)
	walkUI.New().Run()
	return
}
