package globalRoom

import (
	"Wave/server/room"
	"strconv"
	"time"
)

// GlobalRoom - main service room
/* - creates and store other rooms
 * - contains ALL online users
 * - provides all service functions
 */
type GlobalRoom struct {
	*room.Room // the room super

	/** internal rooms:
	 * 	- chats
	 *	- game lobbies
	 */
	internalRooms map[room.RoomID]room.IRoom

	lastRoomID int64
	lastUserID int64
}

func New(id room.RoomID, step time.Duration) *GlobalRoom {
	gr := &GlobalRoom{
		Room: room.NewRoom(id, step),
	}
	gr.Roures["lobby_list"] = gr.OnGetLobbyList
	gr.Roures["lobby_create"] = gr.OnLobbyCreate
	gr.Roures["lobby_delete"] = gr.OnLobbyDelete
	return gr
}

// ----------------| methods

func (gr *GlobalRoom) GetNextUserID() room.UserID {
	gr.lastUserID++
	return room.UserID(strconv.FormatInt(gr.lastUserID, 36))
}

func (gr *GlobalRoom) GetNextRoomID() room.RoomID {
	gr.lastRoomID++
	return room.RoomID(strconv.FormatInt(gr.lastRoomID, 36))
}

// ----------------| handlers

func (gr *GlobalRoom) OnGetLobbyList(u room.IUser, im room.IInMessage) room.IRouteResponce {
	type LobbyListItem struct {
		ID       room.RoomID
		RoomType room.RoomType
	}
	data := []LobbyListItem{}
	for _, r := range gr.internalRooms {
		data = append(data, LobbyListItem{
			ID:       r.GetID(),
			RoomType: r.GetType(),
		})
	}

	return room.RouteResponce{
		Status: room.StatusOK,
	}.WithStruct(data)
}

func (gr *GlobalRoom) OnLobbyCreate(u room.IUser, im room.IInMessage) room.IRouteResponce {
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
		r := factory(gr.GetNextRoomID(), gr.Step)
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

func (gr *GlobalRoom) OnLobbyDelete(u room.IUser, im room.IInMessage) room.IRouteResponce {
	type DeleteLobby struct {
		RoomID room.RoomID
	}
	cmd := DeleteLobby{}
	if im.ToStruct(cmd) != nil {
		return room.RouteResponce{
			Status: room.StatusError,
		}.WithStruct("Wrong input")
	}

	if r, ok := gr.internalRooms[cmd.RoomID]; ok {
		r.Stop()
		delete(gr.internalRooms, cmd.RoomID)
		return &room.RouteResponce{
			Status: room.StatusOK,
		}
	}
	return room.RouteResponce{
		Status: room.StatusError,
	}.WithStruct("Wrong id")
}
