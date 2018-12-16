package room

import (
	"errors"
)

const (
	StatusError = "STATUS_ERROR"
	StatusTick  = "STATUS_TICK"
	StatusOK    = "STATUS_OK"
)

var (
	MessageError         = RouteResponse{Status: StatusError}.WithReason("")
	MessageWrongUserID   = RouteResponse{Status: StatusError}.WithReason("wrong user id")
	MessageWrongRoomID   = RouteResponse{Status: StatusError}.WithReason("wrong room id")
	MessageWrongRoomType = RouteResponse{Status: StatusError}.WithReason("wrong room type")
	MessageWrongFormat   = RouteResponse{Status: StatusError}.WithReason("wrong format")
	MessageForbiden      = RouteResponse{Status: StatusError}.WithReason("forbiden")
	MessageTick          = RouteResponse{Status: StatusTick}.WithStruct("")
	MessageOK            = RouteResponse{Status: StatusOK}.WithStruct("")
)

var ( // const is inavaliable here
	ErrorNil              = errors.New("Nil input")
	ErrorForbiden         = errors.New("Forbiden")
	ErrorNotExists        = errors.New("Not exists")
	ErrorNotFound         = errors.New("Not found")
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
