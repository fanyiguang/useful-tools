package controller

import (
	"github.com/sagernet/sing/common/x/linkedhashmap"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/model"
)

var _ adapter.Controller = (*Draft)(nil)

type Draft struct {
	Base
	leftDocs  linkedhashmap.Map[string, model.Doc]
	rightDocs linkedhashmap.Map[string, model.Doc]
}

func NewDraft() *Draft {
	return &Draft{}
}

func (d *Draft) LeftDocs() []model.Doc {
	return d.leftDocs.Values()
}

func (d *Draft) RightDocs() []model.Doc {
	return d.rightDocs.Values()
}

func (d *Draft) AddLeftDocs(title, content, placeHolder string) {
	d.leftDocs.Put(title, model.Doc{
		Title:       title,
		Content:     content,
		PlaceHolder: placeHolder,
	})
}

func (d *Draft) AddRightDocs(title, content, placeHolder string) {
	d.rightDocs.Put(title, model.Doc{
		Title:       title,
		Content:     content,
		PlaceHolder: placeHolder,
	})
}

func (d *Draft) FindRightNextDocsIndex(title string) (model.Doc, int) {
	for i, doc := range d.rightDocs.Values() {
		if title == doc.Title {
			a := i + 1
			if a >= d.rightDocs.Size() {
				return model.Doc{}, -1
			}
			return doc, a
		}
	}
	return model.Doc{}, -1
}

func (d *Draft) FindRightPrevDocsIndex(title string) (model.Doc, int) {
	for i, doc := range d.rightDocs.Values() {
		if title == doc.Title {
			a := i - 1
			return doc, a
		}
	}
	return model.Doc{}, -1
}

func (d *Draft) FindLeftNextDocsIndex(title string) (model.Doc, int) {
	for i, doc := range d.leftDocs.Values() {
		if title == doc.Title {
			a := i + 1
			if a >= d.leftDocs.Size() {
				return model.Doc{}, -1
			}
			return doc, a
		}
	}
	return model.Doc{}, -1
}

func (d *Draft) FindLeftPrevDocsIndex(title string) (model.Doc, int) {
	for i, doc := range d.leftDocs.Values() {
		if title == doc.Title {
			a := i - 1
			return doc, a
		}
	}
	return model.Doc{}, -1
}

func (d *Draft) ClearCache() {
	for _, val := range d.leftDocs.Values() {
		val.Content = ""
	}
	for _, val := range d.rightDocs.Values() {
		val.Content = ""
	}
}
