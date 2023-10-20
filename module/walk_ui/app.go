package walkUI

import (
	"bytes"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"useful-tools/common/config"
	"useful-tools/module/logic/app"
	"useful-tools/module/walk_ui/aes"
	"useful-tools/module/walk_ui/common"
	"useful-tools/module/walk_ui/dns"
	"useful-tools/module/walk_ui/proxy"
	"useful-tools/module/walk_ui/systray"
	"useful-tools/module/walk_ui/tcp_udp"
)

type PageFactoryFunc func(parent walk.Container, IsConvenientMode bool) (common.Page, error)

type MultiPageMainWindow struct {
	*walk.MainWindow
	logicControl                *app.App
	navTB                       *walk.ToolBar
	pageCom                     *walk.Composite
	action2NewPage              map[*walk.Action]PageFactoryFunc
	pageActions                 []*walk.Action
	currentAction               *walk.Action
	currentPage                 common.Page
	currentPageChangedPublisher walk.EventPublisher
	systrayMainWindow           *walk.MainWindow
}

type AppMainWindow struct {
	*MultiPageMainWindow
}

type PageConfig struct {
	Title   string
	Image   string
	NewPage PageFactoryFunc
}

type MultiPageMainWindowConfig struct {
	Name                 string
	Enabled              Property
	Visible              Property
	Font                 Font
	MinSize              Size
	MaxSize              Size
	ContextMenuItems     []MenuItem
	OnKeyDown            walk.KeyEventHandler
	OnKeyPress           walk.KeyEventHandler
	OnKeyUp              walk.KeyEventHandler
	OnMouseDown          walk.MouseEventHandler
	OnMouseMove          walk.MouseEventHandler
	OnMouseUp            walk.MouseEventHandler
	OnSizeChanged        walk.EventHandler
	OnCurrentPageChanged walk.EventHandler
	Title                string
	Size                 Size
	MenuItems            []MenuItem
	ToolBar              ToolBar
	PageCfgs             []PageConfig
}

func (mpmw *MultiPageMainWindow) CurrentAction() *walk.Action {
	return mpmw.currentAction
}

func (mpmw *MultiPageMainWindow) CurrentPage() common.Page {
	return mpmw.currentPage
}

func (mpmw *MultiPageMainWindow) CurrentPageTitle() string {
	if mpmw.currentAction == nil {
		return ""
	}

	return mpmw.currentAction.Text()
}

func (mpmw *MultiPageMainWindow) CurrentPageChanged() *walk.Event {
	return mpmw.currentPageChangedPublisher.Event()
}

func (mpmw *MultiPageMainWindow) newPageAction(title, image string, newPage PageFactoryFunc) (*walk.Action, error) {
	img, err := walk.Resources.Bitmap(image)
	if err != nil {
		return nil, err
	}

	action := walk.NewAction()
	action.SetCheckable(true)
	action.SetExclusive(true)
	action.SetImage(img)
	action.SetText(title)

	mpmw.action2NewPage[action] = newPage

	action.Triggered().Attach(func() {
		mpmw.SetCurrentAction(action)
	})

	return action, nil
}

func (mpmw *MultiPageMainWindow) SetCurrentAction(action *walk.Action) error {
	defer func() {
		if !mpmw.pageCom.IsDisposed() {
			mpmw.pageCom.RestoreState()
		}
	}()

	mpmw.SetFocus()

	if prevPage := mpmw.currentPage; prevPage != nil {
		mpmw.pageCom.SaveState()
		prevPage.SetVisible(false)
		prevPage.(walk.Widget).SetParent(nil)
		prevPage.Dispose()
	}

	newPage := mpmw.action2NewPage[action]

	page, err := newPage(mpmw.pageCom, ConvenientModeMenu.Checked())
	if err != nil {
		return err
	}

	action.SetChecked(true)

	mpmw.currentPage = page
	mpmw.currentAction = action

	mpmw.currentPageChangedPublisher.Publish()

	return nil
}

func (mpmw *MultiPageMainWindow) updateNavigationToolBar() error {
	mpmw.navTB.SetSuspended(true)
	defer mpmw.navTB.SetSuspended(false)

	actions := mpmw.navTB.Actions()

	if err := actions.Clear(); err != nil {
		return err
	}

	for _, action := range mpmw.pageActions {
		if err := actions.Add(action); err != nil {
			return err
		}
	}

	if mpmw.currentAction != nil {
		if !actions.Contains(mpmw.currentAction) {
			for _, action := range mpmw.pageActions {
				if action != mpmw.currentAction {
					if err := mpmw.SetCurrentAction(action); err != nil {
						return err
					}

					break
				}
			}
		}
	}

	return nil
}

func NewMultiPageMainWindow(cfg *MultiPageMainWindowConfig) (*MultiPageMainWindow, error) {
	mpmw := &MultiPageMainWindow{
		action2NewPage: make(map[*walk.Action]PageFactoryFunc),
		logicControl:   app.New(),
	}

	if err := (MainWindow{
		AssignTo:         &mpmw.MainWindow,
		Name:             cfg.Name,
		Title:            cfg.Title,
		Enabled:          cfg.Enabled,
		Visible:          cfg.Visible,
		Font:             cfg.Font,
		MinSize:          cfg.MinSize,
		MaxSize:          cfg.MaxSize,
		MenuItems:        cfg.MenuItems,
		ToolBar:          cfg.ToolBar,
		ContextMenuItems: cfg.ContextMenuItems,
		OnKeyDown:        cfg.OnKeyDown,
		OnKeyPress:       cfg.OnKeyPress,
		OnKeyUp:          cfg.OnKeyUp,
		OnMouseDown:      cfg.OnMouseDown,
		OnMouseMove:      cfg.OnMouseMove,
		OnMouseUp:        cfg.OnMouseUp,
		OnSizeChanged:    cfg.OnSizeChanged,
		Layout:           HBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			ScrollView{
				HorizontalFixed: true,
				Layout:          VBox{MarginsZero: true},
				Children: []Widget{
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							ToolBar{
								AssignTo:    &mpmw.navTB,
								Orientation: Vertical,
								ButtonStyle: ToolBarButtonImageAboveText,
								MaxTextRows: 2,
							},
						},
					},
				},
			},
			Composite{
				AssignTo: &mpmw.pageCom,
				Name:     "pageCom",
				Layout:   HBox{MarginsZero: true, SpacingZero: true},
			},
			CheckBox{
				Name:    "openHiddenCB",
				Text:    "Open Hidden",
				Checked: false,
				Visible: false,
			},
		},
	}).Create(); err != nil {
		return nil, err
	}

	initMenuItems()

	var handleClosing int
	handleClosing = mpmw.MainWindow.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		mpmw.MainWindow.Closing().Detach(handleClosing)
	})

	succeeded := false
	defer func() {
		if !succeeded {
			mpmw.Dispose()
		}
	}()

	for _, pc := range cfg.PageCfgs {
		action, err := mpmw.newPageAction(pc.Title, pc.Image, pc.NewPage)
		if err != nil {
			return nil, err
		}

		mpmw.pageActions = append(mpmw.pageActions, action)
	}

	if err := mpmw.updateNavigationToolBar(); err != nil {
		return nil, err
	}

	if len(mpmw.pageActions) > 0 {
		if err := mpmw.SetCurrentAction(mpmw.pageActions[0]); err != nil {
			return nil, err
		}
	}

	if cfg.OnCurrentPageChanged != nil {
		mpmw.CurrentPageChanged().Attach(cfg.OnCurrentPageChanged)
	}

	icon, err := walk.Resources.Image("icon.png")
	if err == nil {
		_ = mpmw.SetIcon(icon)
	}
	common.WinCenter(mpmw.Handle())
	//win.RemoveMenu(win.GetSystemMenu(mpmw.Handle(), false), win.SC_SIZE, win.MF_BYCOMMAND)  //禁止改变窗体大小
	//currStyle := win.GetWindowLong(mpmw.MainWindow.Handle(), win.GWL_STYLE)
	//win.SetWindowLong(mpmw.MainWindow.Handle(), win.GWL_STYLE, currStyle&^win.WS_MAXIMIZEBOX) //禁用最大化
	//mpmw.MainWindow.Activating().Attach(func() {
	//	common.WinReSize(mpmw.MainWindow.Handle(), 912, 592)
	//})

	mpmw.systrayMainWindow = systray.New(func() {
		if mpmw != nil {
			_ = mpmw.systrayMainWindow.Close()
		}
	}, func() {
		if mpmw != nil {
			mpmw.Show()
		}
	})
	mpmw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if reason == walk.CloseReasonUnknown {
			*canceled = true // 阻止系统关闭事件
			mpmw.Hide()      // 改为隐藏
		}
	})
	succeeded = true
	return mpmw, nil
}

func (mw *AppMainWindow) updateTitle(prefix string) {
	var buf bytes.Buffer
	if prefix != "" {
		buf.WriteString(prefix)
		buf.WriteString(" - ")
	}
	buf.WriteString("useful-tools ")
	buf.WriteString(config.Version)
	mw.SetTitle(buf.String())
}

func (mw *AppMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw,
		"About Walk Multiple Pages Example",
		"An example that demonstrates a main window that supports multiple pages.",
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

func (mw *AppMainWindow) openAction_Triggered() {
	walk.MsgBox(mw, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
}

func New() *walk.MainWindow {
	mw := new(AppMainWindow)
	cfg := &MultiPageMainWindowConfig{
		Name:      "mainWindow",
		MinSize:   Size{1000, 550},
		MaxSize:   Size{1000, 550},
		Size:      Size{1000, 550},
		MenuItems: MenuItems(mw),
		OnCurrentPageChanged: func() {
			mw.updateTitle(mw.CurrentPageTitle())
		},
		PageCfgs: []PageConfig{
			{"代理检测", "proxy.png", proxy.NewPage},
			{"端口检测", "tcp_udp.png", tcp_udp.NewPage},
			{"DNS检测", "dns.png", dns.NewPage},
			{"AES转换", "aes.png", aes.NewPage},
		},
	}

	mpmw, err := NewMultiPageMainWindow(cfg)
	if err != nil {
		panic(err)
	}

	//addRecentFileActions := func(texts ...string) {
	//	for _, text := range texts {
	//		a := walk.NewAction()
	//		_ = a.SetText(text)
	//		a.Triggered().Attach(mw.openAction_Triggered)
	//		_ = modeMenu.Actions().Add(a)
	//	}
	//}
	//addRecentFileActions("Foo", "Bar")

	mw.MultiPageMainWindow = mpmw
	mw.updateTitle(mw.CurrentPageTitle())
	return mw.systrayMainWindow
}
