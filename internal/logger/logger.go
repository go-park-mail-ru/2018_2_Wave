package logger

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

const (
	path    = "./logs/"
	logFile = "wavelog"
)

func logfileExists() bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		_, err := os.OpenFile(path+logFile, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return false
		}
		return true
	} else if !os.IsNotExist(err) {
		return true
	}

	return false
}

func Construct() *Logger {
	if !logfileExists() {
		return &Logger{}
	}

	rawJSON := []byte(`{
	"level": "debug",
	"encoding": "json",
	"outputPaths": ["stdout", "./logs/wavelog"],
	"errorOutputPaths": ["stderr"],
	"encoderConfig": {
	  "messageKey": "message",
	  "levelKey": "level",
	  "levelEncoder": "lowercase"
	}
  }`)

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
