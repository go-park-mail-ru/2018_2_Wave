package logger

import (
	"encoding/json"
	"log"
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

const path = "./logs/wavelog"

func logfileExists() bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			panic(err)
		}
		return true
	} else if !os.IsNotExist(err) { //ахуенна
		return true
	}

	return false
}

func Construct() Logger {
	if !logfileExists() {
		log.Println("Caution: logger output file missing, no logging utility set.")

		return Logger{}
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

	//"initialFields": {"source": "undefined"},
	var cfg zap.Config
	var err error
	var sugarredLogger Logger

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	basicLogger, err := cfg.Build()
	sugarredLogger.Sugar = basicLogger.Sugar()
	if err != nil {
		panic(err)
	}

	//defer l.Sugar.Sync()
	sugarredLogger.Sugar.Infow("gofuck tou",
		"source", "logger.go")
	return sugarredLogger
}
