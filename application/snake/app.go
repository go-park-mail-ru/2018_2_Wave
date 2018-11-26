package snake

import (
	"Wave/application/manager"
	"Wave/application/room"
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
func New(id room.RoomID, step time.Duration, db interface{}) room.IRoom {
	s := &App{
		Room: room.NewRoom(id, RoomType, step),
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
	var (
		info = a.world.GetGameInfo()
		msg  = room.MessageTick.WithStruct(info)
	)
	a.Broadcast(msg)
}

// get information about map and current users
func (a *App) onGameInfo(u room.IUser, im room.IInMessage) room.IRouteResponse {
	return room.MessageOK.WithStruct(a.world.GetGameInfo())
}

// receive game action (control)
func (a *App) onGameAction(u room.IUser, im room.IInMessage) room.IRouteResponse {
	type Action struct {
		ActionName string
	}

	ac := &Action{}
	if im.ToStruct(ac) != nil {
		return room.MessageWrongFormat
	}

	switch ac.ActionName {
	case "move_left":
		return a.withSnake(u, func(s *snake) { s.movement = left })
	case "move_right":
		return a.withSnake(u, func(s *snake) { s.movement = right })
	case "move_up":
		return a.withSnake(u, func(s *snake) { s.movement = up })
	case "move_down":
		return a.withSnake(u, func(s *snake) { s.movement = down })
	default:
		return messageUnknownCommand
	}
}

// place the user into a game scene and allow him play
func (a *App) onGamePlay(u room.IUser, im room.IInMessage) room.IRouteResponse {
	if _, err := a.world.CreateSnake(u, 6); err != nil {
		return messageAlreadyPlays
	}
	return nil
}

// exit from the game
func (a *App) onGameExit(u room.IUser, im room.IInMessage) room.IRouteResponse {
	if err := a.world.DeleteSnake(u); err != nil {
		return messageNoSnake
	}
	return nil
}

// ----------------| helpers

var (
	messageNoSnake        = room.RouteResponse{Status: room.StatusError}.WithStruct("No snake")
	messageAlreadyPlays   = room.RouteResponse{Status: room.StatusError}.WithStruct("already plays")
	messageUnknownCommand = room.RouteResponse{Status: room.StatusError}.WithStruct("unknown command")
)

func (a *App) withSnake(u room.IUser, next func(s *snake)) room.IRouteResponse {
	if s, err := a.world.GetSnake(u); err == nil {
		next(s)
		return nil
	}
	return messageNoSnake
}
