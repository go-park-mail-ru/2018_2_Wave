package proto

import (
	"Wave/internal/logger"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// ----------------| UserToken

// UserToken - stringified user id
type UserToken string

var (
	wsCloseErrors = func() (res []int) {
		for i := 1000; i <= 1015; i++ {
			res = append(res, i)
		}
		return res
	}()
)

// ----------------| IUser

// IUser interface
type IUser interface {
	GetToken() UserToken    // get user token
	GetName() string        //
	Send(IOutMessage) error // send the message to the user
	EnterRoom(IRoom) error  // add the user into the room
	ExitRoom(IRoom) error   // remove the user from the room

	Start() // run the user
	Stop()  // stop the user

	IActor
}

// ----------------| User

// User - base user class
type User struct {
	Actor // base class

	Name  string
	token UserToken

	LG       logger.ILogger
	conn     *websocket.Conn
	Rooms    map[RoomToken]IRoom
	manager  IManager
	bStopped int32

	input  chan IInMessage
	output chan IOutMessage
	cancel chan interface{}
}

// NewUser - constructor
func NewUser(token UserToken, conn *websocket.Conn, manager IManager) (*User, error) {
	if conn == nil || manager == nil {
		return nil, ErrorNil
	}
	u := &User{
		Actor:   MakeActor(100),
		token:   token,
		conn:    conn,
		manager: manager,
		Rooms:   make(map[RoomToken]IRoom),
		input:   make(chan IInMessage, 100),
		output:  make(chan IOutMessage, 100),
		cancel:  make(chan interface{}, 3),
	}
	u.LogMeta = func() []interface{} {
		return []interface{}{"user", u.GetToken(), "name", u.GetName()}
	}
	u.Rooms[""] = manager
	return u, nil
}

// ------| << IUser

func (u *User) GetToken() UserToken { return u.token }
func (u *User) GetName() string     { return u.Name }

func (u *User) Send(m IOutMessage) error {
	if m == nil {
		return ErrorNil
	}
	u.output <- m
	return nil
}

func (u *User) EnterRoom(r IRoom) error {
	if r == nil {
		return ErrorNil
	}
	return r.Task(func() { r.addUser(u) })
}

func (u *User) ExitRoom(r IRoom) error {
	if r == nil {
		return ErrorNil
	}
	return r.Task(func() { r.removeUser(u) })
}

func (u *User) Start() {
	u.Logf("user started")
	go u.PanicRecovery(u.receiveWorker)
	go u.PanicRecovery(u.sendWorker)

	for {
		select {
		case t := <-u.T:
			u.PanicRecovery(t)
		case <-u.cancel:
			u.Logf("user stopped")
			return
		}
	}
}

func (u *User) Stop() {
	if atomic.AddInt32(&u.bStopped, 1) > 1 {
		u.cancel <- ""
		u.cancel <- ""
	}
}

// ------| workers

func (u *User) receiveWorker() {
	for {
		m := &InMessage{}

		// read a message
		if err := u.conn.ReadJSON(m); err != nil {
			u.Logf("wrong_message")
			u.onDisconnected()
			u.Stop()
			return
		}

		u.input <- m
	}
}

func (u *User) sendWorker() {
	for {
		select {
		case m := <-u.output:
			u.conn.WriteJSON(m)
		case <-u.cancel:
			return
		}
	}
}

// ----------------| internal

func (u *User) onDisconnected() {
	for _, r := range u.Rooms {
		r.Task(func() { r.onDisconnected(u) })
	}
}
