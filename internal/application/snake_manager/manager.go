package manager

import (
	"Wave/internal/application/proto"
	"Wave/internal/metrics"

	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//go:generate easyjson .

// Manager - main service room
/* - creates and store other rooms
 * - contains ALL online users
 * - provides all service functions
 */
type Manager struct {
	*proto.Room // the room super
	rooms      map[proto.RoomToken]proto.IRoom
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
func New(token proto.RoomToken, step time.Duration, db interface{}, prof *metrics.Profiler) *Manager {
	m := &Manager{
		Room:    proto.NewRoom(token, RoomType, nil, step),
		rooms:   map[proto.RoomToken]proto.IRoom{},
		builder: newBuilder(),
		prof:    prof,
		db:      db,
	}
	m.OnUserRemove = m.onUserRemoved
	m.Routes["lobby_list"] = m.onGetLobbyList
	m.Routes["lobby_create"] = m.withRoomType(m.onLobbyCreate)
	m.Routes["lobby_delete"] = m.withRoomToken(m.onLobbyDelete)

	m.Routes["add_to_room"] = m.withRoomToken(m.onAddToRoom)
	m.Routes["remove_from_room"] = m.withRoomToken(m.onRemoveFromRoom)

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
func (m *Manager) GetNextUserID() proto.UserToken {
	m.idsMutex.Lock()
	defer m.idsMutex.Unlock()

	m.lastUserID++
	return proto.UserToken(strconv.FormatInt(m.lastUserID, 36))
}

// GetNextRoomID returns next room id
func (m *Manager) GetNextRoomID() proto.RoomToken {
	m.idsMutex.Lock()
	defer m.idsMutex.Unlock()

	m.lastRoomID++
	return proto.RoomToken(strconv.FormatInt(m.lastRoomID, 36))
}

// CreateLobby -
func (m *Manager) CreateLobby(roomToken proto.RoomToken, u proto.IUser, roomType proto.RoomType) (proto.IRoom, error) {
	if factory, ok := type2Factory[roomType]; ok {
		r := factory(roomToken, m, m.db, m.GetTickTime())
		if r == nil {
			return nil, proto.ErrorNotFound
		}
		r.SetLogger(m.LG)
		m.rooms[roomToken] = r
		go r.Start()

		// profiler
		if m.prof != nil {
			m.prof.ActiveRooms.Inc()
		}

		return r, nil
	}
	return nil, proto.ErrorNotExists
}

// RemoveLobby -
func (m *Manager) RemoveLobby(roomToken proto.RoomToken, u proto.IUser) error {
	if r, ok := m.rooms[roomToken]; ok {
		if r.IsAbleToRemove(u) || u == nil {
			r.Stop()
			delete(m.rooms, roomToken)

			// profiler
			if m.prof != nil {
				m.prof.ActiveRooms.Dec()
			}
			return nil
		}
		return proto.ErrorForbiden
	}
	return proto.ErrorNotFound
}

func (m *Manager) CreateUser(name string, conn *websocket.Conn) (proto.IUser, error) {
	if conn == nil {
		return nil, proto.ErrorNil
	}

	u, err := proto.NewUser(m.GetNextUserID(), conn, m)
	if err != nil {
		return nil, err
	}
	u.Name = name
	u.SetLogger(m.LG)
	u.EnterRoom(m)

	return u, nil
}

// ----------------| handlers

func (m *Manager) onUserRemoved(u proto.IUser) {
	m.builder.RemoveUser(u) // kick from quick serach
}

func (m *Manager) onGetLobbyList(u proto.IUser, im proto.IInMessage) {
	data := []roomInfoPayload{}
	for _, r := range m.rooms {
		data = append(data, roomInfoPayload{
			RoomToken: r.GetToken(),
			RoomType:  r.GetType(),
		})
	}
	m.SendTo(u, messageOK.WithStruct(data))
}

func (m *Manager) onLobbyCreate(u proto.IUser, roomType proto.RoomType) {
	r, err := m.CreateLobby(m.GetNextRoomID(), u, roomType)
	if err != nil {
		return
	}
	m.SendTo(u, messageOK.WithStruct(roomTokenPayload{
		RoomToken: r.GetToken(),
	}))
}

func (m *Manager) onLobbyDelete(u proto.IUser, cmd proto.RoomToken) {
	m.RemoveLobby(cmd, u)
}

func (m *Manager) onAddToRoom(u proto.IUser, cmd proto.RoomToken) {
	if r, ok := m.rooms[cmd]; ok {
		u.Task(func() { u.EnterRoom(r) })
	}
}

func (m *Manager) onRemoveFromRoom(u proto.IUser, cmd proto.RoomToken) {
	if r, ok := m.rooms[cmd]; ok {
		u.Task(func() { u.ExitRoom(r) })
	}
}

// ------| quick serarch

// -> quick_search
func (m *Manager) onQSBegin(u proto.IUser, im proto.IInMessage) {
	p := &QSPayload{}
	if err := im.ToStruct(p); err != nil {
		return
	}
	if !p.IsValid() {
		return
	}
	m.builder.AddUser(u, p.RoomType, p.PlayerCount)
}

// -> quick_search_abort
func (m *Manager) onQSAbort(u proto.IUser, im proto.IInMessage) {
	m.builder.RemoveUser(u)
}

// -> quick_search_accept
func (m *Manager) onQSAccept(u proto.IUser, im proto.IInMessage) {
	p := &QSAcceptPayload{}
	if err := im.ToStruct(p); err != nil {
		return
	}
	m.builder.Accept(u, p.Status)
}

// <- quick_search_removed | quick_search_kick
func (m *Manager) onQSRemoved(f *former, u proto.IUser) {
	m.SendTo(u, messageQSKick)
	m.onQSStatus(f, messageQSRemoved)
	m.Logf("search removed: u=%s", u.GetToken())
}

// <- quick_search_added
func (m *Manager) onQSAdded(f *former, u proto.IUser) {
	m.onQSStatus(f, messageQSAdded)
	m.Logf("search added: u=%s", u.GetToken())
}

// <- quick_search_accept_status
func (m *Manager) onQSAcceptStatus(f *former, u proto.IUser) {
	p := &QSStatusPayload{}
	for _, u := range f.users {
		if !u.bAccepted {
			continue
		}
		p.Members = append(p.Members, QSStatusMemberPayload{
			UserToken:  u.GetToken(),
			UserName:   u.GetName(),
			UserSerial: f.GetUserSerial(u),
		})
	}
	om := messageQSAcceptStatus.WithStruct(p)
	for _, u := range f.users {
		m.SendTo(u, om)
	}
	m.Logf("search accept", u.GetToken())
}

func (m *Manager) onQSStatus(f *former, om *proto.Response) {
	p := &QSStatusPayload{}
	for _, u := range f.users {
		p.Members = append(p.Members, QSStatusMemberPayload{
			UserToken:  u.GetToken(),
			UserName:   u.GetName(),
			UserSerial: f.GetUserSerial(u),
		})
	}
	om = om.WithStruct(p)
	for _, u := range f.users {
		m.SendTo(u, om)
	}
}

// <- quick_search_ready
func (m *Manager) onQSReady(f *former) {
	p := &QSReadyPayload{
		AcceptTimeout: m.builder.acceptTime,
	}
	om := messageQSReady.WithStruct(p)
	for _, u := range f.users {
		m.SendTo(u, om)
	}
	m.Logf("search is ready")
}

// <- quick_search_done
func (m *Manager) onQSDone(f *former) {
	r, err := m.CreateLobby(m.GetNextRoomID(), nil, f.rType)
	if err != nil {
		return // TODO::
	}

	om := messageQSDone.WithStruct(roomTokenPayload{
		RoomToken: r.GetToken(),
	})
	for _, u := range f.users {
		m.SendTo(u, om)
		u.Task(func() { u.EnterRoom(r) }) 
	}
	m.Logf("search done")
}

// <- quick_search_failed
func (m *Manager) onQSFailed(f *former) {
	for _, u := range f.users {
		m.SendTo(u, messageQSFailed)
	}
	m.Logf("search failed")
}

// ----------------| helper functions

func (m *Manager) withRoomToken(next func(proto.IUser, proto.RoomToken)) proto.Route {
	return func(u proto.IUser, im proto.IInMessage) {
		cmd := &roomTokenPayload{}
		if im.ToStruct(cmd) == nil {
			next(u, cmd.RoomToken)
		} else {
			m.Log("message is not a roomTokenPayload", 
				"who", "withRoomToken",
				"where", "manager.go")
		}
	}
}

func (m *Manager) withRoomType(next func(proto.IUser, proto.RoomType)) proto.Route {
	return func(u proto.IUser, im proto.IInMessage) {
		cmd := &roomTypePayload{}
		if im.ToStruct(cmd) == nil {
			next(u, cmd.RoomType)
		}
		m.Log("message is not a roomTypePayload", 
				"who", "withRoomToken",
				"where", "manager.go")
	}
}
