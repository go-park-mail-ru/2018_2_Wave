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
	MessageError         = RouteResponce{Status: StatusError}.WithStruct("")
	MessageWrongUserID   = RouteResponce{Status: StatusError}.WithStruct("wrong user id")
	MessageWrongRoomID   = RouteResponce{Status: StatusError}.WithStruct("wrong room id")
	MessageWrongRoomType = RouteResponce{Status: StatusError}.WithStruct("wrong room type")
	MessageWrongFormat   = RouteResponce{Status: StatusError}.WithStruct("wrong format")
	MessageForbiden      = RouteResponce{Status: StatusError}.WithStruct("forbiden")
	MessageTick          = RouteResponce{Status: StatusTick}.WithStruct("")
	MessageOK            = RouteResponce{Status: StatusOK}.WithStruct("")
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
