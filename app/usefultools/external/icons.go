package external

import "fyne.io/fyne/v2"

//go:generate go run ../../fyne bundle -package external -o bundled.go assets

var (
	Logo *fyne.StaticResource
)
