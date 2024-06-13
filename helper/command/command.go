package command

import (
	"io"
)

type Cmd interface {
	Start() error
	Wait() error
	Run() error
	SetOutput(w io.Writer)
	Stop() error
}
