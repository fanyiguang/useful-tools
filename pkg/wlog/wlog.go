package wlog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	log  *logs.BeeLogger
)

type Config struct {
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	MaxLines int    `json:"maxlines"`
	MaxSize  int    `json:"maxsize"`
	MaxDays  int    `json:"maxdays"`
	Daily    bool   `json:"daily"`
}

func init() {
	log = logs.NewLogger()
	executable, _ := os.Executable()
	NewLog(filepath.Dir(executable), "debug")
}

func NewLog(logPath string, level string) {
	dir := strings.Trim(strings.Trim(logPath, "/"), "\\")
	if !IsExist(dir) {
		_ = os.MkdirAll(dir, 0666)
	}
	logPath = filepath.Join(dir, "parser-tool.log")
	config := Config{
		Filename: logPath,
		Level:    getLogsLevel(level),
		MaxLines: 50000,
		MaxSize:  1 << 28,
		MaxDays:  3,
		Daily:    true,
	}
	jsonConfig, _ := json.Marshal(config)
	_ = log.SetLogger(logs.AdapterFile, string(jsonConfig))
	log.EnableFuncCallDepth(true)
}

func getLogsLevel(level string) (l int) {
	switch level {
	case "debug":
		l = logs.LevelDebug
	case "warning":
		l = logs.LevelWarning
	case "warn":
		l = logs.LevelWarn
	case "info":
		l = logs.LevelInfo
	default:
		l = logs.LevelDebug
	}
	return
}

func Warm(format string, args ...interface{}) {
	log.Warn(format, args)
}

func Error(format string, args ...interface{}) {
	log.Error(format, args)
}

func Debug(format string, args ...interface{}) {
	log.Debug(format, args)
}

func Info(format string, args ...interface{}) {
	log.Info(format, args)
}

func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	}
	//mac
	if strings.Contains(err.Error(), "file name too long") {
		return false
	}
	return !os.IsNotExist(err)
}
