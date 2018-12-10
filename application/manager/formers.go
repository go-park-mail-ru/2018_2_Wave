package manager

import (
	"Wave/application/room"
	"sync"
)

// ----------------| room former

type roomFormerUser struct {
	room.IUser
	accepted bool
}

type roomFormer struct {
	users    []roomFormerUser
	roomType room.RoomType
	players  int

	onUserRemoved func(*roomFormer)
	onUserAdded   func(*roomFormer)
	onFormed      func(*roomFormer)
	onDone        func(*roomFormer)
}

func (r *roomFormer) AddUser(u room.IUser) {
	r.users = append(r.users, roomFormerUser{
		IUser: u,
	})
	if r.onUserAdded != nil {
		r.onUserAdded(r)
	}
	if r.IsFormed() {
		r.onFormed(r)
	}
}

func (r *roomFormer) RemoveUser(u room.IUser) {
	for i, expectant := range r.users {
		if expectant == u {
			r.users = append(r.users[:i], r.users[i+1:]...)
			if r.onUserRemoved != nil {
				r.onUserRemoved(r)
			}
			return
		}
	}
}

func (r *roomFormer) Accept(u room.IUser) {
	accepted := 0
	for i, expectant := range r.users {
		if expectant.accepted {
			accepted++
		}
		if expectant == u {
			r.users[i].accepted = true
		}
	}
	if accepted == len(r.users) && r.onDone != nil {
		r.onDone(r)
	}
}

func (r *roomFormer) IsFormed() bool {
	return len(r.users) >= r.players
}

// ----------------| rooms former

type roomsFormer struct {
	formers map[room.RoomType][]*roomFormer
	u2f     map[room.IUser]*roomFormer
	mu      sync.Mutex

	OnUserAdded   func(*roomFormer)
	OnUserRemoved func(*roomFormer)
	OnFormed      func(*roomFormer)
	OnDone        func(*roomFormer)
}

func newRoomsFormer() *roomsFormer {
	return &roomsFormer{
		formers: map[room.RoomType][]*roomFormer{},
		u2f:     map[room.IUser]*roomFormer{},
	}
}

func (rf *roomsFormer) AddUser(u room.IUser, roomType room.RoomType, players int) {
	rf.mu.Lock()
	{
		defer rf.mu.Unlock()

		// already searches
		if _, ok := rf.formers[roomType]; ok {
			return
		}
		f := rf.getFormer(roomType, players)
		rf.u2f[u] = f
		f.AddUser(u)
	}
}

func (rf *roomsFormer) RemoveUser(u room.IUser) {
	rf.mu.Lock()
	{
		defer rf.mu.Unlock()

		// not found
		f, ok := rf.u2f[u]
		if !ok {
			return
		}
		f.RemoveUser(u)
		delete(rf.u2f, u)
	}
}

func (rf *roomsFormer) Accept(u room.IUser) {
	rf.mu.Lock()
	{
		defer rf.mu.Unlock()

		// not found
		f, ok := rf.u2f[u]
		if !ok {
			return
		}
		f.Accept(u)
	}
}

func (rf *roomsFormer) GetUserFormer(u room.IUser) (*roomFormer, bool) {
	rf.mu.Lock()
	{
		defer rf.mu.Unlock()

		f, ok := rf.u2f[u]
		return f, ok
	}
}

// NOTE: the mutex must be taken before
func (rf *roomsFormer) getFormer(roomType room.RoomType, players int) *roomFormer {
	ff, ok := rf.formers[roomType]
	if !ok {
		ff = []*roomFormer{}
	}

	// try to find an existing former
	for _, f := range ff {
		if f.IsFormed() {
			continue
		}
		if f.players == players {
			return f
		}
	}

	// create a new former
	f := &roomFormer{
		players:       players,
		roomType:      roomType,
		onUserAdded:   rf.OnUserAdded,
		onUserRemoved: rf.OnUserRemoved,
		onFormed:      rf.OnFormed,
		onDone: func(f *roomFormer) {
			// remove the former from the former list and send it into the callback
			for i, expectant := range ff {
				if expectant != f {
					continue
				}
				ff = append(ff[:i], ff[i+1:]...)

				if rf.OnDone != nil {
					rf.OnDone(f)
				}
				return
			}
		},
	}

	rf.formers[roomType] = append(ff, f)
	return f
}
