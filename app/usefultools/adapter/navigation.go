package adapter

import "fyne.io/fyne/v2"

type Navigation interface {
	CreateNavigation(setPage func(page Page), loadPrevious bool) fyne.CanvasObject
	Tutorials() map[string]Page
	ClearCache()
}
