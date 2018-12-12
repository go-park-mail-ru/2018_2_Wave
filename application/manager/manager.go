package manager

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

	former *roomsFormer
}

const RoomType = "manager"

// ----------------|

// New applicarion room
func New(id room.RoomToken, step time.Duration, db interface{}, prof *metrics.Profiler) *App {
	a := &App{
		Room:   room.NewRoom(id, RoomType, step),
		rooms:  map[room.RoomToken]room.IRoom{},
		former: newRoomsFormer(),
		prof:   prof,
		db:     db,
	}
	a.Routes["lobby_list"] = a.onGetLobbyList
	a.Routes["lobby_create"] = withRoomType(a.onLobbyCreate)
	a.Routes["lobby_delete"] = withRoomID(a.onLobbyDelete)

	a.Routes["add_to_room"] = withRoomID(a.onAddToRoom)
	a.Routes["remove_from_room"] = withRoomID(a.onRemoveFromRoom)

	a.Routes["quick_search"] = a.onQuickSearch
	a.Routes["quick_search_abort"] = a.onQuickSearchAbort
	a.Routes["quick_search_accept"] = a.onQuickSearchAccept
	a.former.OnUserAdded = a.onQuickSearchStatus
	a.former.OnUserRemoved = a.onQuickSearchStatus
	a.former.OnFormed = a.onQuickSearchReady
	a.former.OnDone = a.onQuickSearchDone

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

// ------| quick serarch

func (a *App) onQuickSearch(u room.IUser, im room.IInMessage) room.IRouteResponse {
	p := quickSearchPayload{}
	if err := im.ToStruct(p); err != nil {
		return nil
	}
	if !p.IsValid() {
		return nil
	}
	a.former.AddUser(u, p.RoomType, p.PlayerCount)
	return nil
}

func (a *App) onQuickSearchAbort(u room.IUser, im room.IInMessage) room.IRouteResponse {
	a.former.RemoveUser(u)
	return nil
}

func (a *App) onQuickSearchAccept(u room.IUser, im room.IInMessage) room.IRouteResponse {
	p := quickSearchAcceptPayload{}
	if err := im.ToStruct(p); err != nil {
		return nil
	}

	// remove the user from former and continue the search
	if p.Status == false {
		f, ok := a.former.GetUserFormer(u)
		if !ok {
			return nil
		}
		a.former.RemoveUser(u)

		for _, u := range f.users {
			a.SendMessageTo(u, messageQuickSearchFailed)
		}
	}
	// accept the game 
	a.former.Accept(u)
	return nil
}

func (a *App) onQuickSearchStatus(f *roomFormer) {
	p := quickSearchStatusPayload{}
	for _, u := range f.users {
		p.Members = append(p.Members, userTokenPayload{
			UserToken: u.GetID(),
		})
	}

	om := messageQuickSearchStatus.WithStruct(p)
	for _, u := range f.users {
		a.SendMessageTo(u, om)
	}
}

func (a *App) onQuickSearchReady(f *roomFormer) {
	p := quickSearchReadyPayload{
		AcceptTimeout: 30,
	}
	om := messageQuickSearchReady.WithStruct(p)
	for _, u := range f.users {
		a.SendMessageTo(u, om)
	}
}

func (a *App) onQuickSearchDone(f *roomFormer) {
	r, err := a.CreateLobby(f.roomType, a.GetNextRoomID())
	if err != nil {
		return // TODO::
	}

	om := messageQuickSearchDone.WithStruct(roomTokenPayload{
		RoomToken: r.GetID(),
	})
	for _, u := range f.users {
		a.SendMessageTo(u, om)
		u.AddToRoom(r)
	}
}

// ----------------| helper functions

var (
	messageQuickSearchStatus       = room.RouteResponse{Status: "quick_search_status"}.WithStruct("")
	messageQuickSearchReady        = room.RouteResponse{Status: "quick_search_ready"}.WithStruct("")
	messageQuickSearchAcceptStatus = room.RouteResponse{Status: "quick_search_accept_status"}.WithStruct("")
	messageQuickSearchDone         = room.RouteResponse{Status: "quick_search_done"}.WithStruct("")
	messageQuickSearchFailed       = room.RouteResponse{Status: "quick_search_failed"}.WithStruct("")
)

// easyjson:json
type roomTokenPayload struct {
	RoomToken room.RoomToken `json:"room_token"`
}

// easyjson:json
type roomTypePayload struct {
	RoomType room.RoomType `json:"room_type"`
}

// easyjson:json
type userTokenPayload struct {
	UserToken room.UserID `json:"user_token"`
}

// easyjson:json
type roomInfoPayload struct {
	RoomToken room.RoomToken `json:"room_token"`
	RoomType  room.RoomType  `json:"room_type"`
}

// easyjson:json
type quickSearchPayload struct {
	PlayerCount int           `json:"player_count"`
	RoomType    room.RoomType `json:"room_type"`
}

func (q *quickSearchPayload) IsValid() bool {
	return q.PlayerCount <= 4 && q.PlayerCount >= 1 && IsRegisteredType(q.RoomType)
}

// easyjson:json
type quickSearchStatusPayload struct {
	Members []userTokenPayload `json:"members"`
}

// easyjson:json
type quickSearchReadyPayload struct {
	AcceptTimeout int `json:"accept_timeout"`
}

// easyjson:json
type quickSearchAcceptPayload struct {
	Status bool `json:"status"`
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
