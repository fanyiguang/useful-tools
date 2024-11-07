package adapter

import "fyne.io/fyne/v2"

type Systray interface {
	CreateTrayMenu(window fyne.Window) *fyne.Menu
}
