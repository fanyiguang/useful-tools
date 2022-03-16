package dns

import (
	"log"
	"strings"
	"useful-tools/helper/Go"
	"useful-tools/module/logic/dns"
	"useful-tools/module/walk_ui/common"
	"useful-tools/pkg/wlog"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Page struct {
	*walk.Composite
	logicControl          *dns.Dns
	dnsServerAddr         *walk.LineEdit
	parserDomain          *walk.LineEdit
	subButton             *walk.PushButton
	viewContent           *walk.TextEdit
	convenientModeContent *walk.TextEdit
}

func (p *Page) normalDns() {
	Go.Go(func() {
		ips, err := p.logicControl.NormalDns(p.dnsServerAddr.Text(), p.parserDomain.Text())
		if err != nil {
			wlog.Warm("p.logicControl.NormalDns failed: %+v", err)
			p.PrintContent(err.Error())
		} else {
			p.PrintContent(strings.Join(ips, "\r\n"))
		}
	})
}

func (p *Page) convenientDns() {
	Go.Go(func() {
		ips, err := p.logicControl.ConvenientDns(p.convenientModeContent.Text())
		if err != nil {
			wlog.Warm("p.logicControl.ConvenientDns failed: %+v", err)
			p.PrintContent(err.Error())
		} else {
			p.PrintContent(strings.Join(ips, "\r\n"))
		}
	})
}

func (p *Page) PrintContent(content string) {
	_ = p.viewContent.SetText(content)
}

func NewPage(parent walk.Container, IsConvenientMode bool) (common.Page, error) {
	p := new(Page)
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
								Title: "Gradient Parameters",
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
										Text:        "默认",
										ToolTipText: "请输入DNS地址",
										MinSize:     Size{Height: 36},
										MaxSize:     Size{Height: 36},
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
								Title:  "Gradient Parameters",
								Layout: VBox{},
								Children: []Widget{
									TextEdit{
										Font:     Font{Family: "MicrosoftYaHei", PointSize: 15},
										AssignTo: &p.convenientModeContent,
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
								Title:  "Gradient Parameters",
								Layout: VBox{},
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
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	if err := walk.InitWrapperWindow(p); err != nil {
		return nil, err
	}
	p.logicControl = dns.New()
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
