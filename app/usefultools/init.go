package usefultools

import (
	"fmt"
	"os"
	"path/filepath"
	"useful-tools/common/config"
	"useful-tools/helper/file"
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
	_, err := file.CopyFile(filepath.Join(config.GetProjectsPath(), config.ProcessUpgradeName), fmt.Sprintf("./%v", config.ProcessUpgradeName))
	if err != nil {
		wlog.Warm("init file error", err)
		return
	}
	_ = os.Remove(fmt.Sprintf("./%v", config.ProcessUpgradeName))
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
