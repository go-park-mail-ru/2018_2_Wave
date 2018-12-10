package app

import (
	"Wave/application/room"
	"Wave/internal/metrics"
	"strconv"
	"sync"
	"time"
)

//go:generate easyjson .

// App - main service room
/* - creates and store other rooms
 * - contains ALL online users
 * - provides all service functions
 */
type App struct {
	*room.Room // the room super
	rooms      map[room.RoomToken]room.IRoom
	db         interface{}
	prof       *metrics.Profiler

	lastRoomID int64
	lastUserID int64
	idsMutex   sync.Mutex
}

const RoomType = "manager"

// ----------------|

// New applicarion room
func New(id room.RoomToken, step time.Duration, db interface{}, prof *metrics.Profiler) *App {
	a := &App{
		Room:  room.NewRoom(id, RoomType, step),
		rooms: map[room.RoomToken]room.IRoom{},
		prof:  prof,
		db:    db,
	}
	a.Routes["lobby_list"] = a.onGetLobbyList
	a.Routes["lobby_create"] = withRoomType(a.onLobbyCreate)
	a.Routes["lobby_delete"] = withRoomID(a.onLobbyDelete)
	a.Routes["add_to_room"] = withRoomID(a.onAddToRoom)
	a.Routes["remove_from_room"] = withRoomID(a.onRemoveFromRoom)
	return a
}

// ----------------| methods

// GetNextUserID returns next user id
func (a *App) GetNextUserID() room.UserID {
	a.idsMutex.Lock()
	defer a.idsMutex.Unlock()

	a.lastUserID++
	return room.UserID(strconv.FormatInt(a.lastUserID, 36))
}

// GetNextRoomID returns next room id
func (a *App) GetNextRoomID() room.RoomToken {
	a.idsMutex.Lock()
	defer a.idsMutex.Unlock()

	a.lastRoomID++
	return room.RoomToken(strconv.FormatInt(a.lastRoomID, 36))
}

// CreateLobby -
func (a *App) CreateLobby(room_type room.RoomType, room_token room.RoomToken) (room.IRoom, error) {
	if factory, ok := type2Factory[room_type]; ok {
		r := factory(room_token, a.Step, a.db)
		if r == nil {
			return nil, room.ErrorNil
		}
		a.rooms[room_token] = r
		go r.Run()

		// profiler
		if a.prof != nil {
			a.prof.ActiveRooms.Inc()
		}

		return r, nil
	}
	return nil, room.ErrorNotExists
}

// ----------------| handlers

func (a *App) onGetLobbyList(u room.IUser, im room.IInMessage) room.IRouteResponse {
	data := []roomInfoPayload{}
	for _, r := range a.rooms {
		data = append(data, roomInfoPayload{
			RoomToken: r.GetID(),
			RoomType:  r.GetType(),
		})
	}
	return room.MessageOK.WithStruct(data)
}

func (a *App) onLobbyCreate(u room.IUser, im room.IInMessage, cmd room.RoomType) room.IRouteResponse {
	r, err := a.CreateLobby(cmd, a.GetNextRoomID())
	if err != nil {
		return nil
	}
	return room.MessageOK.WithStruct(roomTokenPayload{
		RoomToken: r.GetID(),
	})
}

func (a *App) onLobbyDelete(u room.IUser, im room.IInMessage, cmd room.RoomToken) room.IRouteResponse {
	if r, ok := a.rooms[cmd]; ok {
		r.Stop()
		delete(a.rooms, cmd)

		// profiler
		if a.prof != nil {
			a.prof.ActiveRooms.Dec()
		}
	}
	return nil
}

func (a *App) onAddToRoom(u room.IUser, im room.IInMessage, cmd room.RoomToken) room.IRouteResponse {
	if r, ok := a.rooms[cmd]; ok {
		u.AddToRoom(r)
	}
	return nil
}

func (a *App) onRemoveFromRoom(u room.IUser, im room.IInMessage, cmd room.RoomToken) room.IRouteResponse {
	if r, ok := a.rooms[cmd]; ok {
		u.RemoveFromRoom(r)
	}
	return nil
}

// ----------------| helper functions

// easyjson:json
type roomTokenPayload struct {
	RoomToken room.RoomToken `json:"room_token"`
}

// easyjson:json
type roomTypePayload struct {
	RoomType room.RoomType `json:"room_type"`
}

// easyjson:json
type roomInfoPayload struct {
	RoomToken room.RoomToken `json:"room_token"`
	RoomType  room.RoomType  `json:"room_type"`
}

func withRoomID(next func(room.IUser, room.IInMessage, room.RoomToken) room.IRouteResponse) room.Route {
	return func(u room.IUser, im room.IInMessage) room.IRouteResponse {
		cmd := &roomTokenPayload{}
		if im.ToStruct(cmd) == nil {
			return next(u, im, cmd.RoomToken)
		}
		return room.MessageWrongFormat
	}
}

func withRoomType(next func(room.IUser, room.IInMessage, room.RoomType) room.IRouteResponse) room.Route {
	return func(u room.IUser, im room.IInMessage) room.IRouteResponse {
		cmd := &roomTypePayload{}
		if im.ToStruct(cmd) == nil {
			return next(u, im, cmd.RoomType)
		}
		return room.MessageWrongFormat
	}
}
