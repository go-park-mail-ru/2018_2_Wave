package manager

import (
	"Wave/internal/application/proto"
	"sync"
	"time"
)

// ----------------| stage

type formingStage int

const (
	stageForming formingStage = iota
	stageAccepting
	stageDone
)

// ----------------| former

type formerUser struct {
	proto.IUser
	bAccepted bool
}

type former struct {
	users   []formerUser
	stage   formingStage
	rType   proto.RoomType
	aim     int
	counter proto.Counter

	onUserRemoved func(*former, proto.IUser)
	onUserAdded   func(*former, proto.IUser)
	onAcceped     func(*former, proto.IUser)
	onFormed      func(*former)
	onFailed      func(*former)
	onDone        func(*former)
}

func (f *former) AddUser(u proto.IUser) {
	if f.stage != stageForming {
		return
	}
	if f.IsFormed() {
		return
	}

	f.users = append(f.users, formerUser{u, false})
	f.counter.Register(u.GetToken())
	if f.onUserAdded != nil {
		f.onUserAdded(f, u)
	}

	if f.IsFormed() {
		f.stage = stageAccepting
		if f.onFormed != nil {
			f.onFormed(f)
		}
	}
}

func (f *former) RemoveUser(u proto.IUser) {
	if f.stage == stageForming {
		f.removeUser(u)
		return
	}
	if f.stage == stageAccepting {
		f.Accept(u, false)
		return
	}
}

func (f *former) Accept(u proto.IUser, bAccept bool) {
	if f.stage != stageAccepting {
		return
	}

	accepted := 0
	for i, expectant := range f.users {
		if expectant.bAccepted {
			accepted++
		}
		if expectant.IUser == u && bAccept {
			f.users[i].bAccepted = true
			accepted++
			if f.onAcceped != nil {
				f.onAcceped(f, u)
			}
		} else if expectant.IUser == u {
			f.removeUser(u)
			f.stage = stageForming
			if f.onFailed != nil {
				f.onFailed(f)
			}
			return
		}
	}
	if accepted == len(f.users) && f.onDone != nil {
		f.stage = stageDone
		f.onDone(f)
	}
}

func (f *former) StopAccepting() {
	if f.stage != stageAccepting {
		return
	}

	for _, expectant := range f.users {
		if expectant.bAccepted {
			continue
		}
		f.removeUser(expectant.IUser)
	}
	f.stage = stageForming
	if f.onFailed != nil {
		f.onFailed(f)
	}
}

func (f *former) GetUserSerial(u proto.IUser) int64 {
	c, _ := f.counter.GetUserID(u)
	return c
}

func (f *former) IsFormed() bool {
	return len(f.users) >= f.aim
}

func (f *former) removeUser(u proto.IUser) {
	for i, expectant := range f.users {
		if expectant.IUser == u {
			f.users = append(f.users[:i], f.users[i+1:]...)
			f.counter.Unregister(u.GetToken())

			if f.onUserRemoved != nil {
				f.onUserRemoved(f, u)
			}
		}
	}
}

// ----------------| builder

type builder struct {
	formers    map[proto.RoomType][]*former
	u2f        map[proto.IUser]*former
	mu         sync.Mutex
	acceptTime int // seconds

	OnUserRemoved func(*former, proto.IUser)
	OnUserAdded   func(*former, proto.IUser)
	OnAcceped     func(*former, proto.IUser)
	OnFormed      func(*former)
	OnFailed      func(*former)
	OnDone        func(*former)
}

func newBuilder() *builder {
	b := &builder{
		formers: make(map[proto.RoomType][]*former),
		u2f:     make(map[proto.IUser]*former),
	}
	return b
}

func (b *builder) AddUser(u proto.IUser, roomType proto.RoomType, players int) {
	b.lock(func() {
		// if not searches
		if _, ok := b.u2f[u]; !ok {
			b.getFormer(roomType, players).AddUser(u)
		}
	})
}

func (b *builder) RemoveUser(u proto.IUser) {
	b.lock(func() {
		b.removeUser(u)
	})
}

func (b *builder) Accept(u proto.IUser, bAccept bool) {
	b.lock(func() {
		if f, ok := b.u2f[u]; ok {
			// if searches
			f.Accept(u, bAccept)
		}
	})
}

// ----------------| internal

func (b *builder) removeUser(u proto.IUser) {
	if f, ok := b.u2f[u]; ok {
		// if searches
		f.RemoveUser(u)
	}
}

func (b *builder) getFormer(roomType proto.RoomType, players int) *former {
	ff, ok := b.formers[roomType]
	if !ok {
		ff = []*former{}
	}

	// try to find an existing former
	for _, f := range ff {
		// already full
		if f.IsFormed() {
			continue
		}
		if f.aim == players {
			return f
		}
	}

	// create a new former
	f := &former{
		aim:       players,
		rType:     roomType,
		counter:   proto.MakeCounter(proto.FillGaps),
		onFailed:  b.OnFailed,
		onAcceped: b.OnAcceped,
		onUserAdded: func(f *former, u proto.IUser) {
			b.u2f[u] = f
			if b.OnUserAdded != nil {
				b.OnUserAdded(f, u)
			}
		},
		onUserRemoved: func(f *former, u proto.IUser) {
			delete(b.u2f, u)
			if b.OnUserRemoved != nil {
				b.OnUserRemoved(f, u)
			}
		},
		onFormed: func(f *former) {
			go func() { // dalay and remove all out of time users
				time.Sleep(time.Duration(b.acceptTime) * time.Second)

				b.mu.Lock()
				defer b.mu.Unlock()
				f.StopAccepting()
			}()
			if b.OnFormed != nil {
				b.OnFormed(f)
			}
		},
		onDone: func(f *former) {
			// remove the former from the former list and send it into the callback
			ff := b.formers[roomType]
			for i, expectant := range ff {
				if expectant != f {
					continue
				}
				// remove the former
				ff = append(ff[:i], ff[i+1:]...)
				b.formers[roomType] = ff
				if b.OnDone != nil {
					b.OnDone(f)
				}
				// remove users
				for _, u := range f.users {
					delete(b.u2f, u.IUser)
				}
				return
			}
		},
	}
	b.formers[roomType] = append(ff, f)
	return f
}

func (b *builder) lock(code func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	code()
}
