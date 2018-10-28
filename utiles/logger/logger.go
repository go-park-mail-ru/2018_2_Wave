package logger

import (
	"errors"
	"io"
	"os"
	"sync"

	"Wave/utiles/walhalla"
	"github.com/sirupsen/logrus"
)

// Log just a log
type Log struct {
	mu      *sync.RWMutex // mutex to protect shared members
	entery  logrus.Entry  // entery
	out     *os.File      // output file
	bClosed *bool         // weather the log was closed
	bAsync  bool          // async writing

	walhalla.ILogger
}

type Config struct {
	File    string // Ouput file ("" => no output file)
	BStdOut bool   // copy the logs to stdOut
	BStdErr bool   // copy the logs to stdErr
	BAsync  bool   // async writing
}

// New - create a new Log
func New(config Config) (log *Log, err error) {
	entery := logrus.WithFields(logrus.Fields{
		// "": "",
	})
	entery.Logger.SetFormatter(&logrus.JSONFormatter{})

	var out *os.File
	{ // setup output
		if config.File != "" {
			if out, err = os.OpenFile(config.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
				return nil, err
			}
		}
		writers := []io.Writer{}
		if out != nil {
			writers = append(writers, out)
		}
		if config.BStdErr {
			writers = append(writers, os.Stderr)
		}
		if config.BStdOut {
			writers = append(writers, os.Stdout)
		}
		entery.Logger.SetOutput(io.MultiWriter(writers...))
	}

	return &Log{
		bClosed: new(bool),
		bAsync:  config.BAsync,
		entery:  *entery,
		out:     out,
		mu:      &sync.RWMutex{},
	}, nil
}

func (log *Log) Close() {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.out != nil {
		*log.bClosed = true
		log.out.Close()
	}
}

func (log *Log) WithFields(fields walhalla.Fields) walhalla.ILogger {
	fds := logrus.Fields(fields)
	return &Log{
		entery:  *log.entery.WithFields(fds),
		out:     log.out,
		mu:      log.mu,
		bAsync:  log.bAsync,
		bClosed: log.bClosed,
	}
}

// NOTE: unsafe
func (log *Log) mustBeValid() {
	if *log.bClosed {
		panic(errors.New("Closed log mustn't be used"))
	}
}

// ----------------| Info

func (log *Log) Info(args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entery.Infoln(args...)
	} else {
		log.entery.Infoln(args...)
	}
}

func (log *Log) Infof(format string, args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entery.Infof(format, args...)
	} else {
		log.entery.Infof(format, args...)
	}
}

// ----------------| Warn

func (log *Log) Warn(args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entery.Warnln(args...)
	} else {
		log.entery.Warnln(args...)
	}
}

func (log *Log) Warnf(format string, args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entery.Warnf(format, args...)
	} else {
		log.entery.Warnf(format, args...)
	}
}

// ----------------| Error

func (log *Log) Error(args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entery.Errorln(args...)
	} else {
		log.entery.Errorln(args...)
	}
}

func (log *Log) Errorf(format string, args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entery.Errorf(format, args...)
	} else {
		log.entery.Errorf(format, args...)
	}
}
