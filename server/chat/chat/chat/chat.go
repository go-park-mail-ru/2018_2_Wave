package chat

import (
	"Wave/server/room"
	"time"
)

type App struct {
	*room.Room
}

// New snake app
func New(id room.RoomID, step time.Duration) room.IRoom {
	s := &App{
		Room: room.NewRoom(id, step),
		world: newWorld(sceneSize{
			X: 900,
			Y: 900,
		}),
	}
	s.OnTick = s.onTick
	s.Routes["game_info"] = s.onGameInfo
	return s
}
