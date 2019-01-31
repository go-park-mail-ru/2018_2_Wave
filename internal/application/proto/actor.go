package proto

import (
	"Wave/internal/logger"
	"fmt"
	"github.com/fanliao/go-promise"
	"runtime/debug"
	"sort"
	"sync"
	"sync/atomic"
)

// ----------------| ActorTask

// ActorTask - message an actor
type ActorTask func()

// ----------------| IActor

// IActor - Actor inteface
type IActor interface {
	Task(IActor, ActorTask) *promise.Promise
	SetLogger(logger.ILogger)
	GetLogger() logger.ILogger

	Sync(others ...IActor) ISyncCall
	setSync(IActor, bool)
	getSync(IActor) bool
	getUID() uint64
}

// ISyncCall - ??
type ISyncCall interface {
	Call(func()) *promise.Promise
}

// ----------------| Actor

var actorUIDCounter uint64

// Actor - IActor interface default realisation
type Actor struct {
	T       chan ActorTask
	LG      logger.ILogger
	LogMeta func() []interface{}
	OnPanic func(interface{}) bool // true - continue the code
	syncMap map[IActor]bool        // sync map
	syncMu  sync.Mutex             // sync map mutex
	uid     uint64
}

type syncCall struct {
	actors []IActor
}

// MakeActor - constructor
func MakeActor(bufferSize int) Actor {
	return Actor{
		T:       make(chan ActorTask, bufferSize),
		syncMap: make(map[IActor]bool),
		uid:     atomic.AddUint64(&actorUIDCounter, 1),
	}
}

// ---------------------| sync

// Sync - crate a struct with an async call
// @see syncCall.Call()
func (a *Actor) Sync(others ...IActor) ISyncCall {
	// to prevent deadlocks we need to sort the actors
	actors := append(others, a)
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].getUID() < actors[j].getUID()
	})
	return &syncCall{actors}
}

// Call - async function that sync the actors and execute the code
func (s *syncCall) Call(code func()) *promise.Promise {
	var (
		a      = s.actors[0]
		others = s.actors[1:]
		chans  = map[IActor]chan interface{}{}
		p      = promise.NewPromise()
	)
	a.Task(nil, func() {
		// sync other actors
		wg := sync.WaitGroup{}
		wg.Add(len(others))
		for _, o := range others {
			// create a stop chanel
			chans[o] = make(chan interface{})
			// lock the actor
			o.Task(a, func() {
				// set sync flags
				o.setSync(a, true)
				a.setSync(o, true)
				// set as sync
				wg.Done()
				// lock the actor
				<-chans[o]
			})
		}
		wg.Wait()
		
		// run the code
		code()

		// unlock other actors
		for _, o := range others {
			// unset sync flags
			o.setSync(a, true)
			a.setSync(o, true)
			// unlock the actor
			close(chans[o])
		}
		p.Resolve(nil)
	})
	return p
}

func (a *Actor) setSync(o IActor, state bool) {
	a.syncMu.Lock()
	defer a.syncMu.Unlock()
	if state == true {
		a.syncMap[o] = true
	} else {
		delete(a.syncMap, o)
	}
}

func (a *Actor) getSync(o IActor) bool {
	return a.syncMap[o]
}

func (a *Actor) getUID() uint64 {
	return a.uid
}

// ---------------------| task

// Task - Add a synchronous task
// @param o - task instigator
// @param t	- task function
// @retrun 	- task promice
func (a *Actor) Task(o IActor, t ActorTask) *promise.Promise {
	// nil task
	if t == nil {
		p := promise.NewPromise()
		p.Reject(ErrorNil)
		return p
	}

	// run the task or store it
	p := promise.NewPromise()
	if o != nil && o.getSync(a) {
		// run the task immediately
		t()
		p.Resolve(nil)
	} else {
		// add the task to the task queue
		a.T <- func() {
			t()
			p.Resolve(nil)
		}
	}
	return p
}

// ---------------------| log

// Logf - format log
// @param pattern 	- pattern string
// @param args		- pattern args
func (a *Actor) Logf(pattern string, args ...interface{}) {
	if a.LG == nil {
		return
	}
	if a.LogMeta != nil {
		a.LG.Info(fmt.Sprintf(pattern, args...), a.LogMeta()...)
	} else {
		a.LG.Info(fmt.Sprintf(pattern, args...))
	}
}

// Log - just a log
// @param msg 	- message string
// @param args	- key-value args
func (a *Actor) Log(msg string, args ...interface{}) {
	if a.LG == nil {
		return
	}
	if a.LogMeta != nil {
		args = append(args, a.LogMeta()...)
	}
	a.LG.Info(msg, args...)
}

// Logn - Log with msg length cutting down to 60 runes
// @param msg 	- message string
// @param args	- key-value args
func (a *Actor) Logn(msg string, args ...interface{}) {
	if len(msg) > 60 {
		end := "..."
		a.Log(msg[:60-len(end)]+end, args...)
	} else {
		a.Log(msg, args...)
	}
}

func (a *Actor) SetLogger(l logger.ILogger) { a.LG = l }
func (a *Actor) GetLogger() logger.ILogger  { return a.LG }

// ---------------------| panic recovery

// PanicRecovery - panic recovery function
// @see OnPanic()
func (a *Actor) PanicRecovery(code func()) {
	for nextLoop := true; nextLoop; {
		nextLoop = false

		func() { // try-catch emulation
			defer a.panicRecoveryEntery(&nextLoop)
			code()
		}()
	}
}

func (a *Actor) panicRecoveryEntery(bNext *bool) {
	if err := recover(); err != nil {
		a.Log("Panic",
			"error", err)
		debug.PrintStack()

		// don't coninue
		if a.OnPanic == nil {
			return
		}
		// start the code again
		if a.OnPanic(err) {
			*bNext = true
		}
	}
}
