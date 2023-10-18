package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"useful-tools/common/config"
	"useful-tools/pkg/wlog"
)

type MenuItem struct {
	proMode    bool
	showPass   bool
	hiddenBody bool
}

func NewMenuItem() *MenuItem {
	m := &MenuItem{}
	if s := m.GetProModeFromFile(); s == 1 {
		m.proMode = true
	}
	if s := m.GetShowPassFromFile(); s == 1 {
		m.showPass = true
	}
	if s := m.GetHiddenBodyFromFile(); s == 1 {
		m.hiddenBody = true
	}
	return m
}

func (m *MenuItem) ProMode() bool {
	return m.proMode
}

func (m *MenuItem) SetProMode(proMode bool) {
	m.proMode = proMode
}

func (m *MenuItem) ShowPass() bool {
	return m.showPass
}

func (m *MenuItem) SetShowPass(showPass bool) {
	m.showPass = showPass
}

func (m *MenuItem) HiddenBody() bool {
	return m.hiddenBody
}

func (m *MenuItem) SetHiddenBody(hiddenBody bool) {
	m.hiddenBody = hiddenBody
}

func (m *MenuItem) SetProModeToFile(state int) {
	m.SetStateToFile("view_mode", state)
}

func (m *MenuItem) GetProModeFromFile() (state int) {
	return m.GetStateFromFile("view_mode")
}

func (m *MenuItem) SetShowPassToFile(state int) {
	m.SetStateToFile("show_pass", state)
}

func (m *MenuItem) GetShowPassFromFile() (state int) {
	return m.GetStateFromFile("show_pass")
}

func (m *MenuItem) SetHiddenBodyToFile(state int) {
	m.SetStateToFile("hidden_body", state)
}

func (m *MenuItem) GetHiddenBodyFromFile() (state int) {
	return m.GetStateFromFile("hidden_body")
}

func (m *MenuItem) SetStateToFile(name string, state int) {
	file, err := os.OpenFile(filepath.Join(config.GetSettingPath(), name), os.O_CREATE|os.O_RDWR, 0666)
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
}

func (m *MenuItem) GetStateFromFile(name string) (state int) {
	file, err := os.OpenFile(filepath.Join(config.GetSettingPath(), name), os.O_CREATE|os.O_RDWR, 0666)
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
