package sys

import (
	log "github.com/Sirupsen/logrus"
)

func initLogger() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "02-01-2006 15:04:05"
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
}

func Init() {
	log.Info("Initializing...")
	initLogger()
}
