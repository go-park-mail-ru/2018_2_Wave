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

	return s
}
