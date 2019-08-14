package main

import (
	"github.com/sirupsen/logrus"
	"github.com/cloud-barista/cb-log"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
	cblog.SetLevel(cblog.ErrorLevel)
	cblog.SetLevel(cblog.WarnLevel)
	cblog.SetLevel(cblog.InfoLevel)
}

func main() {

	cblogger.Info("Log Info message")
	cblogger.Infof("Log Info message:%s", "abc")
	cblogger.Warningln("Log Waring message")
	cblogger.Errorln("Log Error message")

cblogger.Info("\n\n")

	cblog.SetLevel(cblog.ErrorLevel)

	cblogger.Info("Log Info message")
	cblogger.Infof("Log Info message:%s", "abc")
	cblogger.Warningln("Log Waring message")
	cblogger.Errorln("Log Error message")
}
