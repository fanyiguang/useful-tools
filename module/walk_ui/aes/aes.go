package aes

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"useful-tools/helper/Go"
	"useful-tools/module/logic/aes"
	"useful-tools/module/walk_ui/base"
	"useful-tools/module/walk_ui/common"
	"useful-tools/pkg/wlog"
)

var (
	logicControl *aes.Aes
	persistence  *base.Persistence
)

func init() {
	logicControl = aes.New()
	persistence = base.NewPersistence()
}

type Page struct {
	*walk.Composite
	persistence           *base.Persistence
	logicControl          *aes.Aes
	convertType           *walk.ComboBox
	key                   *walk.LineEdit
	iv                    *walk.LineEdit
	inputContent          *walk.TextEdit
	subButton             *walk.PushButton
	viewContent           *walk.TextEdit
	convenientModeContent *walk.TextEdit
}

func (p *Page) aesAnalysis() {
	convertTyp := p.convertType.Text()
	key := p.key.Text()
	iv := p.iv.Text()
	inputContent := p.inputContent.Text()
	if base.MenuItemLogic.SaveAesKey() {
		Go.Go(func() {
			base.MenuItemLogic.SetAesKeyToFile(key)
			base.MenuItemLogic.SetAesIVToFile(iv)
		})
	}
	encodeParams := p.persistence.SetLatestParams(convertTyp, key, iv, inputContent)
	Go.Go(func() {
		var (
			data string
			err  error
		)
		switch convertTyp {
		case "加密":
			data, err = p.logicControl.Encode(key, iv, inputContent)
		case "解密":
			data, err = p.logicControl.Decode(key, iv, inputContent)
		}
		if p.persistence.Equal(encodeParams) {
			if err != nil {
				p.PrintContent(err.Error())
			} else {
				p.PrintContent(data)
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
								//Background: SolidColorBrush{ // 增加背景颜色
								//	Color: walk.RGB(14, 149, 9),
								//},
								Children: []Widget{
									Composite{
										Layout: Grid{
											Rows: 10,
										},
										//Background: SolidColorBrush{ // 增加背景颜色
										//	Color: walk.RGB(114, 49, 9),
										//},
										Children: []Widget{
											//VSpacer{MinSize: Size{Height: 5}},
											Label{
												Text:      "转换类型:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											ComboBox{
												Font:          Font{PointSize: 16},
												AssignTo:      &p.convertType,
												Value:         getConvertType(p.logicControl.ConvertType()),
												Model:         convertType(),
												DisplayMember: "Name",
												BindingMember: "Key",
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.key.SetFocus()
													}
												},
												OnCurrentIndexChanged: func() {
													p.logicControl.SetConvertType(p.convertType.Text())
													_ = p.subButton.SetText(p.convertType.Text())
												},
											},
											Label{
												Text:      "KEY:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.key,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText: "以逗号分割AES KEY",
												MinSize:     Size{Height: 36},
												MaxSize:     Size{Height: 36},
												Text:        getText(base.MenuItemLogic.AesKey(), p.logicControl.Key()),
												OnTextChanged: func() {
													p.logicControl.SetKey(p.key.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.iv.SetFocus()
														return
													}
												},
											},
											Label{
												Text:      "IV:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											LineEdit{
												AssignTo:    &p.iv,
												TextColor:   walk.RGB(40, 40, 42),
												Background:  TransparentBrush{},
												Font:        Font{Family: "MicrosoftYaHei", PointSize: 14},
												ToolTipText: "以逗号分割AES IV",
												MinSize:     Size{Height: 36},
												MaxSize:     Size{Height: 36},
												Text:        getText(base.MenuItemLogic.AesIV(), p.logicControl.Iv()),
												OnTextChanged: func() {
													p.logicControl.SetIv(p.iv.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.inputContent.SetFocus()
														return
													}
												},
											},
											//VSpacer{Size: 20},
											Label{
												Text:      "参数:",
												TextColor: walk.RGB(91, 92, 96),
												Font:      Font{PointSize: 12, Family: "MicrosoftYaHei"},
											},
											TextEdit{
												Font:     Font{Family: "MicrosoftYaHei", PointSize: 14},
												AssignTo: &p.inputContent,
												VScroll:  true,
												Text:     p.logicControl.InputContent(),
												//MinSize:  Size{Height: 355},
												MaxSize: Size{Height: 450},
												OnTextChanged: func() {
													p.logicControl.SetInputContent(p.inputContent.Text())
												},
												OnKeyPress: func(key walk.Key) {
													if key == walk.KeyTab || key == walk.KeyReturn {
														_ = p.subButton.SetFocus()
													}
												},
											},
											//VSpacer{Size: 20},
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
														Text:     "解密",
														OnClicked: func() {
															p.aesAnalysis()
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
															//_ = p.key.SetText("")
															//_ = p.iv.SetText("")
															_ = p.inputContent.SetText("")
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
										//Background: SolidColorBrush{ // 增加背景颜色
										//	Color: walk.RGB(114, 49, 9),
										//},
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
												OnMouseDown: func(x, y int, button walk.MouseButton) {
													if button == walk.LeftButton {
														if p.logicControl.DoubleClicked() {
															_ = p.viewContent.SetText(p.logicControl.FormatJson(p.viewContent.Text()))
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
	return true
}

func convenientModeState(mode bool) Property {
	return mode
}
