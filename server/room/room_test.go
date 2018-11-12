package room

import (
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

// ----------------| User test

type testUserHelper struct {
	users    map[UserID]IUser
	Con      *websocket.Conn
	upgrader websocket.Upgrader
}

func newTestUserHelper() *testUserHelper {
	return &testUserHelper{
		users: map[UserID]IUser{},
	}
}
func (tr *testUserHelper) GetID() RoomID {
	return "test"
}
func (tr *testUserHelper) AddUser(u IUser) error {
	tr.users[u.GetID()] = u
	return nil
}
func (tr *testUserHelper) RemoveUser(u IUser) error {
	delete(tr.users, u.GetID())
	return nil
}
func (tr *testUserHelper) SendMessage(im IInMessage) error {
	for _, u := range tr.users {
		u.Send(&OutMessage{
			RoomID:  im.GetRoomID(),
			Payload: im.(*InMessage).Payload,
			Status:  "OK",
		})
	}
	return nil
}
func (tr *testUserHelper) Run() error {
	return nil
}
func (tr *testUserHelper) Stop() error {
	return nil
}
func (tr *testUserHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	con, err := tr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	tr.Con = con
}

func TestUserSimple(t *testing.T) {
	const usrID = "test1"
	var (
		room        = newTestUserHelper()
		con, _, err = wstest.NewDialer(room).Dial("ws:/localhost:6660/", nil)
		user        = NewUser(usrID, con)
	)
	if err != nil {
		t.Errorf(err.Error())
	}

	{ // returning values
		if ID := user.GetID(); ID != usrID {
			t.Errorf("Unexpected ID: %v\n", ID)
		}
	}
	{ // start listening
		go func() {
			if err := user.Listen(); err != nil {
				t.Errorf("Unexpected error: %v\n", err)
			}
		}()
	}
	{ // add to room
		if err := user.AddToRoom(room); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if u, ok := room.users[usrID]; !ok {
			t.Errorf("User expected but not found")
		} else if u != user {
			t.Errorf("Unexpected user: %v\n", u)
		}
	}
	{ // send message
		const testText = "test text"
		to := &InMessage{
			RoomID:  "test",
			Signal:  "test",
			Payload: []byte(testText),
		}
		res := &OutMessage{}

		if err := room.Con.WriteJSON(to); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if err := room.Con.ReadJSON(res); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if res.RoomID != "test" {
			t.Errorf("Unexpected roomID: %v\n", res.RoomID)
		}
		if string(res.Payload) != testText {
			t.Errorf("Unexpected payload: %v\n", string(res.Payload))
		}
	}
	{ // remove from room
		if err := user.RemoveFromRoom(room); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if u, ok := room.users[usrID]; ok {
			t.Errorf("User wasn't expected but found: %v\n", u)
		}
	}
	{ // stop listening
		if err := user.StopListening(); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
	}
}

// ----------------| Room test

// TODO:: test room

// func TestRoomSimple(t *testing.T) {
// 	var room IRoom = NewRoom("test0", 30*time.Millisecond)
// 	var user IUser = NewUser("test1", nil)
// 	println(room, user)
// }
