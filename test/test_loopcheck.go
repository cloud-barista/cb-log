// Dynamic log-level test: demonstrates that editing loglevel in the config
// file is picked up automatically by the always-on file watcher.
//
// Usage:
//
//	export CBLOG_ROOT=<path-to-cb-log>
//	cd $CBLOG_ROOT/test
//	go run test_dynamic_loglevel.go
//
// While the program is running, edit $CBLOG_ROOT/conf/log_conf.yaml and
// change the 'loglevel' value. The new level will take effect immediately.
package main

import (
	"fmt"
	"time"

	cblog "github.com/cloud-barista/cb-log"
	"github.com/sirupsen/logrus"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
}

func main() {

	for {
		fmt.Printf("\n####LogLevel: %s\n", cblog.GetLevel())
		cblogger.Info("Log Info message")
		cblogger.Warning("Log Waring message")
		cblogger.Error("Log Error message")
		cblogger.Errorf("Log Error message:%s", errorMsg())
		time.Sleep(time.Second * 2)
	}
}

func errorMsg() error {
	return fmt.Errorf("internal error message")
}
