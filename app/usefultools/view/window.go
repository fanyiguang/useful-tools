package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/model"
	"useful-tools/app/usefultools/resource"
)

var topWindow fyne.Window

func Run(runOpt model.RunOptions) {
	application := app.NewWithID(runOpt.Id)
	application.SetIcon(resource.AppLogo)
	if desk, ok := application.(desktop.App); ok {
		desk.SetSystemTrayMenu(createTrayMenu())
	}
	setLifeCycle(application)
	content := container.NewStack()
	title := widget.NewLabel(runOpt.DefaultLabel)
	intro := widget.NewLabel(runOpt.DefaultLabel)
	intro.Wrapping = fyne.TextWrapWord

	window := application.NewWindow(runOpt.AppTitle)

	setPage := func(t Page) {
		if fyne.CurrentDevice().IsMobile() {
			child := application.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow, ViewMode(application.Preferences().Int(NavStatePreferenceProMode))))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = window
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)
		content.Objects = []fyne.CanvasObject{t.View(window, ViewMode(application.Preferences().Int(NavStatePreferenceProMode)))}
		content.Refresh()
	}

	window.SetMainMenu(createMenu(application, window, setPage))
	window.SetMaster()

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		window.SetContent(createNavigation(setPage, false))
	} else {
		split := container.NewHSplit(createNavigation(setPage, true), tutorial)
		split.Offset = 0.2
		window.SetContent(split)
	}
	window.MainMenu()
	window.Resize(fyne.NewSize(1100, 690))
	window.ShowAndRun()
}

func setLifeCycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(controller.OnStarted)
	a.Lifecycle().SetOnStopped(controller.OnStopped)
	a.Lifecycle().SetOnEnteredForeground(controller.OnEnteredForeground)
	a.Lifecycle().SetOnExitedForeground(controller.OnExitedForeground)
}
