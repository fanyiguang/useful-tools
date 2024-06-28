package usefultools

import (
	"useful-tools/app/usefultools/widget"
	_ "useful-tools/common/config"
)

func App() (err error) {
	err = initLogic()
	if err != nil {
		return err
	}

	backGround()
	widget.Run()
	return
}
