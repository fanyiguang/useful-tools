package usefultools

import (
	"fmt"
	"useful-tools/app/usefultools/model"
	"useful-tools/app/usefultools/view"
	"useful-tools/common/config"
)

func App() (err error) {
	err = initLogic()
	if err != nil {
		return err
	}

	backGround()
	view.Run(model.RunOptions{
		Id:       "useful-tools",
		Version:  config.Version,
		AppTitle: titleFormat(),
	})
	return
}

func titleFormat() string {
	return fmt.Sprintf("useful-tools v%v", config.Version)
}
