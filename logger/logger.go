package logger

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Log() *logrus.Logger {
	if log != nil{
		return log
	}
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}