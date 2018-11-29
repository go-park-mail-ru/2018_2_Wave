package room

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMessages(t *testing.T) {
	type S struct {
		F1 string
		F2 int
	}
	type Case struct {
		src interface{}
		res interface{}
	}
	cases := []Case{
		{"test string", ""},
		{&S{"gg", 8569}, &S{}},
	}

	for i, cs := range cases {
		bin, err := json.Marshal(cs.src)
		if err != nil {
			t.Fatal(err)
		}
		{ // InMessage
			msg := InMessage{Payload: bin}
			if err := msg.ToStruct(&cs.res); err != nil {
				t.Fatal(err, i)
			}
			if !reflect.DeepEqual(cs.src, cs.res) {
				t.Fatal("Not equal", i, "InMessage")
			}
		}
		{ // OutMessage
			msg := OutMessage{}
			if err := msg.FromStruct(cs.src); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(bin, msg.GetPayload()) {
				t.Fatal("Not equal", i, "OutMessage")
			}
		}
		{ // RouteResponse
			msg := RouteResponse{}
			if err := msg.FromStruct(cs.src); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(bin, msg.GetPayload()) {
				t.Fatal("Not equal", i, "RouteResponse")
			}
		}
	}
}

func TestPayload(t *testing.T) {
	data := []byte(`{
		"room_token": "u1",
		"signal": "s1",
		"payload": {
			"test": "val"
		}
	}`)
	expected := `{"room_token":"","status":"","payload":{"test":"val"}}`
	type Pl struct {
		Test string `json:"test"`
	}

	im := &InMessage{}
	if err := json.Unmarshal(data, im); err != nil {
		t.Fatal(err)
	}

	pl := &Pl{}
	if err := im.ToStruct(pl); err != nil {
		t.Fatal(err)
	}
	if pl.Test != "val" {
		t.Fatal("unexpected test value:", pl.Test)
	}

	om := &OutMessage{
		Payload: im.Payload,
	}
	res, err := json.Marshal(om)
	if err != nil {
		t.Fatal(err)
	}
	if string(res) != expected {
		t.Fatal("unexpected result:", string(res))
	}
}
