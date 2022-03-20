package config

import (
	"os"
	"path/filepath"
	"useful-tools/helper/path"

	"github.com/astaxie/beego/utils"
)

var (
	projectsPath = ""
	tempPath     = ""
	logPath      = ""
	settingPath  = ""
)

func init() {
	localPath := path.GetAppDataLocalPath()
	SetProjectsPath(filepath.Join(localPath, "useful-tools"))
	SetTempPath(filepath.Join(GetProjectsPath(), "temp"))
	SetLogPath(filepath.Join(GetProjectsPath(), "log"))
	SetSettingPath(filepath.Join(GetProjectsPath(), "setting"))
	//fmt.Println(GetTempPath(), GetLogPath(), GetProjectsPath(), GetSettingPath())
	if !utils.FileExists(GetTempPath()) {
		_ = os.MkdirAll(GetTempPath(), 0666)
	}

	if !utils.FileExists(GetLogPath()) {
		_ = os.MkdirAll(GetLogPath(), 0666)
	}

	if !utils.FileExists(GetSettingPath()) {
		_ = os.MkdirAll(GetSettingPath(), 0666)
	}
}

func GetProjectsPath() string {
	return projectsPath
}

func SetProjectsPath(path string) {
	projectsPath = path
}

func GetTempPath() string {
	return tempPath
}

func SetTempPath(path string) {
	tempPath = path
}

func GetLogPath() string {
	return logPath
}

func SetLogPath(path string) {
	logPath = path
}

func GetSettingPath() string {
	return settingPath
}

func SetSettingPath(path string) {
	settingPath = path
}
