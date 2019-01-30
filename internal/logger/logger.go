package logger

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

// ----------------| ILogger

// ILogger - application logger interface
type ILogger interface {
	Infof(pattern string, args ...interface{})
	Info(msg string, args ...interface{})
}

// ----------------| Logger

// Logger - default ilogger implementation
type Logger struct {
	Sugar *zap.SugaredLogger
}

func logfileExists(logPath, logFile string) bool {
	if _, err := os.Stat(logPath + logFile); err != nil {
		os.Mkdir(logPath, 0777)
		f, err := os.Create(logPath + logFile)
		if err != nil {
			return false
		}
		f.Chmod(0777)
		f.Close()
		return true
	}
	return true
}

// Construct - constructor
func Construct(logPath, logFile string) *Logger {
	if !logfileExists(logPath, logFile) {
		panic("it's impossible to open or create the log file (" + logPath + logFile + ")")
	}

	JSON := (`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "` + logPath + logFile + `"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	rawJSON := []byte(JSON)

	var cfg zap.Config
	var err error

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return &Logger{}
	}

	basicLogger, err := cfg.Build()
	if err != nil {
		return &Logger{}
	}
	sugarredLogger := &Logger{}
	sugarredLogger.Sugar = basicLogger.Sugar()
	sugarredLogger.Info("Logger started",
		"Who", "Construct",
		"Where", "logger.go")
	return sugarredLogger
}

// Infof - inforamtion patterned message
func (l *Logger) Infof(pattern string, args ...interface{}) {
	l.Sugar.Infof(pattern, args...)
}

// Info - inforamtion message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Sugar.Infow(msg, args...)
}
