package snake

import (
	"Wave/internal/application/room"
	"Wave/internal/application/snake/core"
	"time"
)

//go:generate easyjson .

type App struct {
	*room.Room // base room
	game       *game
}

const RoomType room.RoomType = "snake"

// ----------------|

// New snake app
func New(id room.RoomToken, step time.Duration, db interface{}) room.IRoom {
	s := &App{
		Room: room.NewRoom(id, RoomType, step),
		game: newGame(core.Vec2i{
			X: 60,
			Y: 40,
		}),
	}
	s.SetCounterType(room.FillGaps)
	s.OnTick = s.onTick
	s.OnUserRemoved = s.onUserRemoved
	s.game.OnSnakeDead = s.onSnakeDead
	s.Routes["game_action"] = s.onGameAction
	s.Routes["game_play"] = s.onGamePlay
	s.Routes["game_exit"] = s.onGameExit
	return s
}

// ----------------| handlers

func (a *App) onTick(dt time.Duration) {
	a.game.Tick(dt)
	info := a.game.GetGameInfo()
	for i, s := range info.Snakes {
		serial, _ := a.GetTokenCounter(s.UserToken)
		info.Snakes[i].Serial = serial
	}
	a.Broadcast(room.MessageTick.WithStruct(info))
}

func (a *App) onUserRemoved(u room.IUser) {
	a.game.DeleteSnake(u)
}

// receive game action (control)
func (a *App) onGameAction(u room.IUser, im room.IInMessage) room.IRouteResponse {
	ac := &gameAction{}
	if im.ToStruct(ac) != nil {
		return nil
	}

	switch ac.ActionName {
	case "move_left":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Left) })
	case "move_right":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Right) })
	case "move_up":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Up) })
	case "move_down":
		a.withSnake(u, func(s *snake) { s.SetDirection(core.Down) })
	}
	return nil
}

// place the user into a game scene and allow him play
func (a *App) onGamePlay(u room.IUser, im room.IInMessage) room.IRouteResponse {
	a.game.CreateSnake(u, 6)
	return nil
}

// exit from the game
func (a *App) onGameExit(u room.IUser, im room.IInMessage) room.IRouteResponse {
	a.game.DeleteSnake(u)
	return nil
}

func (a *App) onSnakeDead(u room.IUser) {
	a.SendMessageTo(u, messageDead)
}

// ----------------| helpers

// easyjson:json
type gameAction struct {
	ActionName string `json:"action"`
}

var (
	messageDead           = room.RouteResponse{Status: "STATUS_DEAD"}.WithStruct("")
	messageNoSnake        = room.RouteResponse{Status: room.StatusError}.WithReason("No snake")
	messageAlreadyPlays   = room.RouteResponse{Status: room.StatusError}.WithReason("already plays")
	messageUnknownCommand = room.RouteResponse{Status: room.StatusError}.WithReason("unknown command")
)

func (a *App) withSnake(u room.IUser, next func(s *snake)) room.IRouteResponse {
	if s, err := a.game.GetSnake(u); err == nil {
		next(s)
		return nil
	}
	return messageNoSnake
}
