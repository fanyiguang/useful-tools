package view

import (
	"fmt"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/model"
	"useful-tools/app/usefultools/resource"
	"useful-tools/app/usefultools/view/constant"
	"useful-tools/app/usefultools/view/menu"
	"useful-tools/app/usefultools/view/navication"
	"useful-tools/app/usefultools/view/systray"
	viewWidget "useful-tools/app/usefultools/view/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Window struct {
	title       *widget.Label
	intro       *widget.Label
	topWindow   fyne.Window
	application fyne.App
	menu        adapter.Menu
	navigation  adapter.Navigation
	systray     adapter.Systray
	window      fyne.Window
	id          string
	version     string
}

func NewWindow(option model.RunOptions) *Window {
	window := &Window{
		menu:       menu.NewNormal(),
		navigation: navication.NewNormal(),
		systray:    systray.NewNormal(),
		title:      widget.NewLabel(option.DefaultLabel),
		intro:      widget.NewLabel(option.DefaultLabel),
		id:         option.Id,
		version:    option.Version,
	}
	window.intro.Wrapping = fyne.TextWrapWord
	return window
}

func (w *Window) Run() {
	w.application = app.NewWithID(w.id)
	w.application.SetIcon(resource.AppLogo)
	w.window = w.application.NewWindow(w.appTitle())
	if desk, ok := w.application.(desktop.App); ok {
		desk.SetSystemTrayMenu(w.systray.CreateTrayMenu(w.window))
	}
	w.setLifeCycle(w.application)
	content := container.NewStack()

	setPage := func(t adapter.Page) {
		if fyne.CurrentDevice().IsMobile() {
			child := w.application.NewWindow(t.GetTitle())
			w.topWindow = child
			child.SetContent(t.Screen(w.topWindow, constant.ViewMode(w.application.Preferences().Int(constant.NavStatePreferenceProMode))))
			child.Show()
			child.SetOnClosed(func() {
				w.topWindow = w.window
			})
			return
		}

		w.title.SetText(t.GetTitle())
		w.intro.SetText(t.GetIntro())
		content.Objects = []fyne.CanvasObject{t.Screen(w.window, constant.ViewMode(w.application.Preferences().Int(constant.NavStatePreferenceProMode)))}
		content.Refresh()
	}

	w.window.SetMainMenu(w.menu.CreateMenu(w.application, w.window, w.navigation.Tutorials(), setPage, w.ClearCache))
	w.window.SetMaster()

	headerLeft := viewWidget.MakeCellSize(10, 10)
	headerRight := viewWidget.MakeCellSize(10, 10)
	header := container.NewBorder(
		nil, nil, headerLeft, headerRight,
		container.NewVBox(w.title, widget.NewSeparator(), w.intro),
	)
	tutorial := container.NewBorder(header, nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.window.SetContent(w.navigation.CreateNavigation(setPage, false))
	} else {
		split := container.NewHSplit(w.navigation.CreateNavigation(setPage, true), tutorial)
		split.Offset = 0.2
		w.window.SetContent(split)
	}
	w.window.MainMenu()
	w.window.Resize(fyne.NewSize(1100, 690))
	w.window.CenterOnScreen()
	w.window.ShowAndRun()
}

func (w *Window) appTitle() string {
	return fmt.Sprintf("useful-tools v%v", w.version)
}

func (w *Window) setLifeCycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(controller.OnStarted)
	a.Lifecycle().SetOnStopped(controller.OnStopped)
	a.Lifecycle().SetOnEnteredForeground(controller.OnEnteredForeground)
	a.Lifecycle().SetOnExitedForeground(controller.OnExitedForeground)
}

func (w *Window) ClearCache() {
	w.navigation.ClearCache()
}
