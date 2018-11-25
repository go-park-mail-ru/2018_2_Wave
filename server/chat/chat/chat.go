package chat

import (
	"Wave/server/chat/room"
	"time"
)

type App struct {
	*room.Room
}

func New(id room.RoomID, step time.Duration) room.IRoom {
	s := &App{
		Room: room.NewRoom(id, step),
	}
	s.Routes["send"] = s.send
	return s
}

func (a *App) send(u room.IUser, im room.IInMessage) room.IRouteResponse {
	room.MessageOK.WithStruct(im.GetPayload())
	return room.MessageOK
}
