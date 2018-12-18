package manager

import (
	"Wave/internal/application/room"
	"Wave/internal/metrics"
	"strconv"
	"sync"
	"time"
)

//go:generate easyjson .

// Manager - main service room
/* - creates and store other rooms
 * - contains ALL online users
 * - provides all service functions
 */
type Manager struct {
	*room.Room // the room super
	rooms      map[room.RoomToken]room.IRoom
	db         interface{}
	prof       *metrics.Profiler

	lastRoomID int64
	lastUserID int64
	idsMutex   sync.Mutex

	builder *builder
}

const RoomType = "manager"

// ----------------|

// New applicarion room
func New(id room.RoomToken, step time.Duration, db interface{}, prof *metrics.Profiler) *Manager {
	m := &Manager{
		Room:    room.NewRoom(id, RoomType, step),
		rooms:   map[room.RoomToken]room.IRoom{},
		builder: newBuilder(),
		prof:    prof,
		db:      db,
	}
	m.OnUserRemoved = m.onUserRemoved
	m.Routes["lobby_list"] = m.onGetLobbyList
	m.Routes["lobby_create"] = withRoomType(m.onLobbyCreate)
	m.Routes["lobby_delete"] = withRoomID(m.onLobbyDelete)

	m.Routes["add_to_room"] = withRoomID(m.onAddToRoom)
	m.Routes["remove_from_room"] = withRoomID(m.onRemoveFromRoom)

	m.Routes["quick_search"] = m.onQSBegin
	m.Routes["quick_search_abort"] = m.onQSAbort
	m.Routes["quick_search_accept"] = m.onQSAccept
	m.builder.OnUserAdded = m.onQSAdded
	m.builder.OnUserRemoved = m.onQSRemoved
	m.builder.OnAcceped = m.onQSAcceptStatus
	m.builder.OnFormed = m.onQSReady
	m.builder.OnFailed = m.onQSFailed
	m.builder.OnDone = m.onQSDone
	m.builder.acceptTime = 30 /*s*/

	return m
}

// ----------------| methods

// GetNextUserID returns next user id
func (m *Manager) GetNextUserID() room.UserID {
	m.idsMutex.Lock()
	defer m.idsMutex.Unlock()

	m.lastUserID++
	return room.UserID(strconv.FormatInt(m.lastUserID, 36))
}

// GetNextRoomID returns next room id
func (m *Manager) GetNextRoomID() room.RoomToken {
	m.idsMutex.Lock()
	defer m.idsMutex.Unlock()

	m.lastRoomID++
	return room.RoomToken(strconv.FormatInt(m.lastRoomID, 36))
}

// CreateLobby -
func (m *Manager) CreateLobby(room_type room.RoomType, room_token room.RoomToken) (room.IRoom, error) {
	if factory, ok := type2Factory[room_type]; ok {
		r := factory(room_token, m.Step, m.db)
		if r == nil {
			return nil, room.ErrorNil
		}
		m.rooms[room_token] = r
		go r.Run()

		// profiler
		if m.prof != nil {
			m.prof.ActiveRooms.Inc()
		}

		return r, nil
	}
	return nil, room.ErrorNotExists
}

// ----------------| handlers

func (m *Manager) onUserRemoved(u room.IUser) {
	m.builder.RemoveUser(u) // kick from quick serach
}

func (m *Manager) onGetLobbyList(u room.IUser, im room.IInMessage) room.IRouteResponse {
	data := []roomInfoPayload{}
	for _, r := range m.rooms {
		data = append(data, roomInfoPayload{
			RoomToken: r.GetID(),
			RoomType:  r.GetType(),
		})
	}
	return room.MessageOK.WithStruct(data)
}

func (m *Manager) onLobbyCreate(u room.IUser, im room.IInMessage, cmd room.RoomType) room.IRouteResponse {
	r, err := m.CreateLobby(cmd, m.GetNextRoomID())
	if err != nil {
		return nil
	}
	return room.MessageOK.WithStruct(roomTokenPayload{
		RoomToken: r.GetID(),
	})
}

func (m *Manager) onLobbyDelete(u room.IUser, im room.IInMessage, cmd room.RoomToken) room.IRouteResponse {
	if r, ok := m.rooms[cmd]; ok {
		r.Stop()
		delete(m.rooms, cmd)

		// profiler
		if m.prof != nil {
			m.prof.ActiveRooms.Dec()
		}
	}
	return nil
}

func (m *Manager) onAddToRoom(u room.IUser, im room.IInMessage, cmd room.RoomToken) room.IRouteResponse {
	if r, ok := m.rooms[cmd]; ok {
		u.AddToRoom(r)
	}
	return nil
}

func (m *Manager) onRemoveFromRoom(u room.IUser, im room.IInMessage, cmd room.RoomToken) room.IRouteResponse {
	if r, ok := m.rooms[cmd]; ok {
		u.RemoveFromRoom(r)
	}
	return nil
}

// ------| quick serarch

// -> quick_search
func (m *Manager) onQSBegin(u room.IUser, im room.IInMessage) room.IRouteResponse {
	p := &QSPayload{}
	if err := im.ToStruct(p); err != nil {
		return nil
	}
	if !p.IsValid() {
		return nil
	}
	m.builder.AddUser(u, p.RoomType, p.PlayerCount)
	return nil
}

// -> quick_search_abort
func (m *Manager) onQSAbort(u room.IUser, im room.IInMessage) room.IRouteResponse {
	m.builder.RemoveUser(u)
	return nil
}

// -> quick_search_accept
func (m *Manager) onQSAccept(u room.IUser, im room.IInMessage) room.IRouteResponse {
	p := &QSAcceptPayload{}
	if err := im.ToStruct(p); err != nil {
		return nil
	}
	m.builder.Accept(u, p.Status)
	return nil
}

// <- quick_search_removed | quick_search_kick
func (m *Manager) onQSRemoved(f *former, u room.IUser) {
	m.SendMessageTo(u, messageQSKick)
	m.onQSStatus(f, messageQSRemoved)
}

// <- quick_search_added
func (m *Manager) onQSAdded(f *former, u room.IUser) {
	m.onQSStatus(f, messageQSAdded)
}

// <- quick_search_accept_status
func (m *Manager) onQSAcceptStatus(f *former, u room.IUser) {
	m.onQSStatus(f, messageQSAcceptStatus)
}

func (m *Manager) onQSStatus(f *former, om *room.RouteResponse) {
	p := &QSStatusPayload{}
	for _, u := range f.users {
		p.Members = append(p.Members, QSStatusMemberPayload{
			UserToken:  u.GetID(),
			UserName:   u.GetName(),
			UserSerial: f.GetUserSerial(u),
		})
	}
	om = om.WithStruct(p)
	for _, u := range f.users {
		m.SendMessageTo(u, om)
	}
}

// <- quick_search_ready
func (m *Manager) onQSReady(f *former) {
	p := &QSReadyPayload{
		AcceptTimeout: m.builder.acceptTime,
	}
	om := messageQSReady.WithStruct(p)
	for _, u := range f.users {
		m.SendMessageTo(u, om)
	}
}

// <- quick_search_done
func (m *Manager) onQSDone(f *former) {
	r, err := m.CreateLobby(f.rType, m.GetNextRoomID())
	if err != nil {
		return // TODO::
	}

	om := messageQSDone.WithStruct(roomTokenPayload{
		RoomToken: r.GetID(),
	})
	for _, u := range f.users {
		m.SendMessageTo(u, om)
		u.AddToRoom(r)
	}
}

// <- quick_search_failed
func (m *Manager) onQSFailed(f *former) {
	for _, u := range f.users {
		m.SendMessageTo(u, messageQSFailed)
	}
}

// ----------------| helper functions

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
