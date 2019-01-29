package manager

import (
	"Wave/internal/application/room"
	"Wave/internal/application/snake"
)

// type2Factory - singletone map converting room type to the room factory
// To add a new factory user @RegisterRoomType function
var type2Factory = map[room.RoomType]room.RoomFactory{}

// RegisterRoomType factory
func RegisterRoomType(roomType room.RoomType, factory room.RoomFactory) error {
	if IsRegisteredType(roomType) {
		return room.ErrorAlreadyExists
	}
	type2Factory[roomType] = factory
	return nil
}

// IsRegisteredType - weather the type factory was gegistered
func IsRegisteredType(roomType room.RoomType) bool {
	_, ok := type2Factory[roomType]
	return ok
}

func init() {
	for Type, Factory := range map[room.RoomType]room.RoomFactory{
		snake.RoomType: snake.New,
	} {
		if err := RegisterRoomType(Type, Factory); err != nil {
			panic(err)
		}
	}
}
