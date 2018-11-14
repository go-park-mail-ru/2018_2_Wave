package logger

import (
	"encoding/json"
	"os"
	"path"

	"go.uber.org/zap"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

var (
	dir     = "logs"
	log     = "wavelog.log"
	logFile = path.Join(dir, log)
)

// New logger
func New() *Logger {

	if _, err := os.Stat(logFile); err != nil {
		println(" -- log file: " + logFile)

		os.Mkdir(dir, 0755)
		file, err := os.Create(logFile)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "` + logFile + `"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	cfg := &zap.Config{}
	if err := json.Unmarshal(rawJSON, cfg); err != nil {
		panic(err)
	}

	basicLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		Sugar: basicLogger.Sugar(),
	}
}
