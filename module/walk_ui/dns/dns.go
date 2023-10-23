package dns

import (
	"github.com/tidwall/gjson"
	"log"
	"strings"
	"useful-tools/helper/Go"
	"useful-tools/module/logic/dns"
	"useful-tools/module/walk_ui/base"
	"useful-tools/module/walk_ui/common"
	"useful-tools/pkg/wlog"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	logicControl *dns.Dns
	persistence  *base.Persistence
)

func init() {
	logicControl = dns.New()
	persistence = base.NewPersistence()
}

type Page struct {
	*walk.Composite
	persistence           *base.Persistence
	logicControl          *dns.Dns
	dnsServerAddr         *walk.LineEdit
	parserDomain          *walk.LineEdit
	subButton             *walk.PushButton
	viewContent           *walk.TextEdit
	convenientModeContent *walk.TextEdit

	LogView *walk.ScrollView
}

func (p *Page) normalDns() {
	server := p.dnsServerAddr.Text()
	domain := p.parserDomain.Text()
	encodeParams := p.persistence.SetLatestParams(server, domain)
	Go.Go(func() {
		ips, err := p.logicControl.NormalDns(server, domain)
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				wlog.Warm("p.logicControl.NormalDns failed: %+v", err)
				p.AppendContent(logFormat(server, domain, err.Error()))
			} else {
				p.AppendContent(logFormat(server, domain, strings.Join(ips, " ")))
			}
		} else {
			wlog.Info("encodeParams(%v) neq p.concurrentParserParams(%v)", encodeParams, p.persistence.GetLatestParams())
		}
	})
}

func (p *Page) convenientDns() {
	content := p.convenientModeContent.Text()
	encodeParams := p.persistence.SetLatestParams(content)
	Go.Go(func() {
		ips, err := p.logicControl.ConvenientDns(content)
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				wlog.Warm("p.logicControl.ConvenientDns failed: %+v", err)
				p.AppendContent(logFormat(gjson.Get(content, "server").String(), gjson.Get(content, "domain").String(), err.Error()))
			} else {
				p.AppendContent(logFormat(gjson.Get(content, "server").String(), gjson.Get(content, "domain").String(), strings.Join(ips, " ")))
			}
		} else {
			wlog.Info("encodeParams(%v) neq p.concurrentParserParams(%v)", encodeParams, p.persistence.GetLatestParams())
		}
	})
}

func (p *Page) PrintContent(content string) {
	_ = p.viewContent.SetText(content)
}

func (p *Page) AppendContent(content string) {
	p.viewContent.AppendText(content)
}

func NewPage(parent walk.Container, IsConvenientMode bool) (common.Page, error) {
	p := new(Page)
	p.logicControl = logicControl
	p.persistence = persistence
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "dnsPage",
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
											Rows: 12,
										},
										Children: []Widget{
											//VSpacer{MinSize: Size{Height: 5}},
											//VSpacer{Size: 20},
											Label{
												Text:      "DNS地址:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.dnsServerAddr,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												Text:        getServer(p.logicControl.Server()),
												ToolTipText: "请输入DNS地址",
												MinSize:     Size{Height: 36},
												MaxSize:     Size{Height: 36},
												OnTextChanged: func() {
													p.logicControl.SetServer(p.dnsServerAddr.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.parserDomain.SetFocus()
													}
												},
											},
											//VSpacer{Size: 20},
											Label{
												Text:      "解析域名:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.parserDomain,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText: "请输入解析域名",
												MinSize:     Size{Height: 36},
												MaxSize:     Size{Height: 36},
												Text:        p.logicControl.Domain(),
												OnTextChanged: func() {
													p.logicControl.SetDomain(p.parserDomain.Text())
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
														Text:     "解析",
														OnClicked: func() {
															p.normalDns()
														},
														//OnKeyPress: func(key walk.Key) {
														//	fmt.Println(key)
														//	if key == walk.KeyReturn {
														//		p.normalDns()
														//	}
														//},
													},
													PushButton{
														Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize: Size{Height: 36},
														MaxSize: Size{Height: 36},
														Text:    "清空",
														OnClicked: func() {
															_ = p.dnsServerAddr.SetText("")
															_ = p.parserDomain.SetText("")
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
					// 解析模式
					Composite{
						Name:          "convenient_mode",
						Visible:       convenientModeState(IsConvenientMode),
						StretchFactor: 1,
						//Background: SolidColorBrush{ // 增加背景颜色
						//	Color: walk.RGB(14, 249, 9),
						//},
						//MinSize: Size{Width: 600},
						Layout: Grid{
							Alignment: AlignHVDefault,
							Rows:      1,
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
											Rows: 2,
										},
										Children: []Widget{
											TextEdit{
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 15},
												AssignTo:    &p.convenientModeContent,
												Text:        getDnsInfo(p.logicControl.ProTemplate(), p.logicControl.RequestInfo()),
												VScroll:     true,
												ToolTipText: "双击格式化内容！",
												OnTextChanged: func() {
													p.logicControl.SetRequestInfo(p.convenientModeContent.Text())
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
														Text:     "解析",
														OnClicked: func() {
															p.convenientDns()
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

func getNetwork() []*common.CompanyItem {
	return []*common.CompanyItem{
		{Key: 1, Name: "TCP"},
		{Key: 2, Name: "UDP"},
	}
}

func getDefaultIFaceList() []*common.CompanyItem {
	return []*common.CompanyItem{
		{Key: 1, Name: "随机"},
	}
}

func createIFaceList(ips []string) []*common.CompanyItem {
	list := getDefaultIFaceList()
	for i, ip := range ips {
		list = append(list, &common.CompanyItem{
			Name: ip,
			Key:  i + 1,
		})
	}
	return list
}
