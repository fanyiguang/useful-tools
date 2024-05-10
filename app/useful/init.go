package useful

import (
	"fmt"
	"os"
	"path/filepath"
	"useful-tools/common/config"
	"useful-tools/helper/path"
	"useful-tools/pkg/wlog"
)

func initLogic() error {
	err := initConfig()
	if err != nil {
		return err
	}

	initFile()

	return nil
}

func initFile() {
	err := os.Rename(fmt.Sprintf("./%v", config.ProcessUpgradeName), filepath.Join(config.GetProjectsPath(), config.ProcessUpgradeName))
	if err != nil {
		wlog.Warm("init file error", err)
	}
}

func initConfig() error {
	s, err := path.Path()
	if err != nil {
		return err
	}
	create, err := os.Create(filepath.Join(config.GetConfigPath(), "path"))
	if err != nil {
		return err
	}
	_, err = create.WriteString(filepath.Dir(s))
	if err != nil {
		return err
	}
	return nil
}
