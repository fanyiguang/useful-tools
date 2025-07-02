package main

import (
	"useful-tools/app/usefultools"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	err := usefultools.App()
	if err != nil {
		logrus.Errorf("%+v", err)
	}
}
