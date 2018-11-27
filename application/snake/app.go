package snake

import (
	"Wave/application/manager"
	"Wave/application/room"
	"Wave/application/snake/core"
	"time"
)

// RoomType - snake type literal
const RoomType room.RoomType = "snake_game"

// App - snake game room
type App struct {
	*room.Room // base room
	game      *game
}

// New snake app
func New(id room.RoomID, step time.Duration, db interface{}) room.IRoom {
	s := &App{
		Room: room.NewRoom(id, RoomType, step),
		game: newGame(core.Vec2i{
			X: 130,
			Y: 10,
		}),
	}
	s.OnTick = s.onTick
	s.OnUserRemoved = s.onUserRemoved
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
	a.game.Tick(dt)
	var (
		info = a.game.GetGameInfo()
		msg  = room.MessageTick.WithStruct(info)
	)
	a.Broadcast(msg)
}

func (a *App) onUserRemoved(u room.IUser) {
	a.game.DeleteSnake(u)
}

// get information about map and current users
func (a *App) onGameInfo(u room.IUser, im room.IInMessage) room.IRouteResponse {
	return room.MessageOK.WithStruct(a.game.GetGameInfo())
}

// receive game action (control)
func (a *App) onGameAction(u room.IUser, im room.IInMessage) room.IRouteResponse {
	type Action struct {
		ActionName string `json:"action"`
	}

	ac := &Action{}
	if im.ToStruct(ac) != nil {
		return room.MessageWrongFormat
	}

	switch ac.ActionName {
	case "move_left":
		return a.withSnake(u, func(s *snake) { s.movement = core.Left })
	case "move_right":
		return a.withSnake(u, func(s *snake) { s.movement = core.Right })
	case "move_up":
		return a.withSnake(u, func(s *snake) { s.movement = core.Up })
	case "move_down":
		return a.withSnake(u, func(s *snake) { s.movement = core.Down })
	default:
		return messageUnknownCommand
	}
}

// place the user into a game scene and allow him play
func (a *App) onGamePlay(u room.IUser, im room.IInMessage) room.IRouteResponse {
	if _, err := a.game.CreateSnake(u, 6); err != nil {
		return messageAlreadyPlays
	}
	return nil
}

// exit from the game
func (a *App) onGameExit(u room.IUser, im room.IInMessage) room.IRouteResponse {
	if err := a.game.DeleteSnake(u); err != nil {
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
	if s, err := a.game.GetSnake(u); err == nil {
		next(s)
		return nil
	}
	return messageNoSnake
}
