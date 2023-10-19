package tcp_udp

import (
	"github.com/tidwall/gjson"
	"log"
	"useful-tools/helper/Go"
	"useful-tools/module/logic/tcp_udp"
	"useful-tools/module/walk_ui/base"
	"useful-tools/module/walk_ui/common"
	"useful-tools/pkg/wlog"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	logicControl *tcp_udp.TcpUdp
	persistence  *base.Persistence
)

func init() {
	logicControl = tcp_udp.New()
	persistence = base.NewPersistence()
}

type Page struct {
	*walk.Composite
	persistence           *base.Persistence
	logicControl          *tcp_udp.TcpUdp
	network               *walk.ComboBox
	iFaceList             *walk.ComboBox
	targetIp              *walk.LineEdit
	targetPort            *walk.LineEdit
	subButton             *walk.PushButton
	viewContent           *walk.TextEdit
	convenientModeContent *walk.TextEdit
}

func (p *Page) normalDial() {
	network := p.network.Text()
	iFace := p.iFaceList.Text()
	ip := p.targetIp.Text()
	port := p.targetPort.Text()
	encodeParams := p.persistence.SetLatestParams(network, iFace, ip, port)
	Go.Go(func() {
		_, err := p.logicControl.NormalDial(network, iFace, ip, port)
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				wlog.Warm("p.logicControl.NormalDial failed: %+v", err)
				p.AppendContent(logFormat(network, iFace, ip, port, err.Error()))
			} else {
				p.AppendContent(logFormat(network, iFace, ip, port, "OK !"))
			}
		} else {
			wlog.Info("encodeParams(%v) neq p.concurrentParserParams(%v)", encodeParams, p.persistence.GetLatestParams())
		}
	})
}

func (p *Page) convenientDial() {
	content := p.convenientModeContent.Text()
	encodeParams := p.persistence.SetLatestParams(content)
	Go.Go(func() {
		_, err := p.logicControl.ConvenientDial(content)
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				wlog.Warm("p.logicControl.ConvenientDial failed: %+v", err)
				p.AppendContent(logFormat(gjson.Get(content, "network").String(), gjson.Get(content, "local_ip").String(), gjson.Get(content, "host").String(), gjson.Get(content, "port").String(), err.Error()))
			} else {
				p.AppendContent(logFormat(gjson.Get(content, "network").String(), gjson.Get(content, "local_ip").String(), gjson.Get(content, "host").String(), gjson.Get(content, "port").String(), "OK !"))
			}
		} else {
			wlog.Info("encodeParams(%v) neq p.concurrentParserParams(%v)", encodeParams, p.persistence.GetLatestParams())
		}
	})
}

func (p *Page) PrintContent(content string) {
	p.viewContent.SetText(content)
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
		Name:     "tcpUdpPage",
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
											Label{
												Text:      "协议类型:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											ComboBox{
												Font:          Font{PointSize: 16},
												AssignTo:      &p.network,
												Value:         getNetwork(p.logicControl.Network()),
												Model:         network(),
												DisplayMember: "Name",
												BindingMember: "Key",
												OnCurrentIndexChanged: func() {
													p.logicControl.SetNetwork(p.network.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.iFaceList.SetFocus()
													}
												},
											},
											Label{
												Text:      "本地网卡:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											ComboBox{
												Font:          Font{PointSize: 16},
												AssignTo:      &p.iFaceList,
												Model:         defaultIFaceList(),
												DisplayMember: "Name",
												BindingMember: "Key",
												OnCurrentIndexChanged: func() {
													p.logicControl.SetIFace(p.iFaceList.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.targetIp.SetFocus()
													}
												},
											},
											//VSpacer{Size: 20},
											Label{
												Text:      "目标地址:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.targetIp,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText: "请输入目标地址",
												MinSize:     Size{Height: 36},
												MaxSize:     Size{Height: 36},
												Text:        p.logicControl.Host(),
												OnTextChanged: func() {
													p.logicControl.SetHost(p.targetIp.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.targetPort.SetFocus()
													}
												},
											},
											//VSpacer{Size: 20},
											Label{
												Text:      "目标端口:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.targetPort,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText: "请输入目标端口",
												MinSize:     Size{Height: 36},
												MaxSize:     Size{Height: 36},
												Text:        p.logicControl.Port(),
												OnTextChanged: func() {
													p.logicControl.SetPort(p.targetPort.Text())
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
														Text:     "连接",
														OnClicked: func() {
															p.normalDial()
														},
														//OnKeyPress: func(key walk.Key) {
														//	fmt.Println(key)
														//	if key == walk.KeyReturn {
														//		p.normalDial()
														//	}
														//},
													},
													PushButton{
														Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize: Size{Height: 36},
														MaxSize: Size{Height: 36},
														Text:    "清空",
														OnClicked: func() {
															_ = p.targetIp.SetText("")
															_ = p.targetPort.SetText("")
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
								Title:  "Parameters",
								Layout: VBox{},
								Children: []Widget{
									TextEdit{
										Font:     Font{Family: "MicrosoftYaHei", PointSize: 15},
										AssignTo: &p.convenientModeContent,
										Text:     getDialInfo(p.logicControl.ProTemplate(), p.logicControl.RequestInfo()),
										VScroll:  true,
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
													p.convenientDial()
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
												ReadOnly: true,
												VScroll:  true,
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
	iFaceList := p.logicControl.GetIFaceList()
	err := p.iFaceList.SetModel(createIFaceList(iFaceList))
	if err != nil {
		wlog.Warm("p.iFaceList.SetModel failed: %v", err)
	}
	err = p.iFaceList.SetCurrentIndex(getInterface(p.logicControl.IFace(), iFaceList))
	if err != nil {
		wlog.Warm("p.iFaceList.SetCurrentIndex failed: %v", err)
	}
	return p, nil
}

func normalModeState(mode bool) Property {
	return !mode
}

func convenientModeState(mode bool) Property {
	return mode
}
