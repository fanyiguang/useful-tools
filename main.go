package main

import (
	"network-tool/app"
	"network-tool/pkg/wlog"
)

func main() {
	err := app.App()
	if err != nil {
		wlog.Error("%+v", err)
	}
}
