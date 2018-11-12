package room

import "time"

type Command struct {
	IOutMessage
	Scope string
}

// Room - default IRoom
// 	- Chat
type Room struct {
	ID     RoomID                      // room ID
	Ticker *time.Ticker                // room tick
	Roures map[string]func(IInMessage) // signal -> handler
	Users  map[UserID]IUser            // room users

	OnTick        func(time.Duration)
	OnUserAdded   func(IUser)
	OnUserRemoved func(IUser)

	Broadcast       chan Command // broadcast message to it's scope
	CancelRoom      chan interface{}
	CancelBroadcast chan interface{}

	step time.Duration
}

func NewRoom(id RoomID, step time.Duration) *Room {
	r := &Room{
		ID:              id,
		Ticker:          time.NewTicker(step),
		Roures:          map[string]func(IInMessage){},
		Users:           map[UserID]IUser{},
		Broadcast:       make(chan Command, 150),
		CancelRoom:      make(chan interface{}, 1),
		CancelBroadcast: make(chan interface{}, 1),
		step:            step,
	}
	return r
}

func (r *Room) GetID() RoomID {
	return r.ID
}

func (r *Room) Run() error {
	go r.runBroadcast()
	for { // infinity cycle
		select {
		case <-r.Ticker.C:
			if r.OnTick != nil {
				r.OnTick(r.step)
			}
		case <-r.CancelRoom:
			return nil
		}
	}
}

func (r *Room) Stop() error {
	r.CancelBroadcast <- ""
	r.CancelRoom <- ""
	return nil
}

// must be runned in a new goroutine
func (r *Room) runBroadcast() {
	for { // infinity cycle
		select {
		case m := <-r.Broadcast:
			for _, p := range r.Users {
				p.Send(m)
			}
		case <-r.CancelBroadcast:
			return
		}
	}
}

func (r *Room) AddUser(usr IUser) error {
	if usr == nil {
		return ErrorNil
	}
	if _, ok := r.Users[usr.GetID()]; !ok {
		r.Users[usr.GetID()] = usr
		if r.OnUserAdded != nil {
			r.OnUserAdded(usr)
		}
		return nil
	}
	return ErrorAlreadyExists
}

func (r *Room) RemoveUser(usr IUser) error {
	if usr == nil {
		return ErrorNil
	}
	if _, ok := r.Users[usr.GetID()]; ok {
		delete(r.Users, usr.GetID())
		if r.OnUserRemoved != nil {
			r.OnUserRemoved(usr)
		}
		return nil
	}
	return ErrorNotExists
}

func (r *Room) SendMessage(im IInMessage) error {
	if im == nil {
		return ErrorNil
	}
	if route, ok := r.Roures[im.GetSignal()]; ok {
		route(im)
		return nil
	}
	return ErrorUnknownSignal
}
