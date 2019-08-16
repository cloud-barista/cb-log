// CB-Log: Logger for Cloud-Barista.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// ref) https://github.com/sirupsen/logrus
// ref) https://github.com/natefinch/lumberjack
// ref) https://github.com/snowzach/rotatefilehook
// by powerkim@etri.re.kr, 2019.08.

package cblog


import (
	"os"
	"fmt"
	"strings"

        "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/cloud-barista/cb-log/formatter"
)

type CBLogger struct {
	loggerName string
	logrus *logrus.Logger
}

// global var.
var (
	thisLogger *CBLogger
	thisFormatter *cblogformatter.Formatter
	logFilePath string = "./log/logs.log"
)

// CB-Log's Logging Level type
//type Level uint32
type Level logrus.Level

// CB-Log's logging level to log
const (
        ErrorLevel Level = Level(logrus.ErrorLevel)
        WarnLevel Level = Level(logrus.WarnLevel)
        InfoLevel Level = Level(logrus.InfoLevel)
)


/*********** ref) logrus's logging level
// These are the different logging levels. You can set the logging level to log
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)
***********/

// You can set up with Framework Name, a Framework Name is one of loggerName.
//func (cbLogger CBLogger)GetLogger(loggerName string) *CBLogger {
func GetLogger(loggerName string) *logrus.Logger {
	if thisLogger != nil {
		return thisLogger.logrus
	}
	thisLogger = new(CBLogger)
	thisLogger.loggerName = loggerName
	thisLogger.logrus =  &logrus.Logger{
        Out:   os.Stderr,
        Level: logrus.DebugLevel,
        Hooks: make(logrus.LevelHooks),
        Formatter: getFormatter(loggerName),
        //Formatter: &cblogformatter.Formatter{
        //    TimestampFormat: "2006-01-02 15:04:05",
        //    LogFormat:       "[" + loggerName + "]." + "[%lvl%]: %time% %func% - %msg%\n",
	//},
	}


	// set default config.
	thisLogger.logrus.SetReportCaller(true)
	SetLevel("info")
	setRotateFileHook(loggerName)
	setRotateFileHook2(loggerName)
//fmt.Printf("====> %#v\n", thisLogger.logrus.Hooks.AllLevels())
	return thisLogger.logrus
}

func SetLevel(strLevel string) {
	err := checkLevel(strLevel)
	if err != nil {
                logrus.Errorf("Failed to set log level: %v", err)
	}
	level, _ := logrus.ParseLevel(strLevel)
	thisLogger.logrus.SetLevel(level)
}

func checkLevel(lvl string) (error) {
	switch strings.ToLower(lvl) {
	case "error":
		return nil
	case "warn", "warning":
		return nil
	case "info":
		return nil
	}
	return fmt.Errorf("not a valid cblog Level: %q", lvl)
}

/* deprecated
func SetLevel(level Level) {
	thisLogger.logrus.SetLevel(logrus.Level(level))	
}
*/

func GetLevel() string {
	return thisLogger.logrus.GetLevel().String()
}

/* deprecated
func GetLevel() Level {
	return Level(thisLogger.logrus.GetLevel())
}
*/

func (level Level) String() string {
	return logrus.Level(level).String()
}

func getFormatter(loggerName string) *cblogformatter.Formatter {

	if thisFormatter != nil {
		return thisFormatter
	}
	thisFormatter = &cblogformatter.Formatter{
            TimestampFormat: "2006-01-02 15:04:05",
            LogFormat:       "[" + loggerName + "]." + "[%lvl%]: %time% %func% - %msg%\n",
        }	
	return thisFormatter
}

func setRotateFileHook(loggerName string) {
        rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
                Filename:   logFilePath,
                //MaxSize:    50, // megabytes
                MaxSize:    10, // megabytes
                MaxBackups: 50,
                MaxAge:     31, //days
                Level:      logrus.InfoLevel,
                //Level:      logrus.ErrorLevel,
                Formatter: getFormatter(loggerName),
                //Formatter: &logrus.JSONFormatter{
                //        TimestampFormat: time.RFC822,
                //},
                //Formatter: (&logrus.TextFormatter{
                //        DisableColors:   true,
                //        ForceColors:   true,
                //        FullTimestamp: true,
                //}),
        })


        if err != nil {
                logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
        }
        thisLogger.logrus.AddHook(rotateFileHook)
}

func setRotateFileHook2(loggerName string) {
        rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
                Filename:   logFilePath+"-2",
                //MaxSize:    50, // megabytes
                MaxSize:    10, // megabytes
                MaxBackups: 50,
                MaxAge:     31, //days
                //Level:      logrus.InfoLevel,
                Level:      logrus.ErrorLevel,
                Formatter: getFormatter(loggerName),
                //Formatter: &logrus.JSONFormatter{
                //        TimestampFormat: time.RFC822,
                //},
                //Formatter: (&logrus.TextFormatter{
                //        DisableColors:   true,
                //        ForceColors:   true,
                //        FullTimestamp: true,
                //}),
        })


        if err != nil {
                logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
        }
        thisLogger.logrus.AddHook(rotateFileHook)
}
