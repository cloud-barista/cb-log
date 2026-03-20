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
	"os"
	"path/filepath"

	cblogformatter "github.com/cloud-barista/cb-log/formatter"
	"github.com/fsnotify/fsnotify"
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

// Get the logger with name you set. The name will be used as below (name: CB-SPIDER)
// [CB-SPIDER].[INFO]: 2020-12-24 16:54:46 sample-with-config-path.go:27, main.main() - start.........
// Read configuration file (log_conf.yaml) by the path set on environment variable (e.g., $CBLOG_ROOT)
func GetLogger(loggerName string) *logrus.Logger {
	return getLoggerHandler(loggerName, "")
}

// Read configuration file (log_conf.yaml) from the path you set
func GetLoggerWithConfigPath(loggerName string, configFilePath string) *logrus.Logger {
	return getLoggerHandler(loggerName, configFilePath)
}

// The handler for GetLogger() and GetLoggerWithConfigPath()
func getLoggerHandler(loggerName string, configFilePath string) *logrus.Logger {

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
	setup(loggerName, configFilePath)
	return thisLogger.logrus
}

func setup(loggerName string, configFilePath string) {
	cblogConfig = GetConfigInfos(configFilePath)
	thisLogger.logrus.SetReportCaller(true)

	if cblogConfig.CBLOG.LOOPCHECK {
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
		go levelSetupWatcher(loggerName, configFilePath)
	} else {
		SetLevel(cblogConfig.CBLOG.LOGLEVEL)
	}

	if cblogConfig.CBLOG.LOGFILE {
		setRotateFileHook(loggerName, &cblogConfig)
	}

	if !cblogConfig.CBLOG.CONSOLE {
		devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			logrus.Fatalf("Failed to open os.DevNull: %v", err)
		}
		thisLogger.logrus.SetOutput(devNull)
	} else {
		thisLogger.logrus.SetOutput(os.Stderr)
	}
}

// levelSetupWatcher watches the config file for changes using fsnotify
// and updates the log level whenever the file is modified.
// ref) https://github.com/fsnotify/fsnotify/blob/master/example_test.go
func levelSetupWatcher(loggerName string, configFilePath string) {
	// Resolve the config file path the same way GetConfigInfos does.
	watchPath := configFilePath
	if watchPath == "" {
		cblogRootPath := os.Getenv("CBLOG_ROOT")
		if cblogRootPath != "" {
			watchPath = filepath.Join(cblogRootPath, "conf", "log_conf.yaml")
		}
	}

	if watchPath == "" {
		logrus.Warn("[cb-log] LOOPCHECK is enabled but no config file path could be determined; file watcher will not start.")
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Errorf("[cb-log] Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(watchPath); err != nil {
		logrus.Errorf("[cb-log] Failed to watch config file %s: %v", watchPath, err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				cblogConfig = GetConfigInfos(configFilePath)
				SetLevel(cblogConfig.CBLOG.LOGLEVEL)
			}
			// Re-add the watch when the file is renamed/removed (e.g. atomic saves by vim/emacs).
			if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
				_ = watcher.Add(watchPath)
				cblogConfig = GetConfigInfos(configFilePath)
				SetLevel(cblogConfig.CBLOG.LOGLEVEL)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logrus.Errorf("[cb-log] File watcher error: %v", err)
		}
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

	level, err := logrus.ParseLevel(strLevel)
	if err != nil {
		thisLogger.logrus.Warnf("Not available logging level: %v. Default logging level will be used: debug", strLevel)
		level = logrus.DebugLevel
	}
	thisLogger.logrus.SetLevel(level)
}

func GetLevel() string {
	return thisLogger.logrus.GetLevel().String()
}

func getFormatter(loggerName string) *cblogformatter.Formatter {

	if thisFormatter != nil {
		return thisFormatter
	}
	thisFormatter = &cblogformatter.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[" + loggerName + "]." + "[%lvl%]: %time% %func% - %msg% \t[%keyvalues%]\n",
	}
	return thisFormatter
}
