package log

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

func Init(path, level string) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Warn(err.Error())
		l = logrus.InfoLevel
	}
	logrus.SetLevel(l)
	logger := &lumberjack.Logger{
		LocalTime: true,
		Filename:  path,
		MaxSize:   20, // 一个文件最大为20M
		MaxAge:    10, // 一个文件最多同时存在10天
		Compress:  false,
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	mutWriter := []io.Writer{logger}
	if level == "debug" {
		mutWriter = append(mutWriter, os.Stdout)
	}
	fileAndStdoutWriter := io.MultiWriter(mutWriter...)
	logrus.SetOutput(fileAndStdoutWriter)
}

func New(path, level string) *logrus.Logger {
	logger := &lumberjack.Logger{
		LocalTime: true,
		Filename:  path,
		MaxSize:   20, // 一个文件最大为20M
		MaxAge:    10, // 一个文件最多同时存在10天
		Compress:  false,
	}
	l := logrus.New()
	l.SetOutput(logger)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Warn(err.Error())
		logrusLevel = logrus.InfoLevel
	}
	l.SetLevel(logrusLevel)
	return l
}
