package page

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
	"time"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
	"useful-tools/helper/Go"
)

var _ adapter.Page = (*DnsQuery)(nil)

type DnsQuery struct {
	BasePage
	logics       *controller.DnsQuery
	view         *widget.Entry
	scroll       *container.Scroll
	latestParams string
}

func NewDnsQuery() *DnsQuery {
	return &DnsQuery{
		BasePage: BasePage{
			ID:         constant.PageIDDnsQuery,
			TitleKey:   i18n.KeyPageDnsTitle,
			IntroKey:   i18n.KeyPageDnsIntro,
			SupportWeb: true,
		},
		logics: controller.NewDnsQuery(),
	}
}

func (d *DnsQuery) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(d.leftScreen(mode), d.rightScreen(w))
}

func (d *DnsQuery) leftScreen(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("dns query mode: %d", mode)
	switch mode {
	case constant.ViewModePro:
		return d.proView()
	default:
		right := viewWidget.MakeCellSize(10, 10)
		left := viewWidget.MakeCellSize(10, 10)
		top := viewWidget.MakeCellSize(10, 10)
		bottom := viewWidget.MakeCellSize(10, 10)
		return container.NewBorder(top, bottom, left, right, d.from())
	}
}

func (d *DnsQuery) proView() fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	multi := widget.NewMultiLineEntryEx(nil, nil, nil, d.logics.FormatJson)
	multi.Wrapping = fyne.TextWrapWord
	if d.logics.PreModeInput() != "" {
		multi.SetText(d.logics.PreModeInput())
	} else {
		multi.SetText(d.logics.PreTemplate())
	}
	multi.OnChanged = func(s string) {
		logrus.Infof("dns query: %s", s)
		d.logics.SetPreModeInput(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonClear),
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("pre clear dns query: %s", multi.Text)
			multi.SetText("")
		},
	}, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonCheck),
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			d.latestParams = multi.Text
			Go.RelativelySafeGo(func() {
				text := multi.Text
				logrus.Infof("pre check dns query: %s", text)
				response, err := d.logics.ProQuery(text)
				if d.latestParams == multi.Text {
					if err != nil {
						logrus.Errorf("pro dns query error: %v", err)
						d.view.SetText(d.logFormat(gjson.Get(text, "server").String(), gjson.Get(text, "domain").String(), err.Error()) + d.view.Text)
					} else {
						logrus.Infof("pro dns query result: %v", response)
						d.view.SetText(d.logFormat(gjson.Get(text, "server").String(), gjson.Get(text, "domain").String(), strings.Join(response, " ")) + text)
					}
				}
			})
		},
	})
	border := container.NewBorder(nil, box, nil, nil, multi)
	return container.NewBorder(top, bottom, left, right, border)
}

func (d *DnsQuery) rightScreen(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	d.view = widget.NewMultiLineEntry()
	d.view.Wrapping = fyne.TextWrapWord
	d.view.Scroll = container.ScrollVerticalOnly
	d.view.TextStyle = fyne.TextStyle{Bold: true}
	if d.logics.ViewText() != "" {
		d.view.SetText(d.logics.ViewText())
	} else {
		d.view.PlaceHolder = i18n.T(i18n.KeyDnsResultPlaceholder)
	}
	d.view.OnChanged = func(s string) {
		logrus.Infof("dns query result: %s", s)
		d.logics.SetViewText(s)
	}
	//view.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonClear),
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("dns query view check clear: %s", d.view.Text)
			d.view.SetText("")
		},
	}, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonCopy),
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("dns query view check copy: %s", d.view.Text)
			viewWidget.CopyToClipboard(w, d.view.Text)
		},
	})
	d.scroll = container.NewVScroll(d.view)
	border := container.NewBorder(nil, box, nil, nil, d.scroll)
	return container.NewBorder(top, bottom, left, right, border)
}

func (d *DnsQuery) from() fyne.CanvasObject {
	host := widget.NewEntry()
	host.SetPlaceHolder(i18n.T(i18n.KeyDnsHostPlaceholder))
	hostValue := d.logics.Host()
	if hostValue == "" || i18n.Matches(i18n.KeyDefault, hostValue) {
		host.SetText(i18n.T(i18n.KeyDefault))
	} else {
		host.SetText(hostValue)
	}
	host.OnChanged = func(s string) {
		logrus.Infof("dns query host: %s", s)
		d.logics.SetHost(s)
		host.SetText(s)
	}

	domain := widget.NewEntry()
	domain.SetPlaceHolder(i18n.T(i18n.KeyDnsDomainPlaceholder))
	domain.SetText(d.logics.Domain())
	domain.OnChanged = func(s string) {
		logrus.Infof("dns query domain: %s", s)
		d.logics.SetDomain(s)
		domain.SetText(s)
	}
	domain.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("dns query validation domain: %s", s)
		if s == "" {
			return nil
		}
		if !strings.Contains(s, ".") {
			return errors.New(i18n.T(i18n.KeyDnsInvalidDomainError))
		}
		return nil
	})

	var form *widget.StyleForm
	form = &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: i18n.T(i18n.KeyDnsHostLabel), Widget: host, HintText: i18n.T(i18n.KeyHintRequired)},
			{Text: i18n.T(i18n.KeyDnsDomainLabel), Widget: domain, HintText: i18n.T(i18n.KeyHintRequired)},
		},
		OnCancel: func() {
			logrus.Infof("dns query page cancelled")
			host.SetText("")
			domain.SetText("")
		},
		OnSubmit: func() {
			d.latestParams = fmt.Sprintf("%s%s", host.Text, domain.Text)
			go func() {
				logrus.Infof("dns query page submitted")
				host := host.Text
				domain := domain.Text
				response, err := d.logics.NormalDns(host, domain)
				if d.latestParams == fmt.Sprintf("%s%s", host, domain) {
					if err != nil {
						logrus.Errorf("dns query error: %v", err)
						d.view.SetText(d.logFormat(host, domain, err.Error()) + d.view.Text)
					} else {
						logrus.Infof("dns query result: %v", response)
						d.view.SetText(d.logFormat(host, domain, strings.Join(response, " ")) + d.view.Text)
					}
				}
			}()
		},
		SubmitText: i18n.T(i18n.KeyButtonCheck),
		CancelText: i18n.T(i18n.KeyButtonClear),
		ButtonLayout: func(cancel *widget.Button, submit *widget.Button) *fyne.Container {
			return container.NewGridWithColumns(2, cancel, submit)
		},
		ContentLayout: func(input *fyne.Container, button *fyne.Container) *fyne.Container {
			return container.NewBorder(nil, button, nil, nil, input)
		},
	}
	return form
}

func (d *DnsQuery) logFormat(server, host, msg string) string {
	return fmt.Sprintf("[%v] %v(%v) => %v\n\n", time.Now().Format(`06-01-02 15:04:05`), host, server, msg)
}

func (d *DnsQuery) ClearCache() {
	d.logics.ClearCache()
}
