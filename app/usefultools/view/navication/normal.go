package navication

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"useful-tools/app/usefultools/adapter"
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
		},
		tutorials:     make(map[string]adapter.Page),
		tutorialIndex: make(map[string][]string),
	}

	var titles []string
	for _, p := range n.pages {
		n.tutorials[p.GetTitle()] = p
		titles = append(titles, p.GetTitle())
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
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := n.tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.GetTitle())
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
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

	if loadPrevious {
		currentPref := app.Preferences().StringWithFallback(constant.NavStatePreferenceCurrentPage, "草稿搭子")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func unsupportedTutorial(t adapter.Page) bool {
	return !t.GetSupportWeb() && fyne.CurrentDevice().IsBrowser()
}

func (n *Normal) Tutorials() map[string]adapter.Page {
	return n.tutorials
}

func (n *Normal) ClearCache() {
	for _, p := range n.pages {
		p.ClearCache()
	}
}
