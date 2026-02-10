package navication

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/i18n"
	"useful-tools/app/usefultools/view/constant"
	"useful-tools/app/usefultools/view/page"
)

// Page defines the data structure for a tutorial
type Page struct {
	Title, Intro string
	View         func(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject
	SupportWeb   bool
}

type Normal struct {
	pages         []adapter.Page
	tutorials     map[string]adapter.Page
	tutorialIndex map[string][]string
	lastSelected  string
	tree          *widget.Tree
}

func NewNormal() *Normal {
	n := &Normal{
		pages: []adapter.Page{
			page.NewDraft(),
			page.NewProxyCheck(),
			page.NewPortCheck(),
			page.NewDnsQuery(),
			page.NewAesConversion(),
			page.NewJsonTools(),
			page.NewFormatConversion(),
		},
		tutorials:     make(map[string]adapter.Page),
		tutorialIndex: make(map[string][]string),
	}

	var titles []string
	for _, p := range n.pages {
		n.tutorials[p.GetID()] = p
		titles = append(titles, p.GetID())
	}

	n.tutorialIndex[""] = titles
	return n
}

func (n *Normal) CreateNavigation(setPage func(page adapter.Page), loadPrevious bool) fyne.CanvasObject {
	app := fyne.CurrentApp()
	var tree *widget.Tree
	tree = &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return n.tutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := n.tutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			icon := widget.NewIcon(nil)
			label := widget.NewLabel("Collection Widgets")
			return container.NewHBox(icon, label)
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := n.tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			node := obj.(*fyne.Container)
			icon := node.Objects[0].(*widget.Icon)
			label := node.Objects[1].(*widget.Label)
			icon.SetResource(navIconForID(t.GetID()))
			label.SetText(t.GetTitle())
			if unsupportedTutorial(t) {
				label.TextStyle = fyne.TextStyle{Italic: true}
			} else {
				label.TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := n.tutorials[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
				app.Preferences().SetString(constant.NavStatePreferenceCurrentPage, uid)
				setPage(t)
				if n.lastSelected != "" && n.lastSelected != uid {
					tree.RefreshItem(n.lastSelected)
				}
				n.lastSelected = uid
				tree.RefreshItem(uid)
			}
		},
	}
	n.tree = tree

	if loadPrevious {
		currentPref := app.Preferences().StringWithFallback(constant.NavStatePreferenceCurrentPage, constant.PageIDDraft)
		if _, ok := n.tutorials[currentPref]; !ok {
			if mapped, ok := legacyTitleToID[currentPref]; ok {
				currentPref = mapped
			}
		}
		app.Preferences().SetString(constant.NavStatePreferenceCurrentPage, currentPref)
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func unsupportedTutorial(t adapter.Page) bool {
	return !t.GetSupportWeb() && fyne.CurrentDevice().IsBrowser()
}

var legacyTitleToID = func() map[string]string {
	mapping := make(map[string]string)
	add := func(id string, key i18n.Key) {
		for _, title := range i18n.All(key) {
			if title != "" {
				mapping[title] = id
			}
		}
	}
	add(constant.PageIDDraft, i18n.KeyPageDraftTitle)
	add(constant.PageIDProxyCheck, i18n.KeyPageProxyTitle)
	add(constant.PageIDPortCheck, i18n.KeyPagePortTitle)
	add(constant.PageIDDnsQuery, i18n.KeyPageDnsTitle)
	add(constant.PageIDAesConvert, i18n.KeyPageAesTitle)
	add(constant.PageIDJsonTools, i18n.KeyPageJsonTitle)
	add(constant.PageIDFormatConversion, i18n.KeyPageFormatConversionTitle)
	return mapping
}()

func navIconForID(id string) fyne.Resource {
	switch id {
	case constant.PageIDDraft:
		return theme.DocumentCreateIcon()
	case constant.PageIDProxyCheck:
		return theme.SettingsIcon()
	case constant.PageIDPortCheck:
		return theme.ComputerIcon()
	case constant.PageIDDnsQuery:
		return theme.SearchIcon()
	case constant.PageIDAesConvert:
		return theme.ContentRedoIcon()
	case constant.PageIDJsonTools:
		return theme.FileTextIcon()
	case constant.PageIDFormatConversion:
		return theme.ListIcon()
	default:
		return theme.InfoIcon()
	}
}

func (n *Normal) Tutorials() map[string]adapter.Page {
	return n.tutorials
}

func (n *Normal) Refresh() {
	if n.tree != nil {
		n.tree.Refresh()
	}
}

func MapLegacyTitleToID(title string) (string, bool) {
	id, ok := legacyTitleToID[title]
	return id, ok
}

func (n *Normal) ClearCache() {
	for _, p := range n.pages {
		p.ClearCache()
	}
}
