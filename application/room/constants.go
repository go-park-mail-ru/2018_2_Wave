package room

import (
	"errors"
)

const (
	StatusError = "status_error"
	StatusTick  = "status_tick"
	StatusOK    = "status_ok"
)

var (
	MessageError         = RouteResponse{Status: StatusError}.WithStruct("")
	MessageWrongUserID   = RouteResponse{Status: StatusError}.WithStruct("wrong user id")
	MessageWrongRoomID   = RouteResponse{Status: StatusError}.WithStruct("wrong room id")
	MessageWrongRoomType = RouteResponse{Status: StatusError}.WithStruct("wrong room type")
	MessageWrongFormat   = RouteResponse{Status: StatusError}.WithStruct("wrong format")
	MessageForbiden      = RouteResponse{Status: StatusError}.WithStruct("forbiden")
	MessageTick          = RouteResponse{Status: StatusTick}.WithStruct("")
	MessageOK            = RouteResponse{Status: StatusOK}.WithStruct("")
)

var ( // const is inavaliable here
	ErrorNil              = errors.New("Nil input")
	ErrorForbiden         = errors.New("Forbiden")
	ErrorNotExists        = errors.New("Not exists")
	ErrorAlreadyExists    = errors.New("Already exists")
	ErrorUnknownSignal    = errors.New("Unknown signal")
	ErrorConnectionClosed = errors.New("WS connection closed unexpected")
	ErrorWrongInputFormat = errors.New("Wrong input format")
)

var (
	wsCloseErrors = func() (res []int) {
		for i := 1000; i <= 1010; i++ {
			res = append(res, i)
		}
		return append(res, 1015)
	}()
)
