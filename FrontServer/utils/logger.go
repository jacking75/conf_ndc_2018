package utils

import (
	logrus "github.com/sirupsen/logrus"
)

var (
	Logger = logrus.New()
)

func SettingLog() {
	Logger.SetLevel(logrus.DebugLevel)
}
