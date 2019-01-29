package snake

import (
	"Wave/internal/application/proto"
	"Wave/internal/application/snake/core"
	"time"
)

//go:generate easyjson .

type App struct {
	*proto.Room // base class
	game       *game
}

const RoomType proto.RoomType = "snake"

// ----------------|

// New snake app
func New(id proto.RoomToken, step time.Duration, m proto.IManager, db interface{}) proto.IRoom {
	s := &App{
		Room: proto.NewRoom(id, RoomType, m, step),
		game: newGame(core.Vec2i{
			X: 60,
			Y: 40,
		}),
	}
	s.SetCounterType(proto.FillGaps)
	s.OnTick = s.onTick
	s.OnUserRemoved = s.onUserRemoved
	s.game.OnSnakeDead = s.onSnakeDead
	s.Routes["game_action"] = s.onGameAction
	s.Routes["game_play"] = s.onGamePlay
	s.Routes["game_exit"] = s.onGameExit
	return s
}

// ----------------| handlers

// <- STATUS_TICK
func (a *App) onTick(dt time.Duration) {
	a.game.Tick(dt)
	info := a.game.GetGameInfo()
	for i, s := range info.Snakes {
		serial, _ := a.GetTokenCounter(s.UserToken)
		info.Snakes[i].Serial = serial
	}
	a.Broadcast(proto.MessageTick.WithStruct(info))
}

func (a *App) onUserRemoved(u proto.IUser) {
	a.game.DeleteSnake(u)
}

// -> game_action
func (a *App) onGameAction(u proto.IUser, im proto.IInMessage) {
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
}

// -> game_play
func (a *App) onGamePlay(u proto.IUser, im proto.IInMessage) {
	a.game.CreateSnake(u, 3)
}

// -> game_exit
func (a *App) onGameExit(u proto.IUser, im proto.IInMessage) {
	a.game.DeleteSnake(u)

	if len(a.game.user2snake) == 0 {
		a.exit()
	}
}

// <- STATUS_DEAD | win
func (a *App) onSnakeDead(u proto.IUser) {
	serial, _ := a.GetUserCounter(u)
	a.Broadcast(messageDead.WithStruct(&playerPayload{
		UserName:   u.GetName(),
		UserToken:  u.GetID(),
		UserSerial: serial,
	}))

	if len(a.game.user2snake) <= 1 {
		if len(a.game.user2snake) == 1 {
			var w proto.IUser
			for w = range a.game.user2snake {
			}

			serial, _ := a.GetUserCounter(w)
			a.Broadcast(messageWin.WithStruct(&playerPayload{
				UserName:   w.GetName(),
				UserToken:  w.GetID(),
				UserSerial: serial,
			}))
		}
		a.exit()
	}
}

func (a *App) exit() {
	for _, u := range a.Users {
		u.Task(func() { u.RemoveFromRoom(a) })
	}
	if a.Manager != nil {
		a.Manager.RemoveLobby(a.GetID(), nil)
	}
}

// ----------------| helpers

// easyjson:json
type gameAction struct {
	ActionName string `json:"action"`
}

type playerPayload struct {
	UserName   string      `json:"user_name"`
	UserToken  proto.UserToken `json:"user_token"`
	UserSerial int64       `json:"user_serial"`
}

var (
	messageWin            = proto.Response{Status: "win"}.WithStruct("")
	messageDead           = proto.Response{Status: "STATUS_DEAD"}.WithStruct("")
	messageNoSnake        = proto.Response{Status: "STATUS_ERROR"}.WithReason("No snake")
	messageAlreadyPlays   = proto.Response{Status: "STATUS_ERROR"}.WithReason("already plays")
	messageUnknownCommand = proto.Response{Status: "STATUS_ERROR"}.WithReason("unknown command")
)

func (a *App) withSnake(u proto.IUser, next func(s *snake)) {
	if s, err := a.game.GetSnake(u); err == nil {
		next(s)
		return nil
	}
	return messageNoSnake
}
