package Go

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

func RelativelySafeGo(fun func()) {
	go func(f func()) {
		defer func() {
			if err := recover(); err != nil {
				Recover(err)
			}
		}()

		f()
	}(fun)
}

func Recover(err interface{}) {
	callers := make([]uintptr, 15)
	_ = runtime.Callers(3, callers)
	for k, _ := range callers {
		frame, _ := runtime.CallersFrames(callers[k : k+1]).Next()
		//frame, _ := runtime.CallersFrames(callers).Next()
		logrus.Errorf("runtime.CallersFrames failed, err: %v file: %v line: %v func: %v", err, frame.File, frame.Line, frame.Function)
	}
	logrus.Error(strings.Repeat("-", 100))
}
