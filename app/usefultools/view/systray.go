package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"log"
)

func createTrayMenu() *fyne.Menu {
	h := fyne.NewMenuItem("Hello", func() {})
	h.Icon = theme.HomeIcon()
	menu := fyne.NewMenu("systray", h)
	h.Action = func() {
		log.Println("System tray menu tapped")
		h.Label = "Welcome"
		menu.Refresh()
	}
	return menu
}
