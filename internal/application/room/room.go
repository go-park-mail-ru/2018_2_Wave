package room

import (
	lg "Wave/internal/logger"
	"fmt"
	"sync"
	"time"
)

type Route func(IUser, IInMessage) IRouteResponse

// Room - default IRoom
type Room struct {
	ID      RoomToken        // room ID
	Type    RoomType         // room type
	Ticker  *time.Ticker     // room tick
	Routes  map[string]Route // signal -> handler
	Users   map[UserID]IUser // room users
	LG      *lg.Logger       // logger
	Manager IRoomManager     // just a boss

	OnTick        func(time.Duration)
	OnUserAdded   func(IUser)
	OnUserRemoved func(IUser)

	broadcast       chan IRouteResponse
	CancelRoom      chan interface{}
	CancelBroadcast chan interface{}
	task            chan func()

	Step    time.Duration
	counter *Counter
}

func NewRoom(id RoomToken, tp RoomType, step time.Duration) *Room {
	r := &Room{
		ID:              id,
		Type:            tp,
		Ticker:          time.NewTicker(step),
		Routes:          map[string]Route{},
		Users:           map[UserID]IUser{},
		broadcast:       make(chan IRouteResponse, 150),
		CancelRoom:      make(chan interface{}, 1),
		CancelBroadcast: make(chan interface{}, 1),
		task:            make(chan func(), 150),
		counter:         NewCounter(Counting),
		Step:            step,
	}
	return r
}

func (r *Room) GetID() RoomToken  { return r.ID }
func (r *Room) GetType() RoomType { return r.Type }

func (r *Room) Run() error {
	r.log("room started")
	go r.runBroadcast()
	for {
		select {
		case <-r.Ticker.C:
			if r.OnTick != nil {
				r.OnTick(r.Step)
			}
		case clb := <-r.task:
			clb()
		case <-r.CancelRoom:
			return nil
		}
	}
}

func (r *Room) Stop() error {
	r.log("room stoped")
	r.CancelBroadcast <- ""
	r.CancelRoom <- ""
	return nil
}

// must be runned in a new goroutine
func (r *Room) runBroadcast() {
	r.log("broadcast started")
	for {
		select {
		case rs := <-r.broadcast:
			for _, u := range r.Users {
				r.SendMessageTo(u, rs)
			}
			// case <-r.CancelBroadcast:
			// return
		}
	}
}

func (r *Room) AddUser(u IUser) (err error) {
	if u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; !ok {
		r.Users[u.GetID()] = u
		r.counter.Add(u.GetID())
		r.log("user added", u.GetID())
		if r.OnUserAdded != nil {
			r.Async(func() { r.OnUserAdded(u) })
		}
		return nil
	}
	return ErrorAlreadyExists
}

func (r *Room) RemoveUser(u IUser) (err error) {
	if u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; ok {
		delete(r.Users, u.GetID())
		r.counter.Delete(u.GetID())
		r.log("user removed", u.GetID())
		if r.OnUserRemoved != nil {
			r.Async(func() { r.OnUserRemoved(u) })
		}
		return nil
	}
	return ErrorNotExists
}

func (r *Room) OnDisconnected(u IUser) {
	r.RemoveUser(u)
}

func (r *Room) Task(t func()) {
	r.task <- t
}

func (r *Room) ApplyMessage(u IUser, im IInMessage) error {
	if im == nil || u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; !ok {
		return ErrorForbiden
	}
	defer func() { // global handler panic
		if err := recover(); err != nil {
			r.log("ws route panic", "who", r.GetID())
		}
	}()

	if route, ok := r.Routes[im.GetSignal()]; ok {
		if om := route(u, im); om != nil {
			return r.SendMessageTo(u, om)
		}
		return nil
	}
	return ErrorUnknownSignal
}

func (r *Room) SendMessageTo(u IUser, rs IRouteResponse) error {
	if u == nil || rs == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; !ok {
		return ErrorForbiden
	}
	return u.Consume(&OutMessage{
		RoomToken: r.GetID(),
		Status:    rs.GetStatus(),
		Payload:   rs.GetPayload(),
	})
}

func (r *Room) Broadcast(rs IRouteResponse) error {
	if rs == nil {
		return ErrorNil
	}
	r.broadcast <- rs
	return nil
}

func (r *Room) GetUserCounter(u IUser) (counter int64, err error) {
	return r.counter.GetUserCounter(u)
}

func (r *Room) GetTokenCounter(t UserID) (counter int64, err error) {
	return r.counter.GetTokenCounter(t)
}

func (r *Room) SetCounterType(CounterType NumerationType) {
	r.counter.UserCounterType = CounterType
}

func (r *Room) IsAbleToRemove(IUser) bool {
	return true
}

// ----------------|

func (r *Room) log(data ...interface{}) {
	data = append([]interface{}{
		"room_token", r.ID,
		"room_type", r.Type,
	}, data...)

	if r.LG != nil {
		r.LG.Sugar.Infof("room_message", data...)
	} else {
		fmt.Println(data...)
	}
}

func (r *Room) doTask(t func()) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	r.task <- func() {
		defer wg.Done()
		t()
	}
	wg.Wait()
}

func (r *Room) Async(t func()) {
	r.task <- t
}
