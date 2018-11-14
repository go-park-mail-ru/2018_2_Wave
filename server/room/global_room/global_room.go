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

func (gr *GlobalRoom) OnGetLobbyList(u room.IUser, im room.IInMessage) {
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
	gr.SendMessageTo(u, room.StatusOK, data)
}

func (gr *GlobalRoom) OnLobbyCreate(u room.IUser, im room.IInMessage) {
	type CreateLobby struct {
		RoomType room.RoomType
	}
	cmd := &CreateLobby{}
	if im.ToStruct(cmd) != nil {
		gr.SendMessageTo(u, room.StatusError, "Wrong input")
		return
	}

	if factory, ok := type2Factory[cmd.RoomType]; !ok {
		gr.SendMessageTo(u, room.StatusError, "Unknown room type")
		return
	} else {
		r := factory(gr.GetNextRoomID(), gr.Step)
		if r == nil {
			gr.SendMessageTo(u, room.StatusError, "Internal error")
			return
		}
		go r.Run()
		u.AddToRoom(r)
	}
}

func (gr *GlobalRoom) OnLobbyDelete(u room.IUser, im room.IInMessage) {
	type DeleteLobby struct {
		RoomID room.RoomID
	}
	cmd := DeleteLobby{}
	if im.ToStruct(cmd) != nil {
		gr.SendMessageTo(u, room.StatusError, "Wrong input")
		return
	}

	if r, ok := gr.internalRooms[cmd.RoomID]; !ok {
		gr.SendMessageTo(u, room.StatusError, "Wrong id")
		return
	} else {
		r.Stop()
		delete(gr.internalRooms, cmd.RoomID)
	}
}
