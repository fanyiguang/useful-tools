package adapter

import (
	"fyne.io/fyne/v2"
)

type Menu interface {
	CreateMenu(a fyne.App, w fyne.Window, tutorials map[string]Page, setPage func(page Page), clearCacheFn func(), onLanguageChange func()) *fyne.MainMenu
}
