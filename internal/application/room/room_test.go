package room

import (
	"reflect"
	"testing"
	"time"
)

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

// ----------------| Room test

func TestRoomSimple(t *testing.T) {
	room := func() IRoom {
		r := NewRoom("test0", "", 30*time.Millisecond)
		r.Type = "test_type"
		r.Routes["echo"] = func(u IUser, im IInMessage) IRouteResponse {
			return &RouteResponse{
				Status:  "OK",
				Payload: im.GetPayload(),
			}
		}
		r.Routes["broad"] = func(u IUser, im IInMessage) IRouteResponse {
			r.Broadcast(&RouteResponse{
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
			t.Fatalf("Unexpected id: %s", id)
		}
		if tp := room.GetType(); tp != "test_type" {
			t.Fatalf("Unexpected type: %s", tp)
		}
	}
	{ // start serving
		go func() {
			if err := room.Run(); err != nil {
				t.Fatalf("Unexpected error: %v\n", err)
			}
		}()
	}
	{ // Add user
		if err := user.AddToRoom(room); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if err := user.AddToRoom(room); err != ErrorAlreadyExists {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if err := user2.AddToRoom(room); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
	}
	{ // apply message
		in1 := &InMessage{
			RoomToken: "test0",
			Signal:    "echo",
			Payload:   []byte("test text"),
		}
		in2 := &InMessage{
			RoomToken: "test0",
			Signal:    "broad",
			Payload:   []byte("test text"),
		}
		msg := &OutMessage{
			RoomToken: "test0",
			Status:    "OK",
			Payload:   []byte("test text"),
		}
		{ // echo (no broadcast)
			if err := room.ApplyMessage(user, in1); err != nil {
				t.Fatalf("Unexpected error: %v\n", err)
			}
			if len(user2.Messages) != 0 {
				t.Fatalf("Unexpected messages in user 2: %v", user2.Messages)
			}
			if len(user.Messages) != 1 {
				t.Fatalf("Expected only message. But: %v", user.Messages)
			}
			if !reflect.DeepEqual(user.Messages[0], msg) {
				t.Fatalf("Expected message: %v\n. Taken: %v", msg, user.Messages[0])
			}
			user.Messages = nil
		}
		{ // broadcast (no response)
			if err := room.ApplyMessage(user, in2); err != nil {
				t.Fatalf("Unexpected error: %v\n", err)
			}
			// @room.broadcast works in another thread
			time.Sleep(20 * time.Microsecond)

			if len(user2.Messages) != 1 {
				t.Fatalf("Expected only message. But: %v", user2.Messages)
			}
			if len(user.Messages) != 1 {
				t.Fatalf("Expected only message. But: %v", user.Messages)
			}
			if !reflect.DeepEqual(user2.Messages[0], msg) {
				t.Fatalf("Expected message: %v\n. Taken: %v", msg, user.Messages[0])
			}
			if !reflect.DeepEqual(user.Messages[0], msg) {
				t.Fatalf("Expected message: %v\n. Taken: %v", msg, user.Messages[0])
			}
			user.Messages = nil
			user2.Messages = nil
		}
	}
	{ // remove user
		if err := user.RemoveFromRoom(room); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if err := user.RemoveFromRoom(room); err != ErrorNotExists {
			t.Fatalf("Unexpected error: %v\n", err)
		}
		if err := user2.RemoveFromRoom(room); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
	}
	{ // stop serving
		if err := room.Stop(); err != nil {
			t.Fatalf("Unexpected error: %v\n", err)
		}
	}
}
