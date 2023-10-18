package base

import (
	"useful-tools/module/logic/common"

	"github.com/pkg/errors"
)

type Base struct {
	executing   bool
	proTemplate string
}

func (b *Base) ProTemplate() string {
	return b.proTemplate
}

func (b *Base) SetProTemplate(proTemplate string) {
	b.proTemplate = proTemplate
}

func (b *Base) IsExecuting() bool {
	return b.executing
}

func (b *Base) SetExecuting() {
	b.executing = true
}

func (b *Base) ResetExecuting() {
	b.executing = false
}

func (b *Base) IsExecutingError(err error) bool {
	return errors.Is(err, common.ExecutingError)
}
