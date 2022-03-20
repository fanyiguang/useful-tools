package common

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

var (
	ConvenientModeMenu *walk.Action
)

func WinCenter(handle win.HWND) {
	var rect win.RECT
	if win.GetWindowRect(handle, &rect) {
		scrWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
		scrHeight := win.GetSystemMetrics(win.SM_CYSCREEN)
		windowWidth := rect.Right - rect.Left
		windowHeight := rect.Bottom - rect.Top
		rect.Right = windowWidth // 设置窗口宽高
		rect.Bottom = windowHeight
		rect.Left = (scrWidth - rect.Right) / 2
		rect.Top = (scrHeight - rect.Bottom) / 2
		// 设置窗体位置 （居中显示）
		win.MoveWindow(handle, rect.Left, rect.Top, rect.Right, rect.Bottom, true) // 居中
	}
}

func WinReSize(handle win.HWND, width int32, height int32) {
	var rect win.RECT
	if win.GetWindowRect(handle, &rect) {
		win.MoveWindow(handle, rect.Left, rect.Top, width, height, true)
	}
}
