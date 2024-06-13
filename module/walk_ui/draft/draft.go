package draft

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"useful-tools/module/logic/aes"
	"useful-tools/module/walk_ui/base"
	"useful-tools/module/walk_ui/common"
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
	persistence       *base.Persistence
	logicControl      *aes.Aes
	leftDraftContent  *walk.TextEdit
	rightDraftContent *walk.TextEdit
	leftGroupBox      *walk.GroupBox
	rightGroupBox     *walk.GroupBox
	leftDraftPages    map[int]string
	leftCurrentPage   int
	rightDraftPages   map[int]string
	rightCurrentPage  int
}

func NewPage(parent walk.Container, IsConvenientMode bool) (common.Page, error) {
	p := new(Page)
	p.logicControl = logicControl
	p.persistence = persistence
	p.leftDraftPages = make(map[int]string)
	p.rightDraftPages = make(map[int]string)
	p.leftCurrentPage = 1
	p.rightCurrentPage = 1
	if err := (Composite{
		AssignTo: &p.Composite,
		Name:     "draftPage",
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
					// 草稿A
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
								Title:    formatGroupBoxTitle(p.leftCurrentPage),
								AssignTo: &p.leftGroupBox,
								Layout:   VBox{},
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
												AssignTo: &p.leftDraftContent,
												VScroll:  true,
												Text:     p.logicControl.ViewContent(),
												OnTextChanged: func() {
													p.logicControl.SetViewContent(p.leftDraftContent.Text())
												},
												OnMouseDown: func(x, y int, button walk.MouseButton) {
													if button == walk.LeftButton {
														if p.logicControl.DoubleClicked() {
															_ = p.leftDraftContent.SetText(p.logicControl.FormatJson(p.leftDraftContent.Text()))
														}
													}
												},
												OnKeyDown: func(key walk.Key) {
													if walk.ControlDown() {
														if key == walk.KeyDown {
															if p.leftCurrentPage < 9527 {
																p.leftDraftPages[p.leftCurrentPage] = p.leftDraftContent.Text()
																p.leftCurrentPage++
																_ = p.leftGroupBox.SetTitle(formatGroupBoxTitle(p.leftCurrentPage))
																_ = p.leftDraftContent.SetText(p.leftDraftPages[p.leftCurrentPage])
																_ = p.leftDraftContent.SetFocus()
															}
														}
														if key == walk.KeyUp {
															if p.leftCurrentPage > 1 {
																p.leftDraftPages[p.leftCurrentPage] = p.leftDraftContent.Text()
																p.leftCurrentPage--
																_ = p.leftGroupBox.SetTitle(formatGroupBoxTitle(p.leftCurrentPage))
																_ = p.leftDraftContent.SetText(p.leftDraftPages[p.leftCurrentPage])
																_ = p.leftDraftContent.SetFocus()
															}
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
														Text:    "上一页",
														OnClicked: func() {
															if p.leftCurrentPage > 1 {
																p.leftDraftPages[p.leftCurrentPage] = p.leftDraftContent.Text()
																p.leftCurrentPage--
																_ = p.leftGroupBox.SetTitle(formatGroupBoxTitle(p.leftCurrentPage))
																_ = p.leftDraftContent.SetText(p.leftDraftPages[p.leftCurrentPage])
																_ = p.leftDraftContent.SetFocus()
															}
														},
													},
													PushButton{
														Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize: Size{Height: 36},
														MaxSize: Size{Height: 36},
														Text:    "下一页",
														OnClicked: func() {
															if p.leftCurrentPage < 9527 {
																p.leftDraftPages[p.leftCurrentPage] = p.leftDraftContent.Text()
																p.leftCurrentPage++
																_ = p.leftGroupBox.SetTitle(formatGroupBoxTitle(p.leftCurrentPage))
																_ = p.leftDraftContent.SetText(p.leftDraftPages[p.leftCurrentPage])
																_ = p.leftDraftContent.SetFocus()
															}
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
					// 草稿B
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
								Title:    formatGroupBoxTitle(p.rightCurrentPage),
								Layout:   VBox{},
								AssignTo: &p.rightGroupBox,
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
												AssignTo: &p.rightDraftContent,
												VScroll:  true,
												Text:     p.logicControl.ViewContent(),
												OnTextChanged: func() {
													p.logicControl.SetViewContent(p.rightDraftContent.Text())
												},
												OnMouseDown: func(x, y int, button walk.MouseButton) {
													if button == walk.LeftButton {
														if p.logicControl.DoubleClicked() {
															_ = p.rightDraftContent.SetText(p.logicControl.FormatJson(p.rightDraftContent.Text()))
														}
													}
												},
												OnKeyDown: func(key walk.Key) {
													if walk.ControlDown() {
														if key == walk.KeyDown {
															if p.rightCurrentPage < 9527 {
																p.rightDraftPages[p.rightCurrentPage] = p.rightDraftContent.Text()
																p.rightCurrentPage++
																_ = p.rightGroupBox.SetTitle(formatGroupBoxTitle(p.rightCurrentPage))
																_ = p.rightDraftContent.SetText(p.rightDraftPages[p.rightCurrentPage])
																_ = p.rightDraftContent.SetFocus()
															}
														}
														if key == walk.KeyUp {
															if p.rightCurrentPage > 1 {
																p.rightDraftPages[p.rightCurrentPage] = p.rightDraftContent.Text()
																p.rightCurrentPage--
																_ = p.rightGroupBox.SetTitle(formatGroupBoxTitle(p.rightCurrentPage))
																_ = p.rightDraftContent.SetText(p.rightDraftPages[p.rightCurrentPage])
																_ = p.rightDraftContent.SetFocus()
															}
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
														Text:    "上一页",
														OnClicked: func() {
															if p.rightCurrentPage > 1 {
																p.rightDraftPages[p.rightCurrentPage] = p.rightDraftContent.Text()
																p.rightCurrentPage--
																_ = p.rightGroupBox.SetTitle(formatGroupBoxTitle(p.rightCurrentPage))
																_ = p.rightDraftContent.SetText(p.rightDraftPages[p.rightCurrentPage])
																_ = p.rightDraftContent.SetFocus()
															}
														},
													},
													PushButton{
														Font:    Font{Family: "MicrosoftYaHei", PointSize: 14},
														MinSize: Size{Height: 36},
														MaxSize: Size{Height: 36},
														Text:    "下一页",
														OnClicked: func() {
															if p.rightCurrentPage < 9527 {
																p.rightDraftPages[p.rightCurrentPage] = p.rightDraftContent.Text()
																p.rightCurrentPage++
																_ = p.rightGroupBox.SetTitle(formatGroupBoxTitle(p.rightCurrentPage))
																_ = p.rightDraftContent.SetText(p.rightDraftPages[p.rightCurrentPage])
																_ = p.rightDraftContent.SetFocus()
															}
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
