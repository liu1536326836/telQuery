package clog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	blog "github.com/weisd/log"
)

var logger *blog.Logger

var logLevels = map[string]string{
	"trace":    "0",
	"debug":    "1",
	"info":     "2",
	"warn":     "3",
	"error":    "4",
	"critical": "5",
}

type LogInfo struct {
	Mode   string
	Level  string
	Buffer int64

	FileName     string
	LogRotate    bool
	MaxLines     int
	MaxSizeShift int
	DailyRotate  bool
	MaxDays      int
}

func InitLog(loginfo LogInfo) {
	str := ""

	level := logLevels[strings.ToLower(loginfo.Level)]
	switch loginfo.Mode {
	case "console":
		str = fmt.Sprintf(`{"level":%s}`, level)

	case "file":
		if len(loginfo.FileName) > 0 {
			logpath, _ := filepath.Abs(loginfo.FileName)
			os.MkdirAll(path.Dir(logpath), os.ModePerm)
		}

		str = fmt.Sprintf(
			`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`,
			level,
			loginfo.FileName,
			loginfo.LogRotate,
			loginfo.MaxLines,
			1<<uint(loginfo.MaxSizeShift),
			loginfo.DailyRotate,
			loginfo.MaxDays,
		)

	default:
		fmt.Printf("未知的日志类型, %s", loginfo.Mode)
		os.Exit(1)
	}

	logger = blog.NewCustomLogger(loginfo.Buffer, loginfo.Mode, str)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Error(4, format, v...)
}
