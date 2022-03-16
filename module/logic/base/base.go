package base

type Base struct {
	executing bool
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
