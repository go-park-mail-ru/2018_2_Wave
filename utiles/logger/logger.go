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
	entry  	logrus.Entry  // entry
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
	entry := logrus.WithFields(logrus.Fields{
		// "": "",
	})
	entry.Logger.SetFormatter(&logrus.JSONFormatter{})

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
		entry.Logger.SetOutput(io.MultiWriter(writers...))
	}

	return &Log{
		bClosed: new(bool),
		bAsync:  config.BAsync,
		entry:  *entry,
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
		entry:  *log.entry.WithFields(fds),
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
		go log.entry.Infoln(args...)
	} else {
		log.entry.Infoln(args...)
	}
}

func (log *Log) Infof(format string, args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entry.Infof(format, args...)
	} else {
		log.entry.Infof(format, args...)
	}
}

// ----------------| Warn

func (log *Log) Warn(args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entry.Warnln(args...)
	} else {
		log.entry.Warnln(args...)
	}
}

func (log *Log) Warnf(format string, args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entry.Warnf(format, args...)
	} else {
		log.entry.Warnf(format, args...)
	}
}

// ----------------| Error

func (log *Log) Error(args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entry.Errorln(args...)
	} else {
		log.entry.Errorln(args...)
	}
}

func (log *Log) Errorf(format string, args ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	log.mustBeValid()
	if log.bAsync {
		go log.entry.Errorf(format, args...)
	} else {
		log.entry.Errorf(format, args...)
	}
}
