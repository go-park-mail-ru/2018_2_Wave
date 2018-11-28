package room

import (
	lg "Wave/internal/logger"

	"github.com/gorilla/websocket"
)

type User struct {
	ID    UserID
	Rooms map[RoomID]IRoom
	Conn  *websocket.Conn
	LG    *lg.Logger

	output  chan IOutMessage
	cancel  chan interface{}
	bClosed bool
}

func NewUser(ID UserID, Conn *websocket.Conn) *User {
	return &User{
		ID:     ID,
		Conn:   Conn,
		cancel: make(chan interface{}, 1),
		output: make(chan IOutMessage, 1000),
		Rooms:  map[RoomID]IRoom{},
	}
}

// ----------------| IUser interface

func (u *User) GetID() UserID {
	return u.ID
}

func (u *User) AddToRoom(r IRoom) error {
	if r == nil {
		return ErrorNil
	}
	if err := r.AddUser(u); err != nil {
		return err
	}
	u.Rooms[r.GetID()] = r
	return nil
}

func (u *User) RemoveFromRoom(r IRoom) error {
	if r == nil {
		return ErrorNil
	}
	if err := r.RemoveUser(u); err != nil {
		return err
	}
	delete(u.Rooms, r.GetID())
	return nil
}

func (u *User) Listen() error {
	defer func() {
		if err := recover(); err != nil {
			u.StopListening()
		}
	}()

	u.LG.Sugar.Infof("User started: id=", u.GetID())

	go u.sendWorker()

	// send current user_id
	u.Conn.WriteJSON(u.GetID())

	for { // stops when connection closes
		m := &InMessage{}

		// read a message
		if err := u.Conn.ReadJSON(m); err != nil {
			if websocket.IsCloseError(err, wsCloseErrors...) {
				u.onDisconnected()
				u.stop()
				if u.bClosed {
					return nil
				}
				return ErrorConnectionClosed
			}
			u.Consume(&OutMessage{
				RoomID:  m.GetRoomID(),
				Status:  StatusError,
				Payload: []byte("Wrong message"),
			})
			continue
		}

		// log input
		// if u.LG != nil {
		// 	data, _ := json.Marshal(m)
		// 	u.LG.Sugar.Infof("in_message: %v", string(data))
		// }

		// apply the message to a room
		if r, ok := u.Rooms[m.GetRoomID()]; ok {
			r.ApplyMessage(u, m)
		} else {
			u.Consume(&OutMessage{
				RoomID:  m.GetRoomID(),
				Status:  StatusError,
				Payload: []byte("Unknown room:" + m.GetRoomID()),
			})
			continue
		}
	}
}

func (u *User) StopListening() error {
	if !u.bClosed {
		u.removeFromAllRooms() // do we need this?
		u.stop()
	}
	return nil
}

func (u *User) Consume(m IOutMessage) error {
	if m == nil {
		return ErrorNil
	}
	u.output <- m
	return nil
}

// ----------------| internal function

func (u *User) sendWorker() {
	for {
		select {
		case m := <-u.output:
			// log input
			// if u.LG != nil {
			// 	data, _ := json.Marshal(m)
			// 	u.LG.Sugar.Infof("out_message: %v", string(data))
			// }

			if err := u.Conn.WriteJSON(m); err != nil {
				u.StopListening()
				u.stop()
			}
		case <-u.cancel:
			return
		}
	}
}

func (u *User) removeFromAllRooms() error {
	for _, r := range u.Rooms {
		err := u.RemoveFromRoom(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) onDisconnected() {
	for _, r := range u.Rooms {
		r.OnDisconnected(u)
	}
}

func (u *User) stop() {
	u.LG.Sugar.Infof("ws closed, uid: %s", u.GetID())
	u.bClosed = true
	u.Conn.Close()
}
