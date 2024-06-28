package main

import (
	"useful-tools/app/usefultools"
	"useful-tools/pkg/wlog"
)

func main() {
	err := usefultools.App()
	if err != nil {
		wlog.Error("%+v", err)
	}
}
