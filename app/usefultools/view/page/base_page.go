package page

type BasePage struct {
	Title      string
	Intro      string
	SupportWeb bool
}

func (b *BasePage) GetTitle() string {
	return b.Title
}

func (b *BasePage) GetIntro() string {
	return b.Intro
}

func (b *BasePage) GetSupportWeb() bool {
	return b.SupportWeb
}
