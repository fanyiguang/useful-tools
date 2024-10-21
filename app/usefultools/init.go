package usefultools

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
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
	currentPath, err := path.Path()
	if err != nil {
		logrus.Warnf("path error: %v", err)
	}
	srcFilePath := filepath.Join(filepath.Dir(currentPath), config.ProcessUpgradeName)
	_, err = file.CopyFile(filepath.Join(config.GetProjectsPath(), config.ProcessUpgradeName), srcFilePath)
	if err != nil {
		logrus.Warnf("copy file error: %v", err)
		return
	}
	err = exec.Command("chmod", "+x", srcFilePath).Run()
	if err != nil {
		logrus.Warnf("chmod error: %v", err)
	}
	_ = os.Remove(srcFilePath)
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
