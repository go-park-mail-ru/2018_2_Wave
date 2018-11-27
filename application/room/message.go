package room

import (
	"encoding/json"
)

// ----------------| InMessage

// InMessage - default IInMessage
type InMessage struct {
	RoomID  RoomID      `json:"room_id"`
	Signal  string      `json:"signal"`
	Payload interface{} `json:"payload"`
}

func (im *InMessage) GetRoomID() RoomID       { return im.RoomID }
func (im *InMessage) GetSignal() string       { return im.Signal }
func (im *InMessage) GetPayload() interface{} { return im.Payload }
func (im *InMessage) ToStruct(s interface{}) error {
	data, err := json.Marshal(im.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

// ----------------| OutMessage

// OutMessage - default IOutMessage
type OutMessage struct {
	RoomID  RoomID      `json:"room_id"`
	Status  string      `json:"status"`
	Payload interface{} `json:"payload"`
}

func (om *OutMessage) GetRoomID() RoomID       { return om.RoomID }
func (om *OutMessage) GetStatus() string       { return om.Status }
func (om *OutMessage) GetPayload() interface{} { return om.Payload }
func (om *OutMessage) FromStruct(s interface{}) (err error) {
	om.Payload, err = json.Marshal(s)
	return err
}

// ----------------| RouteResponse

// RouteResponse - default IOutMessage
type RouteResponse struct {
	Status  string      `json:"status"`
	Payload interface{} `json:"payload"`
}

func (om *RouteResponse) GetStatus() string       { return om.Status }
func (om *RouteResponse) GetPayload() interface{} { return om.Payload }
func (om *RouteResponse) FromStruct(s interface{}) (err error) {
	om.Payload = s
	return nil
}

func (om RouteResponse) WithStruct(s interface{}) *RouteResponse {
	return &RouteResponse{
		Status:  om.Status,
		Payload: s,
	}
}

func (om RouteResponse) WithReason(reason string) *RouteResponse {
	type Reason struct {
		Reason string `json:"reason"`
	}
	return &RouteResponse{
		Status:  om.Status,
		Payload: Reason{reason},
	}
}
