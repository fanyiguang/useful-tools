package path

import (
	"os"
	"os/user"
	"path"
)

func GetCurrentUserPath() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}

// GetAppDataPath 获取用户AddData路径
func GetAppDataPath() string {
	return path.Join(GetCurrentUserPath(), "AppData")
}

// GetAppDataLocalPath 获取用户AddData路径
func GetAppDataLocalPath() string {
	return path.Join(GetAppDataPath(), "Local")
}

func GetCurrentPath() string {
	userPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	return userPath
}
