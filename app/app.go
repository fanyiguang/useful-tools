package app

import walkUI "network-tool/module/walk_ui"

func App() (err error) {
	walkUI.New().Run()
	return
}
