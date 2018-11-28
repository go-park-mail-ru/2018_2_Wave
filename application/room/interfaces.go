package room

import "time"

/** How does it works?
 *	. ## begining of a WS handler
 *  . In WS handler we crate an instance if IUser interface
 *	. Than we call a Listen method in the instance in goroutine
 *	. The Listen() goes into a cycle and listen to the WS connection
 *	. Then in a main thread we call a AddToRoom() and send a main room into
 *	. ## and of the handler
 */

//go:generate easyjson .

type RoomID string
type UserID string
type RoomType string
type RoomFactory func(id RoomID, step time.Duration, db interface{}) IRoom

// IInMessage - message from a user
type IInMessage interface {
	GetRoomID() RoomID          // target room id
	GetSignal() string          // message method
	GetPayload() interface{}    // message payload
	ToStruct(interface{}) error // unmurshal data to struct
}

// IOutMessage - message to a client
type IOutMessage interface {
	GetRoomID() RoomID            // message room id
	GetPayload() interface{}      // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// IRouteResponse - response from route
type IRouteResponse interface {
	GetPayload() interface{}      // message payload
	GetStatus() string            // message status
	FromStruct(interface{}) error // marshal from struct
}

// IUser - client websocket wrapper
type IUser interface {
	GetID() UserID              // User id
	AddToRoom(IRoom) error      // order to add self into the room
	RemoveFromRoom(IRoom) error // order to romve self from the room
	Listen() error              // Listen to messages
	StopListening() error       // Stop listening
	Consume(IOutMessage) error  // Send message to user
}

// IRoom - abstruct room inteface
type IRoom interface {
	GetID() RoomID                        // room id
	GetType() RoomType                    // room type
	Run() error                           // run thr room
	Stop() error                          // stop the room
	AddUser(IUser) error                  // add the user to the room
	RemoveUser(IUser) error               // remove  the user from the room
	OnDisconnected(IUser)                 // inform the room the user was disconnected
	ApplyMessage(IUser, IInMessage) error // send message to the room
}
