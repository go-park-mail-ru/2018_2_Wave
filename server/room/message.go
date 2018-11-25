package room

import (
	"encoding/json"
)

// ----------------| InMessage

// InMessage - default IInMessage
type InMessage struct {
	RoomID  RoomID `json:"room_id"`
	Signal  string `json:"signal"`
	Payload []byte `json:"payload"`
}

func (im *InMessage) GetRoomID() RoomID  { return im.RoomID }
func (im *InMessage) GetSignal() string  { return im.Signal }
func (im *InMessage) GetPayload() []byte { return im.Payload }
func (im *InMessage) ToStruct(s interface{}) error {
	return json.Unmarshal(im.Payload, s)
}

// ----------------| OutMessage

// OutMessage - default IOutMessage
type OutMessage struct {
	RoomID  RoomID `json:"room_id"`
	Status  string `json:"status"`
	Payload []byte `json:"payload"`
}

func (om *OutMessage) GetRoomID() RoomID  { return om.RoomID }
func (om *OutMessage) GetStatus() string  { return om.Status }
func (om *OutMessage) GetPayload() []byte { return om.Payload }
func (om *OutMessage) FromStruct(s interface{}) (err error) {
	om.Payload, err = json.Marshal(s)
	return err
}

// ----------------| RouteResponse

// RouteResponse - default IOutMessage
type RouteResponse struct {
	Status  string `json:"status"`
	Payload []byte `json:"payload"`
}

func (om *RouteResponse) GetStatus() string  { return om.Status }
func (om *RouteResponse) GetPayload() []byte { return om.Payload }
func (om *RouteResponse) FromStruct(s interface{}) (err error) {
	om.Payload, err = json.Marshal(s)
	return err
}

// for usability
func (om RouteResponse) WithStruct(s interface{}) *RouteResponse {
	if err := om.FromStruct(s); err != nil {
		panic(err)
	}
	return &RouteResponse{
		Status:  om.Status,
		Payload: om.Payload,
	}
}
