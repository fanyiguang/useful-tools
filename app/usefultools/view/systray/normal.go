package systray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/sirupsen/logrus"
	"useful-tools/app/usefultools/adapter"
)

var _ adapter.Systray = (*Normal)(nil)

type Normal struct {
}

func NewNormal() *Normal {
	return &Normal{}
}

func (n *Normal) CreateTrayMenu(window fyne.Window) *fyne.Menu {
	h := fyne.NewMenuItem("打开工具", func() {
		logrus.Infof("open tools")
		window.Show()
	})
	h.Icon = theme.HomeIcon()
	menu := fyne.NewMenu("systray", h)
	return menu
}
