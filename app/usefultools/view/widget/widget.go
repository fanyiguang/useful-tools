package widget

import (
	"image/color"
	"strings"
	"useful-tools/app/usefultools/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

func CopyToClipboard(w fyne.Window, content string) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return
	}
	w.Clipboard().SetContent(trimmed)
	if app := fyne.CurrentApp(); app != nil {
		app.SendNotification(&fyne.Notification{
			Title:   "useful-tools",
			Content: i18n.T(i18n.KeyNotificationCopySuccess),
		})
	}
}
