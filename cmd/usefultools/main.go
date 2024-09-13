package main

import (
	"github.com/sirupsen/logrus"
	"useful-tools/app/usefultools"
)

func main() {
	err := usefultools.App()
	if err != nil {
		logrus.Errorf("%+v", err)
	}
}
