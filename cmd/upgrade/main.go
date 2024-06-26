package main

import (
	"os"
	"useful-tools/app/upgrade"
	"useful-tools/pkg/wlog"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}
	err := upgrade.Upgrade(os.Args[1], os.Args[2])
	if err != nil {
		wlog.Warm("upgrade error: %v", err)
	}
}
