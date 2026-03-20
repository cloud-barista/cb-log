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
"context"
"os"
"path/filepath"
"time"

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
watcherCancel context.CancelFunc // stops the background watcher goroutine
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

SetLevel(cblogConfig.CBLOG.LOGLEVEL)
ctx, cancel := context.WithCancel(context.Background())
watcherCancel = cancel
go levelSetupWatcher(ctx, loggerName, configFilePath)

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
// It stops when ctx is cancelled.
//
// To handle atomic saves (used by vim, emacs, sed -i, etc.) the watcher
// monitors the parent *directory* of the config file. Atomic saves
// remove/rename the original file and replace it with a new one, which
// would break a direct file watch. Watching the directory ensures we
// always receive the subsequent Create event when the new file appears.
// ref) https://github.com/fsnotify/fsnotify/blob/master/example_test.go
func levelSetupWatcher(ctx context.Context, loggerName string, configFilePath string) {
// Resolve the config file path the same way GetConfigInfos does.
watchPath := configFilePath
if watchPath == "" {
cblogRootPath := os.Getenv("CBLOG_ROOT")
if cblogRootPath != "" {
watchPath = filepath.Join(cblogRootPath, "conf", "log_conf.yaml")
}
}

if watchPath == "" {
logrus.Warn("[cb-log] No config file path could be determined; file watcher will not start.")
return
}

// Watch the parent directory instead of the file itself so that
// atomic-save operations (remove + create) are handled correctly.
watchDir := filepath.Dir(watchPath)
watchBase := filepath.Base(watchPath)

watcher, err := fsnotify.NewWatcher()
if err != nil {
logrus.Errorf("[cb-log] Failed to create file watcher: %v", err)
return
}
defer watcher.Close()

if err := watcher.Add(watchDir); err != nil {
logrus.Errorf("[cb-log] Failed to watch config directory %s: %v", watchDir, err)
return
}

// debounce timer to coalesce rapid events from atomic saves.
var debounceTimer *time.Timer

for {
select {
case <-ctx.Done():
return
case event, ok := <-watcher.Events:
if !ok {
return
}
// Only react to events on our config file.
if filepath.Base(event.Name) != watchBase {
continue
}
if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
// Debounce: wait briefly for any follow-up events to settle.
if debounceTimer != nil {
debounceTimer.Stop()
}
debounceTimer = time.AfterFunc(200*time.Millisecond, func() {
reloadConfig(configFilePath)
})
}
// Rename/Remove events are expected during atomic saves;
// the directory watch survives them, and the subsequent
// Create event (handled above) will trigger the reload.
case err, ok := <-watcher.Errors:
if !ok {
return
}
logrus.Errorf("[cb-log] File watcher error: %v", err)
}
}
}

// reloadConfig re-reads the config file and applies the new log level.
// Errors are logged and silently ignored so that the process is not
// terminated by a transient file-system condition (e.g. mid-atomic-save).
func reloadConfig(configFilePath string) {
newConfig, err := GetConfigInfosSafe(configFilePath)
if err != nil {
logrus.Warnf("[cb-log] Failed to reload config (keeping current settings): %v", err)
return
}
cblogConfig = newConfig
SetLevel(cblogConfig.CBLOG.LOGLEVEL)
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
