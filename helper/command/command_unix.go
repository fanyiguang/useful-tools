//go:build linux

package command

import (
	"io"
	"os/exec"
	"syscall"
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
	c.cmd.Stdout = w
	c.cmd.Stderr = w
}

func (c unixCmd) Stop() error {
	return c.cmd.Process.Kill()
}

func New(name string, args ...string) (Cmd, error) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
	}
	return &unixCmd{
		cmd: cmd,
	}, nil
}
