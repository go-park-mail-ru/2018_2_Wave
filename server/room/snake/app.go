package snake

import (
	"Wave/server/room"
	"Wave/server/room/app"
	"time"
)

// RoomType - snake type literal
const RoomType room.RoomType = "snake_game"

// App - snake game room
type App struct {
	*room.Room        // base room
	world      *world // game world
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
	s.Routes["game_action"] = s.onGameAction
	s.Routes["game_play"] = s.onGamePlay
	s.Routes["game_exit"] = s.onGameExit
	return s
}

func init() {
	// register the room type
	if err := app.RegisterRoomType(RoomType, New); err != nil {
		panic(err)
	}
}

// ----------------| handlers

func (a *App) onTick(dt time.Duration) {
	a.world.Tick(dt)

	// TODO: broadcast tick message
}

// get information about map and current users
func (a *App) onGameInfo(u room.IUser, im room.IInMessage) room.IRouteResponce {
	return room.RouteResponce{
		Status: room.StatusOK,
	}.WithStruct(a.world.GetGameInfo())
}

// receive game action (control)
func (a *App) onGameAction(u room.IUser, im room.IInMessage) room.IRouteResponce {
	type Action struct {
		ActionName string
	}

	ac := &Action{}
	if im.ToStruct(ac) != nil {
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("Incorrect data")
	}

	switch ac.ActionName {
	case "move_left":
		return a.withSnake(u, func(s *snake) room.IRouteResponce {
			s.movement = left
			return nil
		})
	case "move_right":
		return a.withSnake(u, func(s *snake) room.IRouteResponce {
			s.movement = right
			return nil
		})
	case "move_up":
		return a.withSnake(u, func(s *snake) room.IRouteResponce {
			s.movement = up
			return nil
		})
	case "move_down":
		return a.withSnake(u, func(s *snake) room.IRouteResponce {
			s.movement = down
			return nil
		})
	default:
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("Unknow command")
	}
}

// place the user into a game scene and allow him play
func (a *App) onGamePlay(u room.IUser, im room.IInMessage) room.IRouteResponce {
	if _, err := a.world.CreateSnake(u, 6); err != nil {
		return nil
	}
	return room.RouteResponce{
		Status: room.StatusError,
	}.WithStruct("already plays")
}

// exit from the game
func (a *App) onGameExit(u room.IUser, im room.IInMessage) room.IRouteResponce {
	if err := a.world.DeleteSnake(u); err != nil {
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("No snake")
	}
	return nil
}

// ----------------| helpers

func (a *App) withSnake(u room.IUser, next func(s *snake) room.IRouteResponce) room.IRouteResponce {
	if s, err := a.world.GetSnake(u); err == nil {
		return next(s)
	}
	return room.RouteResponce{
		Status: room.StatusError,
	}.WithStruct("No snake")
}
