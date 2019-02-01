package proto

import (
	"Wave/internal/logger"
	"fmt"
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
	Task(IActor, ActorTask) IPipe
	SetLogger(logger.ILogger)
	GetLogger() logger.ILogger

	Sync(others ...IActor) ISyncCall
	setSync(IActor, bool)
	getSync(IActor) bool
	getUID() uint64
}

// ISyncCall - ??
type ISyncCall interface {
	Call(func()) IPipe
}

// IPipe - ??
type IPipe interface {
	Then(func()) IPipe
}

// **********************************************|
// **************************| 	Actor

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
	return &syncCall{actors, a}
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
func (a *Actor) getSync(o IActor) bool { return a.syncMap[o] }
func (a *Actor) getUID() uint64        { return a.uid }

// ---------------------| task

// Task - Add a synchronous task
// @param o - task instigator
// @param t	- task function
// @retrun 	- task promice
func (a *Actor) Task(o IActor, task ActorTask) IPipe {
	p := newPipe(a)
	if task != nil {
		clb := func() {
			task()
			p.done()
		}
		if o != nil && o.getSync(a) {
			clb()
		} else {
			a.T <- clb
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

// PanicRecoveryAsync - async call with panic recovery
func (a *Actor) PanicRecoveryAsync(code func()) {
	go func() { a.PanicRecovery(code) }()
}

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

// **********************************************|
// **************************| 	sync Call

type syncCall struct {
	actors []IActor
	caller IActor
}

// Call - async function that sync the actors and execute the code
func (s *syncCall) Call(code func()) IPipe {
	var (
		a      = s.actors[0]
		others = s.actors[1:]
		chans  = map[IActor]chan interface{}{}
		p      = newPipe(s.caller)
	)
	a.Task(nil, func() {
		// sync other actors
		wg := sync.WaitGroup{}
		wg.Add(len(others))
		for _, o := range others {
			chans[o] = make(chan interface{})
			o.Task(a, func() {
				o.setSync(a, true)
				a.setSync(o, true)
				wg.Done()
				<-chans[o]
			})
		}
		wg.Wait()

		code()

		// unlock other actors
		for _, o := range others {
			o.setSync(a, false)
			a.setSync(o, false)
			close(chans[o])
		}
		p.done()
	})
	return p
}

// **********************************************|
// **************************| 	pipe

type pipe struct {
	callbacks []func()
	a         IActor
}

func newPipe(a IActor) *pipe {
	return &pipe{[]func(){}, a}
}

func (p *pipe) Then(code func()) IPipe {
	p.callbacks = append(p.callbacks, code)
	return p
}

func (p *pipe) done() {
	if len(p.callbacks) > 0 {
		callback := p.callbacks[0]
		p.callbacks = p.callbacks[1:]
		p.a.Task(nil, func() {
			callback()
			p.done()
		})
	}
}
