package manager

import "Wave/internal/application/proto"

var (
	messageOK             = proto.Response{Status: "STATUS_OK"}.WithStruct("")
	messageQSAcceptStatus = proto.Response{Status: "quick_search_accept_status"}.WithStruct("")
	messageQSFailed       = proto.Response{Status: "quick_search_failed"}.WithStruct("")
	messageQSRemoved      = proto.Response{Status: "quick_search_removed"}.WithStruct("")
	messageQSAdded        = proto.Response{Status: "quick_search_added"}.WithStruct("")
	messageQSReady        = proto.Response{Status: "quick_search_ready"}.WithStruct("")
	messageQSDone         = proto.Response{Status: "quick_search_done"}.WithStruct("")
	messageQSKick         = proto.Response{Status: "quick_search_kick"}.WithStruct("")
)

// easyjson:json
type roomTokenPayload struct {
	RoomToken proto.RoomToken `json:"room_token"`
}

// easyjson:json
type roomTypePayload struct {
	RoomType proto.RoomType `json:"room_type"`
}

// easyjson:json
type userTokenPayload struct {
	UserToken proto.UserToken `json:"user_token"`
}

// easyjson:json
type roomInfoPayload struct {
	RoomToken proto.RoomToken `json:"room_token"`
	RoomType  proto.RoomType  `json:"room_type"`
}

// easyjson:json
type QSStatusMemberPayload struct {
	UserName   string          `json:"user_name"`
	UserToken  proto.UserToken `json:"user_token"`
	UserSerial int64           `json:"user_serial"`
}

// easyjson:json
type QSPayload struct {
	PlayerCount int            `json:"player_count"`
	RoomType    proto.RoomType `json:"room_type"`
}

func (q *QSPayload) IsValid() bool {
	return q.PlayerCount <= 4 && q.PlayerCount >= 1 && IsRegisteredType(q.RoomType)
}

// easyjson:json
type QSStatusPayload struct {
	Members []QSStatusMemberPayload `json:"members"`
}

// easyjson:json
type QSReadyPayload struct {
	AcceptTimeout int `json:"accept_timeout"`
}

// easyjson:json
type QSAcceptPayload struct {
	Status bool `json:"status"`
}
