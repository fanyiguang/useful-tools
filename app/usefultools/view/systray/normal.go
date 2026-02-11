package systray

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/sirupsen/logrus"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/i18n"
)

var _ adapter.Systray = (*Normal)(nil)

type Normal struct {
}

func NewNormal() *Normal {
	return &Normal{}
}

func (n *Normal) CreateTrayMenu(window fyne.Window) *fyne.Menu {
	h := fyne.NewMenuItem(i18n.T(i18n.KeySystrayOpen), func() {
		logrus.Infof("open tools")
		window.Show()
	})
	h.Icon = theme.HomeIcon()
	menu := fyne.NewMenu("systray", h)
	return menu
}
