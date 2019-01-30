package proto

import (
	"Wave/internal/logger"
)

// ----------------| ActorTask

// ActorTask - message an actor
type ActorTask func()

// ----------------| IActor

// IActor - Actor inteface.
// @See Actor pattern
type IActor interface {
	Task(ActorTask) error
	SetLogger(logger.ILogger)
	GetLogger() logger.ILogger
}

// ----------------| Actor

// Actor - IActor interface default realisation
type Actor struct {
	T       chan ActorTask
	LG      logger.ILogger
	LogMeta func() []interface{}
	OnPanic func(interface{}) bool // true - continue the code
}

// MakeActor - constructor
func MakeActor(bufferSize int) Actor {
	return Actor{
		T: make(chan ActorTask, bufferSize),
	}
}

// AddTask to the actor message queue
func (a *Actor) Task(t ActorTask) error {
	if t == nil {
		a.T <- t
		return nil
	}
	return ErrorNil
}

// ->> log

// Logf - patterned log
func (a *Actor) Logf(pattern string, args ...interface{}) {
	if a.LG == nil {
		return
	}
	if a.LogMeta != nil {
		args = append(args, a.LogMeta()...)
	}
	a.LG.Infof(pattern, args...)
}

// Log - log
func (a *Actor) Log(msg string, args ...interface{}) {
	if a.LG == nil {
		return
	}
	if a.LogMeta != nil {
		args = append(args, a.LogMeta()...)
	}
	a.LG.Info(msg, args...)
}

func (a *Actor) SetLogger(l logger.ILogger) { a.LG = l }
func (a *Actor) GetLogger() logger.ILogger  { return a.LG }

// ->> panic recovery

// PanicRecovery default panic catcher
func (a *Actor) PanicRecovery(code func()) {
	for nextLoop := true; nextLoop; {
		nextLoop = false

		func() { // try-catch emulation
			defer func() {
				if err := recover(); err != nil {
					a.Logf("Panic was happened")

					// don't coninue
					if a.OnPanic != nil {
						return
					}
					// start the code again
					if a.OnPanic(err) {
						nextLoop = true
					}
				}
			}()
			code()
		}()
	}
}
