package app

import (
	"Wave/server/chat/room"
	"strconv"
	"sync"
	"time"
)

// App - main service room
/* - creates and store other rooms
 * - contains ALL online users
 * - provides all service functions
 */
type App struct {
	*room.Room // the room super

	/** internal rooms:
	 * 	- chats
	 *	- game lobbies	*/
	internalRooms map[room.RoomID]room.IRoom

	lastRoomID int64
	lastUserID int64
	idsMutex   sync.Mutex
}

// New applicarion room
func New(id room.RoomID, step time.Duration) *App {
	a := &App{
		Room: room.NewRoom(id, step),
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
func (a *App) GetNextRoomID() room.RoomID {
	a.idsMutex.Lock()
	defer a.idsMutex.Unlock()

	a.lastRoomID++
	return room.RoomID(strconv.FormatInt(a.lastRoomID, 36))
}

// ----------------| handlers

func (a *App) onGetLobbyList(u room.IUser, im room.IInMessage) room.IRouteResponse {
	type Response struct {
		RoomID   room.RoomID
		RoomType room.RoomType
	}
	data := []Response{}
	for _, r := range a.internalRooms {
		data = append(data, Response{
			RoomID:   r.GetID(),
			RoomType: r.GetType(),
		})
	}

	return room.MessageOK.WithStruct(data)
}

func (a *App) onLobbyCreate(u room.IUser, im room.IInMessage, cmd room.RoomType) room.IRouteResponse {
	if factory, ok := type2Factory[cmd]; ok {
		r := factory(a.GetNextRoomID(), a.Step)
		if r == nil {
			return room.MessageError
		}
		go r.Run()

		u.AddToRoom(r)
		return room.MessageOK.WithStruct(r.GetID())
	}
	return room.MessageWrongRoomType
}

func (a *App) onLobbyDelete(u room.IUser, im room.IInMessage, cmd room.RoomID) room.IRouteResponse {
	if r, ok := a.internalRooms[cmd]; ok {
		r.Stop()
		delete(a.internalRooms, cmd)
		return room.MessageOK
	}
	return room.MessageWrongRoomID
}

func (a *App) onAddToRoom(u room.IUser, im room.IInMessage, cmd room.RoomID) room.IRouteResponse {
	if r, ok := a.internalRooms[cmd]; ok {
		if err := u.AddToRoom(r); err == nil {
			return room.MessageOK
		}
		return room.MessageForbiden
	}
	return room.MessageWrongRoomID
}

func (a *App) onRemoveFromRoom(u room.IUser, im room.IInMessage, cmd room.RoomID) room.IRouteResponse {
	if r, ok := a.internalRooms[cmd]; ok {
		if err := u.RemoveFromRoom(r); err == nil {
			return room.MessageOK
		}
		return room.MessageForbiden
	}
	return room.MessageWrongRoomID
}

// ----------------| helper functions

type roomIDPayload struct {
	RoomID room.RoomID
}

type roomTypePayload struct {
	RoomType room.RoomType
}

func withRoomID(next func(room.IUser, room.IInMessage, room.RoomID) room.IRouteResponse) room.Route {
	return func(u room.IUser, im room.IInMessage) room.IRouteResponse {
		cmd := &roomIDPayload{}
		if im.ToStruct(cmd) == nil {
			return next(u, im, cmd.RoomID)
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
