package systray

import (
	"fmt"
	"useful-tools/helper/Go"
	"useful-tools/pkg/wlog"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/pkg/errors"
)

func New(closeAction func(), openAction func()) (systrayMainWindow *walk.MainWindow) {
	systrayMainWindow = new(walk.MainWindow)
	mainWindow := declarative.MainWindow{
		AssignTo: &systrayMainWindow,
		Size:     declarative.Size{Width: 1, Height: 1},
		MinSize:  declarative.Size{Width: 1, Height: 1},
		MaxSize:  declarative.Size{Width: 1, Height: 1},
		Layout:   declarative.VBox{},
	}
	_ = mainWindow.Create()
	err := initSystray(systrayMainWindow, closeAction, openAction)
	if err != nil {
		wlog.Warm("initSystray failed: %v", err)
	}
	dwExStyle := win.GetWindowLong(systrayMainWindow.Handle(), win.GWL_STYLE)
	dwExStyle &= ^(win.WS_VISIBLE)
	dwExStyle |= win.WS_EX_TOOLWINDOW
	dwExStyle &= ^(win.WS_EX_APPWINDOW)
	win.SetWindowLong(systrayMainWindow.Handle(), win.GWL_STYLE, dwExStyle)
	win.ShowWindow(systrayMainWindow.Handle(), win.SW_SHOW)
	win.ShowWindow(systrayMainWindow.Handle(), win.SW_HIDE)
	Go.Go(func() {
		win.ShowWindow(systrayMainWindow.Handle(), win.SW_HIDE)
	})
	return
}

func initSystray(form walk.Form, closeActionHandle func(), openActionHandle func()) (err error) {
	ni, err := walk.NewNotifyIcon(form)
	if err != nil {
		err = errors.Wrap(err, "walk.NewNotifyIcon")
		return
	}

	icon, err := walk.Resources.Image("./logo.png")
	fmt.Println(icon, err)
	if icon != nil {
		err = ni.SetIcon(icon)
	}

	openAction := walk.NewAction()
	_ = openAction.SetText("打开面板")
	openAction.Triggered().Attach(func() {
		fmt.Println("点击打开面板")
		openActionHandle()
	})
	_ = ni.ContextMenu().Actions().Add(openAction)

	closeAction := walk.NewAction()
	_ = closeAction.SetText("关闭程序")
	closeAction.Triggered().Attach(func() {
		fmt.Println("点击关闭程序")
		closeActionHandle()
	})
	_ = ni.ContextMenu().Actions().Add(closeAction)
	_ = ni.SetVisible(true)
	return
}
