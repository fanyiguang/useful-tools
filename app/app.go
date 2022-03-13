package app

import walkUI "useful-tools/module/walk_ui"

func App() (err error) {
	walkUI.New().Run()
	return
}
