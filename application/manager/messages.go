package manager

import "Wave/application/room"

var (
	messageQSAcceptStatus = room.RouteResponse{Status: "quick_search_accept_status"}.WithStruct("")
	messageQSFailed       = room.RouteResponse{Status: "quick_search_failed"}.WithStruct("")
	messageQSStatus       = room.RouteResponse{Status: "quick_search_status"}.WithStruct("")
	messageQSReady        = room.RouteResponse{Status: "quick_search_ready"}.WithStruct("")
	messageQSDone         = room.RouteResponse{Status: "quick_search_done"}.WithStruct("")
	messageQSKick         = room.RouteResponse{Status: "quick_search_kick"}.WithStruct("")
)

// easyjson:json
type roomTokenPayload struct {
	RoomToken room.RoomToken `json:"room_token"`
}

// easyjson:json
type roomTypePayload struct {
	RoomType room.RoomType `json:"room_type"`
}

// easyjson:json
type userTokenPayload struct {
	UserToken room.UserID `json:"user_token"`
}

// easyjson:json
type roomInfoPayload struct {
	RoomToken room.RoomToken `json:"room_token"`
	RoomType  room.RoomType  `json:"room_type"`
}

// easyjson:json
type QSPayload struct {
	PlayerCount int           `json:"player_count,string"`
	RoomType    room.RoomType `json:"room_type"`
}

func (q *QSPayload) IsValid() bool {
	return q.PlayerCount <= 4 && q.PlayerCount >= 1 && IsRegisteredType(q.RoomType)
}

// easyjson:json
type QSStatusPayload struct {
	Members []userTokenPayload `json:"members"`
}

// easyjson:json
type QSReadyPayload struct {
	AcceptTimeout int `json:"accept_timeout"`
}

// easyjson:json
type QSAcceptPayload struct {
	Status bool `json:"status"`
}
