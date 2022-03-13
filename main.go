package main

import (
	"useful-tools/app"
	"useful-tools/pkg/wlog"
)

func main() {
	err := app.App()
	if err != nil {
		wlog.Error("%+v", err)
	}
}
