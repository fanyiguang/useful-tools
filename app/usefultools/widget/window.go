package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"useful-tools/app/usefultools/external"
	"useful-tools/app/usefultools/logics"
)

const preferenceCurrentPage = "currentPage"

var topWindow fyne.Window

func Run() {
	application := app.NewWithID("useful-tools")
	application.SetIcon(external.Logo)
	if desk, ok := application.(desktop.App); ok {
		desk.SetSystemTrayMenu(createTrayMenu())
	}
	setLifeCycle(application)

	window := application.NewWindow("useful-tools")
	window.SetMainMenu(createMenu(application, window))
	window.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setPage := func(t Page) {
		if fyne.CurrentDevice().IsMobile() {
			child := application.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = window
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(window)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		window.SetContent(createNavigation(setPage, false))
	} else {
		split := container.NewHSplit(createNavigation(setPage, true), tutorial)
		split.Offset = 0.2
		window.SetContent(split)
	}
	window.Resize(fyne.NewSize(640, 460))
	window.ShowAndRun()
}

func setLifeCycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(logics.OnStarted)
	a.Lifecycle().SetOnStopped(logics.OnStopped)
	a.Lifecycle().SetOnEnteredForeground(logics.OnEnteredForeground)
	a.Lifecycle().SetOnExitedForeground(logics.OnExitedForeground)
}
