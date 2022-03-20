package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"useful-tools/common/config"
	"useful-tools/helper/Go"
	"useful-tools/pkg/wlog"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) SetViewModeState(state int) {
	Go.Go(func() {
		file, err := os.OpenFile(filepath.Join(config.GetSettingPath(), "view_mode"), os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			wlog.Warm("os.Open view_mode failed: %v", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(strconv.Itoa(state))
		if err != nil {
			wlog.Warm("file.WriteString state failed: %v", err)
			return
		}
	})
}

func (a *App) GetViewModeState() (state int) {
	file, err := os.OpenFile(filepath.Join(config.GetSettingPath(), "view_mode"), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		wlog.Warm("os.Open view_mode failed: %v", err)
		return
	}
	defer file.Close()

	var tempState []byte
	tempState, err = ioutil.ReadAll(file)
	if err != nil {
		wlog.Warm("file.WriteString state failed: %v", err)
		return
	}

	state, _ = strconv.Atoi(string(tempState))
	return
}
