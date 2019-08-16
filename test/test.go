package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/cloud-barista/cb-log"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
	cblog.SetLevel("error")
	cblog.SetLevel("warn")
	cblog.SetLevel("info")
}

func main() {

	fmt.Printf("####LogLevel: %s\n", cblog.GetLevel())
	cblogger.Info("Log Info message")
	cblogger.Infof("Log Info message:%s", "abc")
	cblogger.Warningln("Log Waring message")
	cblogger.Errorln("Log Error message")

	cblog.SetLevel("error")
	fmt.Printf("####LogLevel: %s\n", cblog.GetLevel())

	cblogger.Info("Log Info message")
	cblogger.Infof("Log Info message:%s", "abc")
	cblogger.Warningln("Log Waring message")
	cblogger.Errorln("Log Error message")
}
