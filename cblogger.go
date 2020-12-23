// CB-Log: Logger for Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// ref) https://github.com/sirupsen/logrus
// ref) https://github.com/natefinch/lumberjack
// ref) https://github.com/snowzach/rotatefilehook
// by CB-Log Team, 2019.08.

package cblog

import (
	"fmt"
	"os"
	"strings"
	"time"

	cblogformatter "github.com/cloud-barista/cb-log/formatter"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type CBLogger struct {
	loggerName string
	logrus     *logrus.Logger
}

// global var.
var (
	thisLogger    *CBLogger
	thisFormatter *cblogformatter.Formatter
	cblogConfig   CBLOGCONFIG
)

func parseArgs(args []interface{}) (string, string){

	var loggerName string
	var configPath string

	for i, arg := range args {
		switch i {
		case 0: // name
			name, ok := arg.(string)
			if !ok {
				fmt.Printf("loggerName is not passed as string")
			}
			loggerName = name
		case 1:
			path, ok := arg.(string)
			if !ok {
				fmt.Printf("confgPath is not passed as string")
			}
			configPath = path
		default:
			fmt.Printf("Wrong parametes passed")
		}
	}
	return loggerName, configPath
}

// You can set up with Framework Name, a Framework Name is one of loggerName.
func GetLogger(args ...interface{}) *logrus.Logger {

	// arg[0]: loggerName
	// arg[1]: configPath

	loggerName, configPath := parseArgs(args)

	if thisLogger != nil {
		return thisLogger.logrus
	}
	thisLogger = new(CBLogger)
	thisLogger.loggerName = loggerName
	thisLogger.logrus = &logrus.Logger{
		Level:     logrus.DebugLevel,
		Out:       os.Stderr,
		Hooks:     make(logrus.LevelHooks),
		Formatter: getFormatter(loggerName),
	}

	// set config.
	setup(loggerName, configPath)
	return thisLogger.logrus
}

func setup(loggerName string, configPath string) {
	cblogConfig = GetConfigInfos(configPath)
	thisLogger.logrus.SetReportCaller(true)

	if cblogConfig.CBLOG.LOOPCHECK {
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
		go levelSetupLoop(loggerName, configPath)
	} else {
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
	}

	if cblogConfig.CBLOG.LOGFILE {
		setRotateFileHook(loggerName, &cblogConfig)
	}
}

// Now, this method is busy wait.
// @TODO must change this  with file watch&event.
// ref) https://github.com/fsnotify/fsnotify/blob/master/example_test.go
func levelSetupLoop(loggerName string, configPath string) {
	for {
		cblogConfig = GetConfigInfos(configPath)
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
		time.Sleep(time.Second * 2)
	}
}

func setRotateFileHook(loggerName string, logConfig *CBLOGCONFIG) {
	level, _ := logrus.ParseLevel(logConfig.CBLOG.LOGLEVEL)

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logConfig.LOGFILEINFO.FILENAME,
		MaxSize:    logConfig.LOGFILEINFO.MAXSIZE, // megabytes
		MaxBackups: logConfig.LOGFILEINFO.MAXBACKUPS,
		MaxAge:     logConfig.LOGFILEINFO.MAXAGE, //days
		Level:      level,
		Formatter:  getFormatter(loggerName),
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}
	thisLogger.logrus.AddHook(rotateFileHook)
}

func SetLevel(strLevel string) {
	err := checkLevel(strLevel)
	if err != nil {
		logrus.Errorf("Failed to set log level: %v", err)
	}
	level, _ := logrus.ParseLevel(strLevel)
	thisLogger.logrus.SetLevel(level)
}

func checkLevel(lvl string) error {
	switch strings.ToLower(lvl) {
	case "error":
		return nil
	case "warn", "warning":
		return nil
	case "info":
		return nil
	case "debug":
		return nil
	}
	return fmt.Errorf("not a valid cblog Level: %q", lvl)
}

func GetLevel() string {
	return thisLogger.logrus.GetLevel().String()
}

func getFormatter(loggerName string) *cblogformatter.Formatter {

	if thisFormatter != nil {
		return thisFormatter
	}
	// 출력 포맷 조정 (keyvalues) 추가 (Formatter.go에서 해당 위치에 실제 데이터로 변경)
	thisFormatter = &cblogformatter.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[" + loggerName + "]." + "[%lvl%]: %time% %func% - %msg% \t[%keyvalues%]\n",
	}
	return thisFormatter
}
