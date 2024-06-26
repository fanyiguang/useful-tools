package proxy

import (
	"log"
	"strings"
	"useful-tools/helper/Go"
	"useful-tools/module/logic/proxy"
	"useful-tools/module/walk_ui/base"
	"useful-tools/module/walk_ui/common"
	"useful-tools/pkg/wlog"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	logicControl *proxy.Proxy
	persistence  *base.Persistence
)

func init() {
	logicControl = proxy.New()
	persistence = base.NewPersistence()
}

type Page struct {
	*walk.Composite
	persistence           *base.Persistence
	logicControl          *proxy.Proxy
	proxyType             *walk.ComboBox
	proxyIp               *walk.LineEdit
	proxyPort             *walk.LineEdit
	proxyUsername         *walk.LineEdit
	proxyPassword         *walk.LineEdit
	proxyReqURLs          *walk.LineEdit
	subButton             *walk.PushButton
	viewContent           *walk.TextEdit
	convenientModeContent *walk.TextEdit
	proxyUsernameLabel    *walk.Label
}

func (p *Page) normalCheckProxy() {
	ip := p.proxyIp.Text()
	typ := p.proxyType.Text()
	port := p.proxyPort.Text()
	urls := p.proxyReqURLs.Text()
	username := p.proxyUsername.Text()
	password := p.proxyPassword.Text()
	encodeParams := p.persistence.SetLatestParams(ip, port, username, password, typ, urls)
	Go.Go(func() {
		checkProxy, err := p.logicControl.NormalCheckProxy(ip, port, username, password, typ, urls, base.MenuItemLogic.HiddenBody())
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				wlog.Warm("p.logicControl.NormalCheckProxy failed: %+v", err)
				p.PrintContent(err.Error())
			} else {
				p.PrintContent(checkProxy)
			}
		} else {
			wlog.Info("encodeParams(%v) neq p.concurrentParserParams(%v)", encodeParams, p.persistence.GetLatestParams())
		}
	})
}

func (p *Page) convenientCheckProxy() {
	content := p.convenientModeContent.Text()
	encodeParams := p.persistence.SetLatestParams(content)
	Go.Go(func() {
		checkProxy, err := p.logicControl.ConvenientCheckProxy(content, base.MenuItemLogic.HiddenBody())
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				wlog.Warm("p.logicControl.ConvenientCheckProxy failed: %+v", err)
				p.PrintContent(err.Error())
			} else {
				p.PrintContent(strings.Replace(checkProxy, "\n", "\r\n", -1))
			}
		} else {
			wlog.Info("encodeParams(%v) neq p.concurrentParserParams(%v)", encodeParams, p.persistence.GetLatestParams())
		}
	})
}

func (p *Page) PrintContent(content string) {
	_ = p.viewContent.SetText(content)
}

func NewPage(parent walk.Container, IsConvenientMode bool) (common.Page, error) {
	p := new(Page)
	p.logicControl = logicControl
	p.persistence = persistence
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "proxyPage",
		Layout: Grid{
			MarginsZero: true,
			Rows:        1,
		},
		//Background: SolidColorBrush{ // 增加背景颜色
		//	Color: walk.RGB(124, 149, 9),
		//},
		Children: []Widget{
			Composite{
				//Background: SolidColorBrush{ // 增加背景颜色
				//	Color: walk.RGB(124, 89, 99),
				//},
				Layout: Grid{
					Rows:        2,
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					// 标准模式
					Composite{
						StretchFactor: 1,
						//Background: SolidColorBrush{ // 增加背景颜色
						//	Color: walk.RGB(14, 249, 9),
						//},
						Visible: normalModeState(IsConvenientMode),
						//MinSize: Size{Width: 600},
						Layout: Grid{
							Rows: 1,
							//Columns: 1,
							//MarginsZero: true,
						},
						Children: []Widget{
							GroupBox{
								Title: "Parameters",
								Layout: Grid{
									Rows: 1,
								},
								Children: []Widget{
									Composite{
										Layout: Grid{
											Rows: 14,
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
												Value:         getProxyType(p.logicControl.Typ()),
												Model:         proxyType(),
												DisplayMember: "Name",
												BindingMember: "Key",
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.proxyIp.SetFocus()
													}
												},
												OnCurrentIndexChanged: func() {
													p.logicControl.SetTyp(p.proxyType.Text())
													if p.proxyType.Text() == "SHADOWSOCKS" {
														_ = p.proxyUsernameLabel.SetText("加密类型:")
													} else {
														_ = p.proxyUsernameLabel.SetText("代理账号:")
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
												Text:        p.logicControl.Host(),
												OnTextChanged: func() {
													p.logicControl.SetHost(p.proxyIp.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.proxyPort.SetFocus()
														return
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
												Text:        p.logicControl.Port(),
												OnTextChanged: func() {
													p.logicControl.SetPort(p.proxyPort.Text())
												},
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
												AssignTo:  &p.proxyUsernameLabel,
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
												Text:        p.logicControl.Username(),
												OnTextChanged: func() {
													p.logicControl.SetUsername(p.proxyUsername.Text())
												},
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
												AssignTo:     &p.proxyPassword,
												TextColor:    walk.RGB(40, 40, 42),
												PasswordMode: true,
												Background:   TransparentBrush{},
												Font:         Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText:  "请输入代理密码",
												MinSize:      Size{Height: 36},
												MaxSize:      Size{Height: 36},
												Text:         p.logicControl.Password(),
												OnTextChanged: func() {
													p.logicControl.SetPassword(p.proxyPassword.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.proxyReqURLs.SetFocus()
													}
												},
											},
											Label{
												Text:      "请求地址(选填):",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.proxyReqURLs,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText: "代理请求的地址，多个地址请用;隔开！",
												//CueBanner:   "非必填",
												MinSize: Size{Height: 36},
												MaxSize: Size{Height: 36},
												Text:    p.logicControl.Urls(),
												OnTextChanged: func() {
													p.logicControl.SetUrls(p.proxyReqURLs.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.subButton.SetFocus()
													}
												},
											},
											VSpacer{},
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
														AssignTo: &p.subButton,
														Font:     Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize:  Size{Height: 36},
														MaxSize:  Size{Height: 36},
														Text:     "检测",
														OnClicked: func() {
															p.normalCheckProxy()
														},
														//OnKeyPress: func(key walk.Key) {
														//	fmt.Println(key)
														//	if key == walk.KeyReturn {
														//		p.normalCheckProxy()
														//	}
														//},
													},
													PushButton{
														Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize: Size{Height: 36},
														MaxSize: Size{Height: 36},
														Text:    "清空",
														OnClicked: func() {
															_ = p.proxyIp.SetText("")
															_ = p.proxyPort.SetText("")
															_ = p.proxyUsername.SetText("")
															_ = p.proxyPassword.SetText("")
														},
													},
												},
											},
											//VSpacer{Size: 10},
										},
									},
								},
							},
						},
					},
					// 专业模式
					Composite{
						Name:          "convenient_mode",
						Visible:       convenientModeState(IsConvenientMode),
						StretchFactor: 1,
						//Background: SolidColorBrush{ // 增加背景颜色
						//	Color: walk.RGB(14, 249, 9),
						//},
						//MinSize: Size{Width: 600},
						Layout: Grid{
							Rows: 1,
							//Columns: 1,
							//MarginsZero: true,
						},
						Children: []Widget{
							GroupBox{
								Title:  "Parameters",
								Layout: VBox{},
								Children: []Widget{
									Composite{
										Layout: Grid{
											Alignment: AlignHVDefault,
											Rows:      2,
										},
										Children: []Widget{
											TextEdit{
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 15},
												AssignTo:    &p.convenientModeContent,
												Text:        getRequestInfo(p.logicControl.ProTemplate(), p.logicControl.RequestInfo()),
												VScroll:     true,
												ToolTipText: "双击格式化内容！",
												OnTextChanged: func() {
													p.logicControl.SetRequestInfo(p.convenientModeContent.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyReturn {
														_ = p.subButton.Checked()
													}
												},
												OnMouseDown: func(x, y int, button walk.MouseButton) {
													if button == walk.LeftButton {
														if p.logicControl.DoubleClicked() {
															_ = p.convenientModeContent.SetText(p.logicControl.FormatJson(p.convenientModeContent.Text()))
														}
													}
												},
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
														AssignTo: &p.subButton,
														Font:     Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize:  Size{Height: 36},
														MaxSize:  Size{Height: 36},
														Text:     "检测",
														OnClicked: func() {
															p.convenientCheckProxy()
														},
													},
													PushButton{
														Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize: Size{Height: 36},
														MaxSize: Size{Height: 36},
														Text:    "清空",
														OnClicked: func() {
															_ = p.convenientModeContent.SetText("")
														},
													},
												},
											},
										},
									},
								},
							},
							//VSpacer{},
						},
					},
					// 输出页面
					Composite{
						//StretchFactor: 1,
						//Background: SolidColorBrush{ // 增加背景颜色
						//	Color: walk.RGB(54, 29, 9),
						//},
						//MinSize: Size{Width: 600},
						Layout: Grid{
							Rows: 1,
						},
						Children: []Widget{
							GroupBox{
								Title:  "View",
								Layout: VBox{},
								Children: []Widget{
									Composite{
										Layout: Grid{
											Rows: 2,
										},
										Children: []Widget{
											TextEdit{
												Font:     Font{Family: "MicrosoftYaHei", PointSize: 15},
												AssignTo: &p.viewContent,
												VScroll:  true,
												ReadOnly: true,
												Text:     p.logicControl.ViewContent(),
												OnTextChanged: func() {
													p.logicControl.SetViewContent(p.viewContent.Text())
												},
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
					},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}
	p.proxyPassword.SetPasswordMode(!base.MenuItemLogic.ShowPass())
	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}
	return p, nil
}

func normalModeState(mode bool) Property {
	return !mode
}

func convenientModeState(mode bool) Property {
	return mode
}
