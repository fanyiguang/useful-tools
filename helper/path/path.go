package path

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
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

func Path() (string, error) {
	// 获取可执行文件的路径
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 将路径转换为绝对路径
	absolutePath, err := filepath.Abs(executablePath)
	if err != nil {
		return "", err
	}

	return absolutePath, nil
}
