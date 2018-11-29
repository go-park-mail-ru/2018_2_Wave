package room

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

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

func (tr *testUserHelper) GetID() RoomToken     { return "test" }
func (tr *testUserHelper) GetType() RoomType { return "test" }
func (tr *testUserHelper) AddUser(u IUser) error {
	tr.users[u.GetID()] = u
	return nil
}

func (tr *testUserHelper) RemoveUser(u IUser) error {
	delete(tr.users, u.GetID())
	return nil
}

func (tr *testUserHelper) ApplyMessage(u IUser, im IInMessage) error {
	for _, u := range tr.users {
		u.Consume(&OutMessage{
			RoomToken:  im.GetRoomID(),
			Payload: im.(*InMessage).Payload,
			Status:  "OK",
		})
	}
	return nil
}

func (tr *testUserHelper) Run() error  { return nil }
func (tr *testUserHelper) Stop() error { return nil }
func (tr *testUserHelper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	con, err := tr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	tr.Con = con
}

// ----------------| User test

func TestUserSimple(t *testing.T) {
	const usrID = "test1"
	var (
		room        = newTestUserHelper()
		con, _, err = wstest.NewDialer(room).Dial("ws:/localhost:6660/", nil)
		user        = NewUser(usrID, con)
	)
	if err != nil {
		t.Fatal(err)
	}

	// ws connection initialises in another goroutine
	// see @room.ServeHTTP() and @wstest.NewDialer()
	time.Sleep(30 * time.Microsecond)

	{ // returning values
		if ID := user.GetID(); ID != usrID {
			t.Fatalf("Unexpected ID: %v\n", ID)
		}
	}
	{ // start listening
		go func() {
			if err := user.Listen(); err != nil {
				t.Fatalf("Unexpected error: %v\n", err)
			}
		}()
	}
	{ // add to room
		if err := user.AddToRoom(room); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if u, ok := room.users[usrID]; !ok {
			t.Fatalf("User expected but not found")
		} else if u != user {
			t.Fatalf("Unexpected user: %v\n", u)
		}
	}
	{ // send message
		const testText = "test text"
		to := &InMessage{
			RoomToken:  "test",
			Signal:  "test",
			Payload: []byte(testText),
		}
		res := &OutMessage{}

		if err := room.Con.WriteJSON(to); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if err := room.Con.ReadJSON(res); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if res.RoomToken != "test" {
			t.Fatalf("Unexpected roomType: %v\n", res.RoomToken)
		}
		if res.Payload.(string) != testText {
			t.Fatalf("Unexpected payload: %v\n", res.Payload.(string))
		}
	}
	{ // remove from room
		if err := user.RemoveFromRoom(room); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if u, ok := room.users[usrID]; ok {
			t.Fatalf("User wasn't expected but found: %v\n", u)
		}
	}
	{ // stop listening
		if err := user.StopListening(); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
	}
}
