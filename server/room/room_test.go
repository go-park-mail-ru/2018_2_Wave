package room

import (
	"net/http"
	"reflect"
	"testing"
	"time"

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
func (tr *testUserHelper) GetType() RoomType {
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
func (tr *testUserHelper) ApplyMessage(u IUser, im IInMessage) error {
	for _, u := range tr.users {
		u.Consume(&OutMessage{
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

type testRoomSimpleHelper struct {
	ID       UserID
	Room     IRoom
	Messages []IOutMessage
}

func (t *testRoomSimpleHelper) GetID() UserID { return t.ID }
func (t *testRoomSimpleHelper) AddToRoom(r IRoom) error {
	t.Room = r
	return r.AddUser(t)
}
func (t *testRoomSimpleHelper) RemoveFromRoom(r IRoom) error {
	t.Room = nil
	return r.RemoveUser(t)
}
func (t *testRoomSimpleHelper) Listen() error        { return nil }
func (t *testRoomSimpleHelper) StopListening() error { return nil }
func (t *testRoomSimpleHelper) Consume(om IOutMessage) error {
	t.Messages = append(t.Messages, om)
	return nil
}

func TestRoomSimple(t *testing.T) {
	room := func() IRoom {
		r := NewRoom("test0", 30*time.Millisecond)
		r.Type = "test_type"
		r.Roures["echo"] = func(u IUser, im IInMessage) IRouteResponce {
			return &RouteResponce{
				Status:  "OK",
				Payload: im.GetPayload(),
			}
		}
		r.Roures["broad"] = func(u IUser, im IInMessage) IRouteResponce {
			r.Broadcast(&RouteResponce{
				Status:  "OK",
				Payload: im.GetPayload(),
			})
			return nil
		}
		return r
	}()
	user := &testRoomSimpleHelper{
		ID: "test ID",
	}
	user2 := &testRoomSimpleHelper{
		ID: "test ID2",
	}

	{ // returning values
		if id := room.GetID(); id != "test0" {
			t.Errorf("Unexpected id: %s", id)
		}
		if tp := room.GetType(); tp != "test_type" {
			t.Errorf("Unexpected type: %s", tp)
		}
	}
	{ // start serving
		go func() {
			if err := room.Run(); err != nil {
				t.Errorf("Unexpected error: %v\n", err)
			}
		}()
	}
	{ // Add user
		if err := user.AddToRoom(room); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if err := user.AddToRoom(room); err != ErrorAlreadyExists {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if err := user2.AddToRoom(room); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
	}
	{ // apply message
		in1 := &InMessage{
			RoomID:  "test0",
			Signal:  "echo",
			Payload: []byte("test text"),
		}
		in2 := &InMessage{
			RoomID:  "test0",
			Signal:  "broad",
			Payload: []byte("test text"),
		}
		msg := &OutMessage{
			RoomID:  "test0",
			Status:  "OK",
			Payload: []byte("test text"),
		}
		{ // echo (no broadcast)
			if err := room.ApplyMessage(user, in1); err != nil {
				t.Errorf("Unexpected error: %v\n", err)
			}
			if len(user2.Messages) != 0 {
				t.Errorf("Unexpected messages in user 2: %v", user2.Messages)
			}
			if len(user.Messages) != 1 {
				t.Errorf("Expected only message. But: %v", user.Messages)
			}
			if !reflect.DeepEqual(user.Messages[0], msg) {
				t.Errorf("Expected message: %v\n. Taken: %v", msg, user.Messages[0])
			}
			user.Messages = nil
		}
		{ // broadcast (no responce)
			if err := room.ApplyMessage(user, in2); err != nil {
				t.Errorf("Unexpected error: %v\n", err)
			}
			time.Sleep(10 * time.Microsecond) // broadcast works in another thread

			if len(user2.Messages) != 1 {
				t.Errorf("Expected only message. But: %v", user2.Messages)
			}
			if len(user.Messages) != 1 {
				t.Errorf("Expected only message. But: %v", user.Messages)
			}
			if !reflect.DeepEqual(user2.Messages[0], msg) {
				t.Errorf("Expected message: %v\n. Taken: %v", msg, user.Messages[0])
			}
			if !reflect.DeepEqual(user.Messages[0], msg) {
				t.Errorf("Expected message: %v\n. Taken: %v", msg, user.Messages[0])
			}
			user.Messages = nil
			user2.Messages = nil
		}
	}
	{ // remove user
		if err := user.RemoveFromRoom(room); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if err := user.RemoveFromRoom(room); err != ErrorNotExists {
			t.Errorf("Unexpected error: %v\n", err)
		}
		if err := user2.RemoveFromRoom(room); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
	}
	{ // stop serving
		if err := room.Stop(); err != nil {
			t.Errorf("Unexpected error: %v\n", err)
		}
	}
}
