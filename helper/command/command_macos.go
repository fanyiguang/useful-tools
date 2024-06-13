//go:build darwin

package command

import (
	"io"
	"os/exec"
)

type unixCmd struct {
	cmd *exec.Cmd
}

func (c *unixCmd) Start() error {
	return c.cmd.Start()
}

func (c *unixCmd) Wait() error {
	return c.cmd.Wait()
}

func (c *unixCmd) Run() error {
	return c.cmd.Run()
}

func (c *unixCmd) SetOutput(w io.Writer) {
	c.cmd.Stderr = w
	c.cmd.Stdout = w
}

func (c unixCmd) Stop() error {
	return c.cmd.Process.Kill()
}

func New(name string, args ...string) (Cmd, error) {
	return &unixCmd{
		cmd: exec.Command(name, args...),
	}, nil
}
