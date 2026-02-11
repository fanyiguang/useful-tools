package page

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/controller"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/app/usefultools/view/constant"
	viewWidget "useful-tools/app/usefultools/view/widget"
	"useful-tools/helper/Go"
)

var _ adapter.Page = (*ProxyCheck)(nil)

type ProxyCheck struct {
	BasePage
	logics       *controller.ProxyCheck
	view         *widget.Entry
	scroll       *container.Scroll
	latestParams string
	loadingStop  chan struct{}
}

func NewProxyCheck() *ProxyCheck {
	return &ProxyCheck{
		BasePage: BasePage{
			ID:         constant.PageIDProxyCheck,
			TitleKey:   i18n.KeyPageProxyTitle,
			IntroKey:   i18n.KeyPageProxyIntro,
			SupportWeb: true,
		},
		logics: controller.NewProxyCheck(),
	}
}

func (p *ProxyCheck) Screen(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject {
	return container.NewHSplit(p.leftProxyCheck(mode), p.rightProxyCheck(w))
}

func (p *ProxyCheck) leftProxyCheck(mode constant.ViewMode) fyne.CanvasObject {
	logrus.Infof("proxy check mode: %d", mode)
	switch mode {
	case constant.ViewModePro:
		return p.proView()
	default:
		right := viewWidget.MakeCellSize(10, 10)
		left := viewWidget.MakeCellSize(10, 10)
		top := viewWidget.MakeCellSize(10, 10)
		bottom := viewWidget.MakeCellSize(10, 10)
		return container.NewBorder(top, bottom, left, right, p.checkFrom())
	}
}

func (p *ProxyCheck) proView() fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	multi := widget.NewMultiLineEntryEx(nil, nil, nil, p.logics.FormatJson)
	multi.Wrapping = fyne.TextWrapWord
	if p.logics.PreModeInput() != "" {
		multi.SetText(p.logics.PreModeInput())
	} else {
		multi.SetText(p.logics.ProTemplate())
	}
	multi.OnChanged = func(s string) {
		logrus.Infof("proxy: %s", s)
		p.logics.SetPreModeInput(s)
	}

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonClear),
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("pre clear proxy: %s", multi.Text)
			multi.SetText("")
		},
	}, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonCheck),
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			p.latestParams = multi.Text
			p.startLoading()
			go func() {
				text := multi.Text
				logrus.Infof("pre check proxy: %s", text)
				response, err := p.logics.ProCheckProxy(text, fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceHideBody))
				if p.latestParams == text {
					p.stopLoading()
					if err != nil {
						logrus.Warnf("pre check proxy failed: %s", err)
						p.view.SetText(err.Error())
						return
					}
					logrus.Infof("pre check proxy result: %s", response)
					p.view.SetText(response)
				}
			}()
		},
	})
	border := container.NewBorder(nil, box, nil, nil, container.NewVScroll(multi))
	return container.NewBorder(top, bottom, left, right, border)
}

func (p *ProxyCheck) rightProxyCheck(w fyne.Window) fyne.CanvasObject {
	right := viewWidget.MakeCellSize(10, 10)
	left := viewWidget.MakeCellSize(10, 10)
	top := viewWidget.MakeCellSize(10, 10)
	bottom := viewWidget.MakeCellSize(10, 10)

	p.view = widget.NewMultiLineEntry()
	p.view.Wrapping = fyne.TextWrapWord
	p.view.TextStyle = fyne.TextStyle{Bold: true}
	if p.logics.ViewText() != "" {
		p.view.SetText(p.logics.ViewText())
	} else {
		p.view.PlaceHolder = i18n.T(i18n.KeyProxyResultPlaceholder)
	}
	p.view.OnChanged = func(s string) {
		logrus.Infof("proxy check result: %s", s)
		p.logics.SetViewText(s)
	}
	//view.Disable()

	box := container.NewGridWithColumns(2, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonClear),
		Icon:       theme.Icon(theme.IconNameContentClear),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("proxy check view check clear: %s", p.view.Text)
			p.view.SetText("")
		},
	}, &widget.Button{
		Text:       i18n.T(i18n.KeyButtonCopy),
		Icon:       theme.Icon(theme.IconNameContentCopy),
		Importance: widget.MediumImportance,
		OnTapped: func() {
			logrus.Infof("proxy check view check copy: %s", p.view.Text)
			w.Clipboard().SetContent(strings.TrimSpace(p.view.Text))
		},
	})
	p.scroll = container.NewVScroll(p.view)
	border := container.NewBorder(nil, box, nil, nil, p.scroll)
	return container.NewBorder(top, bottom, left, right, border)
}

func (p *ProxyCheck) checkFrom() fyne.CanvasObject {
	proxyTypeSelect := widget.NewSelect(p.logics.SupportProxyTypeList(), func(s string) {
		logrus.Infof("proxy check type: %s", s)
		p.logics.SetTyp(s)
	})
	if p.logics.Typ() != "" {
		proxyTypeSelect.SetSelected(p.logics.Typ())
	} else {
		proxyTypeSelect.SetSelected("SOCKS5")
	}

	host := widget.NewEntry()
	host.SetPlaceHolder(i18n.T(i18n.KeyProxyHostPlaceholder))
	host.SetText(p.logics.Host())
	host.OnChanged = func(s string) {
		logrus.Infof("proxy check host: %s", s)
		p.logics.SetHost(s)
		host.SetText(s)
	}
	host.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("proxy check validation host: %s", s)
		if s == "" {
			return nil
		}
		if strings.Contains(strings.TrimSpace(s), ".") {
			return nil
		} else {
			return errors.New(i18n.T(i18n.KeyProxyInvalidHostError))
		}
	})

	port := widget.NewEntry()
	port.SetPlaceHolder(i18n.T(i18n.KeyProxyPortPlaceholder))
	port.SetText(p.logics.Port())
	port.OnChanged = func(s string) {
		logrus.Infof("proxy check port: %s", s)
		p.logics.SetPort(s)
		port.SetText(s)
	}
	port.Validator = validation.NewAllStrings(func(s string) error {
		logrus.Infof("proxy check validation domain: %s", s)
		if s == "" {
			return nil
		}
		iPort, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			logrus.Warnf("strconv.Atoi error: %v", err)
			return errors.New(i18n.T(i18n.KeyProxyInvalidPortError))
		}
		if iPort > 65535 || iPort < 0 {
			return errors.New(i18n.T(i18n.KeyProxyInvalidPortRangeError))
		}
		return nil
	})

	username := widget.NewEntry()
	username.SetPlaceHolder(i18n.T(i18n.KeyProxyUsernamePlaceholder))
	username.SetText(p.logics.Username())
	username.OnChanged = func(s string) {
		logrus.Infof("proxy check username: %s", s)
		p.logics.SetUsername(s)
	}

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(i18n.T(i18n.KeyProxyPasswordPlaceholder))
	password.SetText(p.logics.Password())
	password.OnChanged = func(s string) {
		logrus.Infof("proxy check password: %s", s)
		p.logics.SetPassword(s)
	}

	urls := widget.NewMultiLineEntry()
	urls.SetPlaceHolder(i18n.T(i18n.KeyProxyURLPlaceholder))
	urls.SetText(p.logics.Urls())
	urls.OnChanged = func(s string) {
		logrus.Infof("proxy check urls: %s", s)
		p.logics.SetUrls(s)
	}

	var form *widget.StyleForm
	form = &widget.StyleForm{
		Items: []*widget.StyleFormItem{
			{Text: i18n.T(i18n.KeyProxyTypeLabel), Widget: proxyTypeSelect, HintText: i18n.T(i18n.KeyHintSelectRequired)},
			{Text: i18n.T(i18n.KeyProxyHostLabel), Widget: host, HintText: i18n.T(i18n.KeyHintRequired)},
			{Text: i18n.T(i18n.KeyProxyPortLabel), Widget: port, HintText: i18n.T(i18n.KeyHintRequired)},
			{Text: i18n.T(i18n.KeyProxyUsernameLabel), Widget: username, HintText: i18n.T(i18n.KeyHintRequired)},
			{Text: i18n.T(i18n.KeyProxyPasswordLabel), Widget: password, HintText: i18n.T(i18n.KeyHintRequired)},
			{Text: i18n.T(i18n.KeyProxyTestURLLabel), Widget: urls, HintText: i18n.T(i18n.KeyHintOptional)},
		},
		OnCancel: func() {
			logrus.Infof("proxy check page cancelled")
			//proxyTypeSelect.SetSelected("SOCKS5")
			host.SetText("")
			port.SetText("")
			username.SetText("")
			password.SetText("")
			urls.SetText("")
		},
		OnSubmit: func() {
			p.latestParams = fmt.Sprintf("%v%v%v%v%v%v%v", host.Text, port.Text, username.Text, password.Text, proxyTypeSelect.Selected, urls.Text, fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceHideBody))
			p.startLoading()
			Go.RelativelySafeGo(func() {
				hostText := host.Text
				portText := port.Text
				usernameText := username.Text
				passwordText := password.Text
				selectText := proxyTypeSelect.Selected
				urlText := urls.Text
				isHideBody := fyne.CurrentApp().Preferences().Bool(constant.NavStatePreferenceHideBody)
				logrus.Infof("proxy check page submitted")
				response, err := p.logics.NormalCheckProxy(hostText, portText, usernameText, passwordText, selectText, urlText, isHideBody)
				if p.latestParams == fmt.Sprintf("%v%v%v%v%v%v%v", hostText, portText, usernameText, passwordText, selectText, urlText, isHideBody) {
					p.stopLoading()
					if err != nil {
						logrus.Errorf("proxy check error: %v", err)
						p.view.SetText(err.Error())
						return
					}
					logrus.Infof("proxy check result: %s", response)
					p.view.SetText(response)
				}
			})
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

func (p *ProxyCheck) startLoading() {
	if p.view == nil {
		return
	}
	p.stopLoading()
	stop := make(chan struct{})
	p.loadingStop = stop
	base := i18n.T(i18n.KeyProxyCheckInProgress)
	p.view.SetText(base + "...")
	go func() {
		ticker := time.NewTicker(350 * time.Millisecond)
		defer ticker.Stop()
		dots := 0
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				dots = (dots + 1) % 4
				msg := base + strings.Repeat(".", dots)
				fyne.Do(func() {
					if p.loadingStop == stop {
						p.view.SetText(msg)
					}
				})
			}
		}
	}()
}

func (p *ProxyCheck) stopLoading() {
	if p.loadingStop == nil {
		return
	}
	close(p.loadingStop)
	p.loadingStop = nil
}

func (p *ProxyCheck) ClearCache() {
	p.logics.ClearCache()
}
