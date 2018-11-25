package room

import (
	"io"

	lg "Wave/utiles/logger"

	"github.com/gorilla/websocket"
)

type User struct {
	ID    UserID
	Rooms map[RoomID]IRoom
	Conn  *websocket.Conn
	LG    *lg.Logger

	bClosed bool
}

func NewUser(ID UserID, Conn *websocket.Conn) *User {
	return &User{
		ID:    ID,
		Conn:  Conn,
		Rooms: map[RoomID]IRoom{},
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

	// send current user_id
	u.Conn.WriteJSON(u.GetID())

	for { // stops when connection closes
		m := &InMessage{}

		// read a message
		err := u.Conn.ReadJSON(m)
		if err == io.ErrClosedPipe {
			u.removeFromAllRooms()
			if u.bClosed {
				return nil
			}
			return ErrorConnectionClosed
		}
		if err != nil {
			u.Consume(&OutMessage{
				RoomID:  m.GetRoomID(),
				Status:  StatusError,
				Payload: []byte("Wrong message"),
			})
			continue
		}

		// log input
		if u.LG != nil {
			u.LG.Sugar.Infof("in_message: %v", m)
		}

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
	u.bClosed = true
	u.Conn.Close()
	return nil
}

func (u *User) Consume(m IOutMessage) error {
	if m == nil {
		return ErrorNil
	}

	// log input
	if u.LG != nil {
		u.LG.Sugar.Infof("out_message: %v", m)
	}

	err := u.Conn.WriteJSON(m)
	if websocket.IsCloseError(err) { // is that correct?
		return ErrorConnectionClosed
	}
	if err != nil {
		return ErrorWrongInputFormat
	}
	return nil
}

// ----------------| internal function

func (u *User) removeFromAllRooms() error {
	for _, r := range u.Rooms {
		err := u.RemoveFromRoom(r)
		if err != nil {
			return err
		}
	}
	return nil
}
