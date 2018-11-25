package chat

import (
	"Wave/server/chat/room"
	"Wave/server/chat/app"
	"time"
)

const RoomType = "chat"

type App struct {
	*room.Room
}

func New(id room.RoomID, step time.Duration, db interface{}) room.IRoom {
	s := &App{
		Room: room.NewRoom(id, step),
	}
	s.Routes["send"] = s.send
	return s
}

func init() {
	// register the room type
	if err := app.RegisterRoomType(RoomType, New); err != nil {
		panic(err)
	}
}

func (a *App) send(u room.IUser, im room.IInMessage) room.IRouteResponse {
	room.MessageOK.WithStruct(im.GetPayload())
	return room.MessageOK
}
