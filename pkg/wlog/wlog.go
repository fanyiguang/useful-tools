package wlog

import (
	"os"
	"path/filepath"
	"strings"
	"useful-tools/common/config"

	"github.com/astaxie/beego/utils"

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
	NewLog(config.GetLogPath(), "debug")
}

func NewLog(logPath string, level string) {
	dir := strings.Trim(strings.Trim(logPath, "/"), "\\")
	if !utils.FileExists(dir) {
		_ = os.MkdirAll(dir, 0666)
	}
	logPath = filepath.Join(dir, "useful-tools.log")
	config := Config{
		Filename: logPath,
		Level:    getLogsLevel(level),
		MaxLines: 50000,
		MaxSize:  1 << 28,
		MaxDays:  3,
		Daily:    true,
	}
	jsonConfig, _ := json.Marshal(config)
	//_ = log.SetLogger(logs.AdapterConsole, string(jsonConfig))
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
