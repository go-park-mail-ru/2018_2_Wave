package room

import (
	"errors"
)

const (
	StatusError = "status_error"
)

var ( // const is inavaliable here
	ErrorNil              = errors.New("Nil input")
	ErrorNotExists        = errors.New("Not exists")
	ErrorAlreadyExists    = errors.New("Already exists")
	ErrorUnknownSignal    = errors.New("Unknown signal")
	ErrorConnectionClosed = errors.New("WS connection closed unexpected")
	ErrorWrongInputFormat = errors.New("Wrong input format")
)
