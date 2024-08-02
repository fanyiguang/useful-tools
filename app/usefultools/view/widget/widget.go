package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func MakeCell() fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.Transparent)
	rect.SetMinSize(fyne.NewSize(30, 30))
	return rect
}

func MakeCellSize(w, h float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(&color.Transparent)
	rect.SetMinSize(fyne.NewSize(w, h))
	return rect
}
