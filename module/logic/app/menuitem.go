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
	saveAesKey bool
	aesKey     string
	aesIV      string
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
	if s := m.GetSaveAesKeyFromFile(); s == 1 {
		m.saveAesKey = true
		m.aesKey = m.GetAesKeyFromFile()
		m.aesIV = m.GetAesIVFromFile()
	}
	return m
}

func (m *MenuItem) AesKey() string {
	return m.aesKey
}

func (m *MenuItem) SetAesKey(aesKey string) {
	m.aesKey = aesKey
}

func (m *MenuItem) AesIV() string {
	return m.aesIV
}

func (m *MenuItem) SetAesIV(aesIV string) {
	m.aesIV = aesIV
}

func (m *MenuItem) SaveAesKey() bool {
	return m.saveAesKey
}

func (m *MenuItem) SetSaveAesKey(saveAesKey bool) {
	m.saveAesKey = saveAesKey
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
	m.SetStateToFile("view_mode", strconv.Itoa(state))
}

func (m *MenuItem) GetProModeFromFile() (state int) {
	s, err := strconv.Atoi(m.GetStateFromFile("view_mode"))
	if err != nil {
		return 0
	}
	return s
}

func (m *MenuItem) SetShowPassToFile(state int) {
	m.SetStateToFile("show_pass", strconv.Itoa(state))
}

func (m *MenuItem) GetShowPassFromFile() (state int) {
	s, err := strconv.Atoi(m.GetStateFromFile("show_pass"))
	if err != nil {
		return 0
	}
	return s

}

func (m *MenuItem) SetSaveAesKeyToFile(state int) {
	m.SetStateToFile("save_aes_key", strconv.Itoa(state))
}

func (m *MenuItem) GetSaveAesKeyFromFile() (state int) {
	s, err := strconv.Atoi(m.GetStateFromFile("save_aes_key"))
	if err != nil {
		return 0
	}
	return s
}

func (m *MenuItem) SetAesKeyToFile(content string) {
	m.SetStateToFile("aes_key", content)
}

func (m *MenuItem) GetAesKeyFromFile() (content string) {
	return m.GetStateFromFile("aes_key")
}

func (m *MenuItem) SetAesIVToFile(content string) {
	m.SetStateToFile("aes_iv", content)
}

func (m *MenuItem) GetAesIVFromFile() (content string) {
	return m.GetStateFromFile("aes_iv")
}

func (m *MenuItem) SetHiddenBodyToFile(state int) {
	m.SetStateToFile("hidden_body", strconv.Itoa(state))
}

func (m *MenuItem) GetHiddenBodyFromFile() (state int) {
	s, err := strconv.Atoi(m.GetStateFromFile("hidden_body"))
	if err != nil {
		return 0
	}
	return s
}

func (m *MenuItem) SetStateToFile(name string, content string) {
	file, err := os.OpenFile(filepath.Join(config.GetSettingPath(), name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		wlog.Warm("os.Open view_mode failed: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		wlog.Warm("file.WriteString state failed: %v", err)
		return
	}
}

func (m *MenuItem) GetStateFromFile(name string) (content string) {
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

	return string(tempState)
}
