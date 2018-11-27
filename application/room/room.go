package room

import (
	lg "Wave/utiles/logger"
	"fmt"
	"time"
)

type Route func(IUser, IInMessage) IRouteResponse

// Room - default IRoom
// - Chat
type Room struct {
	ID     RoomID           // room ID
	Type   RoomType         // room type
	Ticker *time.Ticker     // room tick
	Routes map[string]Route // signal -> handler
	Users  map[UserID]IUser // room users
	LG     *lg.Logger

	OnTick        func(time.Duration)
	OnUserAdded   func(IUser)
	OnUserRemoved func(IUser)

	broadcast       chan IRouteResponse
	CancelRoom      chan interface{}
	CancelBroadcast chan interface{}

	Step time.Duration
}

func NewRoom(id RoomID, tp RoomType, step time.Duration) *Room {
	r := &Room{
		ID:              id,
		Type:            tp,
		Ticker:          time.NewTicker(step),
		Routes:          map[string]Route{},
		Users:           map[UserID]IUser{},
		broadcast:       make(chan IRouteResponse, 150),
		CancelRoom:      make(chan interface{}, 1),
		CancelBroadcast: make(chan interface{}, 1),
		Step:            step,
	}
	return r
}

func (r *Room) GetID() RoomID     { return r.ID }
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
		case <-r.CancelBroadcast:
			return
		}
	}
}

func (r *Room) AddUser(u IUser) error {
	if u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; !ok {
		r.Users[u.GetID()] = u
		if r.OnUserAdded != nil {
			r.log("user added", u.GetID())
			r.OnUserAdded(u)
		}
		return nil
	}
	return ErrorAlreadyExists
}

func (r *Room) RemoveUser(u IUser) error {
	if u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; ok {
		delete(r.Users, u.GetID())
		if r.OnUserRemoved != nil {
			r.log("user removed", u.GetID())
			r.OnUserRemoved(u)
		}
		return nil
	}
	return ErrorNotExists
}

func (r *Room) OnDisconnected(u IUser) {
	r.RemoveUser(u)
}

func (r *Room) ApplyMessage(u IUser, im IInMessage) error {
	if im == nil || u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetID()]; !ok {
		return ErrorForbiden
	}
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
		RoomID:  r.GetID(),
		Status:  rs.GetStatus(),
		Payload: rs.GetPayload(),
	})
}

func (r *Room) Broadcast(rs IRouteResponse) error {
	if rs == nil {
		return ErrorNil
	}
	r.broadcast <- rs
	return nil
}

// ----------------|

func (r *Room) log(data ...interface{}) {
	data = append([]interface{}{
		"room_id", r.ID,
		"room_type", r.Type,
	}, data...)

	if r.LG != nil {
		r.LG.Sugar.Infof("room_message", data...)
	} else {
		fmt.Println(data...)
	}
}
