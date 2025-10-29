//go:build windows
// +build windows

package proc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CodyGuo/win"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	winExecError = map[uint32]string{
		0:  "The system is out of memory or resources.",
		2:  "The .exe file is invalid.",
		3:  "The specified file was not found.",
		11: "The specified path was not found.",
	}
)

func RunProcByWin32Api(cmd string, isStart bool, args ...string) error {
	var temCmd, strArgs string
	if len(args) > 0 {
		strArgs = strings.Join(args, " ")
	}
	if strArgs == "" {
		if isStart {
			temCmd = fmt.Sprintf("cmd /c start %v", cmd)
		} else {
			temCmd = fmt.Sprintf("cmd /c %v", cmd)
		}
	} else {
		if isStart {
			temCmd = fmt.Sprintf("start \"%v\" %v", cmd, strArgs)
		} else {
			temCmd = fmt.Sprintf("%v %v", cmd, strArgs)
		}
	}

	s, err := simplifiedchinese.GB18030.NewEncoder().String(temCmd)
	if err != nil {
		return err
	}
	lpCmdLine := win.StringToBytePtr(s)
	ret := win.WinExec(lpCmdLine, win.SW_HIDE)
	//ret := win.WinExec(lpCmdLine, win.SW_SHOWMINIMIZED)
	if ret <= 31 {
		return errors.New(winExecError[ret])
	}
	return nil
}
