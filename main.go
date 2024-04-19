package main

import (
	"useful-tools/app/useful"
	"useful-tools/pkg/wlog"
)

func main() {
	err := useful.App()
	if err != nil {
		wlog.Error("%+v", err)
	}
}
