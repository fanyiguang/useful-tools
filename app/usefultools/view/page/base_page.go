package page

import "useful-tools/app/usefultools/i18n"

type BasePage struct {
	ID         string
	TitleKey   i18n.Key
	IntroKey   i18n.Key
	SupportWeb bool
}

func (b *BasePage) GetID() string {
	return b.ID
}

func (b *BasePage) GetTitle() string {
	return i18n.T(b.TitleKey)
}

func (b *BasePage) GetIntro() string {
	return i18n.T(b.IntroKey)
}

func (b *BasePage) GetSupportWeb() bool {
	return b.SupportWeb
}
