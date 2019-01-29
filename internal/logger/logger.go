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
	Info(args ...interface{})
}

// ----------------| Logger

// Logger - default ilogger implementation
type Logger struct {
	Sugar *zap.SugaredLogger
}

func logfileExists(logPath, logFile string) bool {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, 0777)
		_, err := os.OpenFile(logPath+logFile, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return false
		}
		return true
	} else if !os.IsNotExist(err) {
		return true
	}

	return false
}

// Construct - constructor
func Construct(logPath, logFile string) *Logger {
	if !logfileExists(logPath, logFile) {
		return &Logger{}
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
	var sugarredLogger Logger

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return &Logger{}
	}

	basicLogger, err := cfg.Build()
	sugarredLogger.Sugar = basicLogger.Sugar()
	if err != nil {
		return &Logger{}
	}

	return &sugarredLogger
}

// Infof - inforamtion patterned message
func (l *Logger) Infof(pattern string, args ...interface{}) {
	l.Sugar.Infof(pattern, args...)
}

// Info - inforamtion message
func (l *Logger) Info(args ...interface{}) {
	l.Sugar.Info(args...)
}
