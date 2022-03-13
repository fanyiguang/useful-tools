package proxy

import (
	"log"
	"useful-tools/helper/Go"
	"useful-tools/module/logic/proxy"
	"useful-tools/module/walk_ui/common"
	"useful-tools/pkg/wlog"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Page struct {
	*walk.Composite
	logicControl  *proxy.Proxy
	proxyType     *walk.ComboBox
	proxyIp       *walk.LineEdit
	proxyPort     *walk.LineEdit
	proxyUsername *walk.LineEdit
	proxyPassword *walk.LineEdit
	subButton     *walk.PushButton
	viewContent   *walk.TextEdit
}

type CompanyItem struct {
	Name          string
	CompanyId     int
	UserTitleType int
	Key           int
	UserId        int
	//BindMsg string
}

func (p *Page) checkProxy() {
	Go.Go(func() {
		checkProxy, err := p.logicControl.CheckProxy(p.proxyIp.Text(), p.proxyPort.Text(), p.proxyUsername.Text(), p.proxyPassword.Text(), p.proxyType.Text())
		if err != nil {
			wlog.Warm("p.logicControl.CheckProxy failed: %+v", err)
		} else {
			p.PrintContent(checkProxy)
		}
	})
}

func (p *Page) PrintContent(content string) {
	p.viewContent.SetText(content)
}

func NewPage(parent walk.Container) (common.Page, error) {
	p := new(Page)
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "fooPage",
		Layout: Grid{
			MarginsZero: true,
			Rows:        1,
		},
		Background: SolidColorBrush{ // 增加背景颜色
			Color: walk.RGB(124, 149, 9),
		},
		Children: []Widget{
			Composite{
				Background: SolidColorBrush{ // 增加背景颜色
					Color: walk.RGB(124, 89, 99),
				},
				Layout: Grid{
					Rows:        2,
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					Composite{
						StretchFactor: 1,
						Background: SolidColorBrush{ // 增加背景颜色
							Color: walk.RGB(14, 249, 9),
						},
						//MinSize: Size{Width: 600},
						Layout: Grid{
							Alignment: AlignHVDefault,
							Rows:      12,
							//Columns: 1,
							//MarginsZero: true,
						},
						Children: []Widget{
							//VSpacer{MinSize: Size{Height: 5}},
							Label{
								Text:      "代理类型:",
								TextColor: walk.RGB(91, 92, 96),
								Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
							},
							ComboBox{
								Font:          Font{PointSize: 16},
								AssignTo:      &p.proxyType,
								Value:         1,
								Model:         getProxyType(),
								DisplayMember: "Name",
								BindingMember: "Key",
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyTab || key == walk.KeyReturn {
										_ = p.proxyIp.SetFocus()
									}
								},
							},
							Label{
								Text:      "代理地址:",
								TextColor: walk.RGB(91, 92, 96),
								Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
							},
							LineEdit{
								AssignTo:    &p.proxyIp,
								TextColor:   walk.RGB(40, 40, 42),
								Background:  TransparentBrush{},
								Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
								ToolTipText: "请输入代理地址",
								MinSize:     Size{Height: 36},
								MaxSize:     Size{Height: 36},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyTab || key == walk.KeyReturn {
										_ = p.proxyPort.SetFocus()
									}
								},
							},
							//VSpacer{Size: 20},
							Label{
								Text:      "代理端口:",
								TextColor: walk.RGB(91, 92, 96),
								Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
							},
							LineEdit{
								AssignTo:    &p.proxyPort,
								TextColor:   walk.RGB(40, 40, 42),
								Background:  TransparentBrush{},
								Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
								ToolTipText: "请输入代理端口",
								MinSize:     Size{Height: 36},
								MaxSize:     Size{Height: 36},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyTab || key == walk.KeyReturn {
										_ = p.proxyUsername.SetFocus()
									}
								},
							},
							//VSpacer{Size: 20},
							Label{
								Text:      "代理账号:",
								TextColor: walk.RGB(91, 92, 96),
								Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
							},
							LineEdit{
								AssignTo:    &p.proxyUsername,
								TextColor:   walk.RGB(40, 40, 42),
								Background:  TransparentBrush{},
								Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
								ToolTipText: "请输入代理账号",
								MinSize:     Size{Height: 36},
								MaxSize:     Size{Height: 36},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyTab || key == walk.KeyReturn {
										_ = p.proxyPassword.SetFocus()
									}
								},
							},
							Label{
								Text:      "代理密码:",
								TextColor: walk.RGB(91, 92, 96),
								Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
							},
							LineEdit{
								AssignTo:  &p.proxyPassword,
								TextColor: walk.RGB(40, 40, 42),
								//PasswordMode: true,
								Background:  TransparentBrush{},
								Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
								ToolTipText: "请输入代理密码",
								MinSize:     Size{Height: 36},
								MaxSize:     Size{Height: 36},
								OnKeyPress: func(key walk.Key) {
									if key == walk.KeyTab || key == walk.KeyReturn {
										_ = p.subButton.SetFocus()
									}
								},
							},
							PushButton{
								AssignTo: &p.subButton,
								Font:     Font{Family: "MicrosoftYaHei", PointSize: 14},
								MinSize:  Size{Height: 36},
								MaxSize:  Size{Height: 36},
								Text:     "检测",
								OnClicked: func() {
									p.checkProxy()
								},
								//OnKeyPress: func(key walk.Key) {
								//	fmt.Println(key)
								//	if key == walk.KeyReturn {
								//		p.checkProxy()
								//	}
								//},
							},
							VSpacer{},
							//VSpacer{Size: 10},
						},
					},
					Composite{
						//StretchFactor: 1,
						//Background: SolidColorBrush{ // 增加背景颜色
						//	Color: walk.RGB(54, 29, 9),
						//},
						//MinSize: Size{Width: 600},
						Layout: Grid{
							Rows:        2,
							MarginsZero: true,
							//SpacingZero: true,
						},
						Children: []Widget{
							TextEdit{
								Font:     Font{Family: "MicrosoftYaHei", PointSize: 15},
								AssignTo: &p.viewContent,
								ReadOnly: true,
							},
							Composite{
								StretchFactor: 1,
								//Background: SolidColorBrush{ // 增加背景颜色
								//	Color: walk.RGB(54, 29, 9),
								//},
								//MinSize: Size{Width: 600},
								Layout: Grid{
									Columns:     2,
									MarginsZero: true,
									SpacingZero: true,
								},
								Children: []Widget{
									PushButton{
										Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
										MinSize: Size{Height: 36},
										MaxSize: Size{Height: 36},
										Text:    "复制",
										OnClicked: func() {
											if err := walk.Clipboard().SetText(p.viewContent.Text()); err != nil {
												log.Print("Copy: ", err)
											}
										},
									},
									PushButton{
										Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
										MinSize: Size{Height: 36},
										MaxSize: Size{Height: 36},
										Text:    "清空",
										OnClicked: func() {
											_ = p.viewContent.SetText("")
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}
	p.logicControl = proxy.New()
	return p, nil
}

func getProxyType() []*CompanyItem {
	return []*CompanyItem{
		{Key: 1, Name: "SOCKS5"},
		{Key: 2, Name: "SSL"},
		{Key: 3, Name: "SSH"},
		{Key: 4, Name: "HTTP"},
		{Key: 5, Name: "HTTPS"},
	}
}
