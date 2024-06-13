//go:build windows
// +build windows

package command

import (
	"bytes"
	"io"
	"os/exec"
)

type winCmd struct {
	cmd *exec.Cmd
	g   ProcessExitGroup
}

func New(name string, args ...string) (Cmd, error) {
	g, err := NewProcessExitGroup()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = bytes.NewBuffer([]byte{})
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	return &winCmd{
		cmd: cmd,
		g:   g,
	}, nil
}

func (c *winCmd) Start() error {
	err := c.cmd.Start()
	if err != nil {
		return err
	}
	err = c.g.AddProcess(c.cmd.Process)
	if err != nil {
		return err
	}
	return nil
}

func (c *winCmd) Wait() error {
	return c.cmd.Wait()
}

func (c *winCmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *winCmd) SetOutput(w io.Writer) {
	c.cmd.Stderr = w
	c.cmd.Stdout = w
}

func (c *winCmd) Stop() error {
	return c.g.Dispose()
}
