package proto

import (
	"encoding/json"
	"testing"
)

type messagePayloadMock struct {
	Name string `json:"name"`
}

func TestInMessageUnmarshal(t *testing.T) {
	const data = `
	{
		"room_token": "zeus",
		"signal": "ggwp",
		"payload": {
			"name": "kate"
		}
	}`
	im := &InMessage{}
	if err := json.Unmarshal([]byte(data), im); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if im.RoomToken != "zeus" {
		t.Fatalf("Unexpected room token %v, expected %v", im.RoomToken, "zeus")
	}
	if im.Signal != "ggwp" {
		t.Fatalf("Unexpected signal %v, expected %v", im.Signal, "ggwp")
	}
	if im.Payload == nil {
		t.Fatalf("Pauload must be non empty")
	}
	payload := &messagePayloadMock{}
	if err := im.ToStruct(payload); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if payload.Name != "kate" {
		t.Fatalf("Unexpected name %v, expected %v", payload.Name, "kate")
	}
}

func TestInMessageMarshal(t *testing.T) {
	var (
		om = OutMessage{
			RoomToken: "zeus",
			Status:    "ok",
			Payload: messagePayloadMock{
				Name: "kate",
			},
		}
		expected = `{"room_token":"zeus","status":"ok","payload":{"name":"kate"}}`
	)
	data, err := json.Marshal(om)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if expected != string(data) {
		t.Fatalf("Unexpected reslt %v expected %v", string(data), expected)
	}
}
