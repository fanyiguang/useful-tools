package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"useful-tools/app/upgrade"
	"useful-tools/common/config"
	"useful-tools/common/log"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("error")
		os.Exit(1)
	}

	log.Init(filepath.Join(config.GetLogPath(), "log.log"), "info")

	err := upgrade.Upgrade(os.Args[1], os.Args[2])
	if err != nil {
		logrus.Warnf("upgrade error: %v", err)
	}
}
