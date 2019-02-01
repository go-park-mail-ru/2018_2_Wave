package proto

import "encoding/json"

// ----------------| IResponse

// IResponse - response from a route
type IResponse interface {
	GetPayload() interface{}      // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// ----------------| Response

// Response - default IOutMessage
type Response struct {
	Status  string      `json:"status"`
	Payload interface{} `json:"payload"`
}

func (om *Response) GetStatus() string       { return om.Status }
func (om *Response) GetPayload() interface{} { return om.Payload }
func (om *Response) FromStruct(s interface{}) (err error) {
	om.Payload = s
	return nil
}

func (om Response) WithStruct(s interface{}) *Response {
	return &Response{
		Status:  om.Status,
		Payload: s,
	}
}

func (om Response) WithReason(reason string) *Response {
	type Reason struct {
		Reason string `json:"reason"`
	}
	return &Response{
		Status:  om.Status,
		Payload: Reason{reason},
	}
}

// ----------------| IInMessage

// IInMessage - message from a user
type IInMessage interface {
	GetRoomToken() RoomToken    // target room token
	GetSignal() string          // message method
	GetPayload() interface{}    // message payload
	ToStruct(interface{}) error // unmurshal data to struct
}

// ----------------| InMessage

// InMessage - default IInMessage
// easyjson:json
type InMessage struct {
	RoomToken RoomToken   `json:"room_token"`
	Signal    string      `json:"signal"`
	Payload   interface{} `json:"payload"`
}

func (im *InMessage) GetRoomToken() RoomToken { return im.RoomToken }
func (im *InMessage) GetSignal() string       { return im.Signal }
func (im *InMessage) GetPayload() interface{} { return im.Payload }
func (im *InMessage) ToStruct(s interface{}) error {
	data, err := json.Marshal(im.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

// ----------------| IOutMessage

// IOutMessage - message to a client
type IOutMessage interface {
	GetRoomToken() RoomToken      // message room token
	GetPayload() interface{}      // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// ----------------| OutMessage

// OutMessage - default IOutMessage
// easyjson:json
type OutMessage struct {
	RoomToken RoomToken   `json:"room_token"`
	Status    string      `json:"status"`
	Payload   interface{} `json:"payload"`
}

func (om *OutMessage) GetRoomToken() RoomToken { return om.RoomToken }
func (om *OutMessage) GetStatus() string       { return om.Status }
func (om *OutMessage) GetPayload() interface{} { return om.Payload }
func (om *OutMessage) FromStruct(s interface{}) (err error) {
	om.Payload, err = json.Marshal(s)
	return err
}
