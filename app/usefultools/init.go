package usefultools

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"useful-tools/common/config"
	"useful-tools/common/log"
	"useful-tools/helper/file"
	"useful-tools/helper/path"
)

func initLogic() error {
	if config.IsTest() {
		log.Init(filepath.Join(config.GetLogPath(), "log.log"), "debug")
	} else {
		log.Init(filepath.Join(config.GetLogPath(), "log.log"), "info")
	}

	err := initConfig()
	if err != nil {
		return err
	}

	initFile()

	printVersion()

	backGround()

	return nil
}

func initFile() {
	_, err := file.CopyFile(filepath.Join(config.GetProjectsPath(), config.ProcessUpgradeName), fmt.Sprintf("./%v", config.ProcessUpgradeName))
	if err != nil {
		logrus.Warnf("copy file error: %v", err)
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
