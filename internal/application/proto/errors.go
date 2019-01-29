package proto

import "errors"

// ----------------| Errors

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
