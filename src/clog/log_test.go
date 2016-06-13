package clog

import (
	"testing"
)

var info = LogInfo{
	Mode:   "file",
	Level:  "info",
	Buffer: 100000,

	FileName:     "test.log",
	LogRotate:    true,
	MaxLines:     10000,
	MaxSizeShift: 28,
	DailyRotate:  true,
	MaxDays:      7,
}

func TestInitLog(t *testing.T) {
	InitLog(info)
}
