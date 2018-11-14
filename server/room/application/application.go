package application

import (
	"Wave/server/room"
	"strconv"
	"sync"
	"time"
)

// Application - main service room
/* - creates and store other rooms
 * - contains ALL online users
 * - provides all service functions
 */
type Application struct {
	*room.Room // the room super

	/** internal rooms:
	 * 	- chats
	 *	- game lobbies
	 */
	internalRooms map[room.RoomID]room.IRoom

	lastRoomID int64
	lastUserID int64
	idsMutex   sync.Mutex
}

func New(id room.RoomID, step time.Duration) *Application {
	a := &Application{
		Room: room.NewRoom(id, step),
	}
	a.Roures["lobby_list"] = a.OnGetLobbyList
	a.Roures["lobby_create"] = a.OnLobbyCreate
	a.Roures["lobby_delete"] = a.OnLobbyDelete
	return a
}

// ----------------| methods

func (a *Application) GetNextUserID() room.UserID {
	a.idsMutex.Lock()
	defer a.idsMutex.Unlock()

	a.lastUserID++
	return room.UserID(strconv.FormatInt(a.lastUserID, 36))
}

func (a *Application) GetNextRoomID() room.RoomID {
	a.idsMutex.Lock()
	defer a.idsMutex.Unlock()

	a.lastRoomID++
	return room.RoomID(strconv.FormatInt(a.lastRoomID, 36))
}

// ----------------| handlers

func (a *Application) OnGetLobbyList(u room.IUser, im room.IInMessage) room.IRouteResponce {
	type LobbyListItem struct {
		ID       room.RoomID
		RoomType room.RoomType
	}
	data := []LobbyListItem{}
	for _, r := range a.internalRooms {
		data = append(data, LobbyListItem{
			ID:       r.GetID(),
			RoomType: r.GetType(),
		})
	}

	return room.RouteResponce{
		Status: room.StatusOK,
	}.WithStruct(data)
}

func (a *Application) OnLobbyCreate(u room.IUser, im room.IInMessage) room.IRouteResponce {
	type CreateLobby struct {
		RoomType room.RoomType
	}
	cmd := &CreateLobby{}
	if im.ToStruct(cmd) != nil {
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("Wrong input")
	}

	if factory, ok := type2Factory[cmd.RoomType]; !ok {
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("Unknown room type")
	} else {
		r := factory(a.GetNextRoomID(), a.Step)
		if r == nil {
			return room.RouteResponce{
				Status: room.StatusError,
			}.WithStruct("Internal error")
		}
		go r.Run()
		u.AddToRoom(r)

		return room.RouteResponce{
			Status: room.StatusOK,
		}.WithStruct(r.GetID())
	}
}

func (a *Application) OnLobbyDelete(u room.IUser, im room.IInMessage) room.IRouteResponce {
	type DeleteLobby struct {
		RoomID room.RoomID
	}
	cmd := DeleteLobby{}
	if im.ToStruct(cmd) != nil {
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("Wrong input")
	}

	if r, ok := a.internalRooms[cmd.RoomID]; ok {
		r.Stop()
		delete(a.internalRooms, cmd.RoomID)
		return &room.RouteResponce{
			Status: room.StatusOK,
		}
	}
	return room.RouteResponce{
		Status: room.StatusError,
	}.WithStruct("Wrong id")
}
