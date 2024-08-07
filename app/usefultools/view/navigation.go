package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"useful-tools/app/usefultools/view/constant"
	"useful-tools/app/usefultools/view/dnsquery"
	"useful-tools/app/usefultools/view/draft"
	"useful-tools/app/usefultools/view/portcheck"
	"useful-tools/app/usefultools/view/proxycheck"
)

// Page defines the data structure for a tutorial
type Page struct {
	Title, Intro string
	View         func(w fyne.Window, mode constant.ViewMode) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Page{
		"草稿搭子": {"草稿搭子",
			"合理且好用的草稿纸",
			draft.Screen,
			true,
		},
		"代理检测": {"代理检测",
			"多协议代理可用性检测",
			proxycheck.Screen,
			true,
		},
		"端口检测": {"端口检测",
			"TCP/UDP端口检测",
			portcheck.Screen,
			true,
		},
		"DNS查询": {"DNS查询",
			"一个简单的DNS查询",
			dnsquery.Screen,
			true,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"": {"草稿搭子", "代理检测", "端口检测", "DNS查询"},
	}
)

func createNavigation(setPage func(page Page), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedTutorial(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := Tutorials[uid]; ok {
				if unsupportedTutorial(t) {
					return
				}
				a.Preferences().SetString(constant.NavStatePreferenceCurrentPage, uid)
				setPage(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(constant.NavStatePreferenceCurrentPage, "cgdz")
		tree.Select(currentPref)
	}

	return container.NewBorder(nil, nil, nil, nil, tree)
}

func unsupportedTutorial(t Page) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}
