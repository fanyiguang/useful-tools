package controller

import "useful-tools/app/usefultools/model"

type Draft struct {
	leftDocs  map[string]model.Doc
	rightDocs map[string]model.Doc
}

func NewDraft() *Draft {
	return &Draft{
		leftDocs:  make(map[string]model.Doc),
		rightDocs: make(map[string]model.Doc),
	}
}

func (d *Draft) LeftDocs() map[string]model.Doc {
	return d.leftDocs
}

func (d *Draft) RightDocs() map[string]model.Doc {
	return d.rightDocs
}

func (d *Draft) AddLeftDocs(title, content, placeHolder string) {
	d.leftDocs[title] = model.Doc{
		Title:       title,
		Content:     content,
		PlaceHolder: placeHolder,
	}
}

func (d *Draft) AddRightDocs(title, content, placeHolder string) {
	d.rightDocs[title] = model.Doc{
		Title:       title,
		Content:     content,
		PlaceHolder: placeHolder,
	}
}
