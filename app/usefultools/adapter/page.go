package adapter

import (
	"fyne.io/fyne/v2"
	"useful-tools/app/usefultools/view/constant"
)

type Page interface {
	Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject
	ClearCache()
	GetSupportWeb() bool
	GetIntro() string
	GetTitle() string
	GetID() string
}
