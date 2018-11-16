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
		{ // RouteResponce
			msg := RouteResponce{}
			if err := msg.FromStruct(cs.src); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(bin, msg.GetPayload()) {
				t.Fatal("Not equal", i, "RouteResponce")
			}
		}
	}
}
