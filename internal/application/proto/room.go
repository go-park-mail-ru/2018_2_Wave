package proto

import (
	"Wave/internal/logger"
	"time"
)

// ----------------| RoomToken && RoomType && Route

// RoomToken - stringified room id
type RoomToken string

// RoomType - room type name
type RoomType string

// Route - ws signal handler
type Route func(IUser, IInMessage)

type inLetter struct {
	u IUser
	m IInMessage
}

// ----------------| IRoom

// IRoom interface
type IRoom interface {
	GetToken() RoomToken             // get room token
	GetType() RoomType               // get room type
	GetManager() IManager            // get room manager
	Receive(IUser, IInMessage) error // receive the message

	addUser(IUser) error    // add the user to the room
	removeUser(IUser) error // remove the user from the room
	onDisconnected(IUser)   //

	Start() // start the room. Locks the thread
	Stop()  // stop the room

	GetUserSerial(u IUser) (serial int64, err error)      // get user local id
	GetTokenSerial(t UserToken) (serial int64, err error) // get user local id
	SetCounterType(CounterType NumerationType)            // set id numeration mode

	IsAbleToRemove(u IUser) bool // wether the user can remove the room
}

// RoomFactory - IRoom factort
type RoomFactory func(_ RoomToken, _ RoomType, _ IManager, db interface{}, step time.Duration) IRoom

// ---------------| Room

// Room base class
type Room struct {
	Actor // base class

	Users   map[UserToken]IUser
	Routes  map[string]Route
	LG      logger.ILogger
	manager IManager

	counter Counter
	ticker  time.Ticker
	step    time.Duration
	token   RoomToken
	rtype   RoomType

	inMessages chan inLetter
	broadcast  chan IResponse
	cancel     chan interface{}

	OnTick           func(time.Duration) // tick function
	OnUserAdded      func(IUser)         // user was added
	OnUserRemove     func(IUser)         // user will be removed
	OnUserDisconnect func(IUser)         // user was disconnected
}

// NewRoom - constructor
func NewRoom(token RoomToken, rtype RoomType, manager IManager, step time.Duration) *Room {
	r := &Room{
		Actor:   MakeActor(100),
		ticker:  *time.NewTicker(step),
		counter: MakeCounter(FillGaps),
		manager: manager,
		token:   token,
		step:    step,

		Routes:     make(map[string]Route),
		Users:      make(map[UserToken]IUser),
		inMessages: make(chan inLetter, 100),
		broadcast:  make(chan IResponse, 100),
		cancel:     make(chan interface{}, 6),
	}
	r.LogMeta = func() []interface{} {
		return []interface{}{"room", r.GetToken(), "type", r.GetType()}
	}
	return r
}

// ------| << IRoom

func (r *Room) GetToken() RoomToken  { return r.token }   // safe
func (r *Room) GetType() RoomType    { return r.rtype }   // safe
func (r *Room) GetManager() IManager { return r.manager } // safe

// Receive the message - safe
func (r *Room) Receive(u IUser, m IInMessage) error {
	if u == nil || m == nil {
		return ErrorNil
	}
	r.inMessages <- inLetter{u, m}
	return nil
}

func (r *Room) addUser(u IUser) error {
	if u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetToken()]; ok {
		return ErrorAlreadyExists
	}

	r.Users[u.GetToken()] = u
	r.onUserAdded(u)
	return nil
}

func (r *Room) removeUser(u IUser) error {
	if u == nil {
		return ErrorNil
	}
	if _, ok := r.Users[u.GetToken()]; !ok {
		return ErrorNotFound
	}

	r.onUserRemove(u)
	delete(r.Users, u.GetToken())
	return nil
}

func (r *Room) onDisconnected(u IUser) {
	if u == nil && r.OnUserDisconnect != nil {
		r.OnUserDisconnect(u)
	}
}

// Start serving
func (r *Room) Start() {
	r.Logf("room started")
	for { // main room loop
		select {
		case <-r.ticker.C:
			r.PanicRecovery(r.doTick)
		case t := <-r.Actor.T:
			r.PanicRecovery(t)
		case m := <-r.broadcast:
			r.PanicRecovery(func() {
				r.doBroadcast(m)
			})
		case p := <-r.inMessages:
			r.PanicRecovery(func() {
				r.doReceive(p)
			})
		case <-r.cancel:
			r.PanicRecovery(r.doCancel)
			r.Logf("room stopped")
			return
		}
	}
}

// Stop the room - safe
func (r *Room) Stop() {
	r.cancel <- ""
}

func (r *Room) GetUserSerial(u IUser) (serial int64, err error)      { return r.counter.GetUserID(u) }
func (r *Room) GetTokenSerial(t UserToken) (serial int64, err error) { return r.counter.GetTokenID(t) }
func (r *Room) SetCounterType(CounterType NumerationType)            { r.counter.UserCounterType = CounterType }
func (r *Room) IsAbleToRemove(u IUser) bool                          { return true }

// ------| Messages - message interface

// SendTo - send the message to the user
func (r *Room) SendTo(u IUser, m IResponse) error {
	if u == nil || m == nil {
		return ErrorNil
	}
	return u.Send(&OutMessage{
		RoomToken: r.GetToken(),
		Status:    m.GetStatus(),
		Payload:   m.GetPayload(),
	})
}

// Broadcast the message
func (r *Room) Broadcast(m IResponse) error {
	if m == nil {
		return ErrorNil
	}
	r.broadcast <- m
	return nil
}

// ------| workers - main cycle workers

func (r *Room) doBroadcast(m IResponse) {
	for _, u := range r.Users {
		r.SendTo(u, m)
	}
}

func (r *Room) doTick() {
	if r.OnTick != nil {
		r.OnTick(r.step)
	}
}

func (r *Room) doReceive(p inLetter) {
	if rt, ok := r.Routes[p.m.GetSignal()]; ok {
		rt(p.u, p.m)
	}
}

func (r *Room) doCancel() {
	if r.manager != nil {
		r.manager.RemoveLobby(r.GetToken(), nil)
	}
}

// ------| events - internal events

func (r *Room) onUserAdded(u IUser) {
	r.counter.Register(u.GetToken())
	if r.OnUserAdded != nil {
		r.OnUserAdded(u)
	}
}

func (r *Room) onUserRemove(u IUser) {
	r.counter.Unregister(u.GetToken())
	if r.OnUserRemove != nil {
		r.OnUserRemove(u)
	}
}
