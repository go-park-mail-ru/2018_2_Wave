package globalRoom

import (
	"Wave/server/room"
)

// type2Factory - singletone map converting room type to the room factory
// All room types MUST assign self into the map with use of @AddTypeFactory()
var type2Factory = map[room.RoomType]room.RoomFactory{}

func AddTypeFactory(roomType room.RoomType, factory room.RoomFactory) error {
	if _, ok := type2Factory[roomType]; ok {
		return room.ErrorAlreadyExists
	}
	type2Factory[roomType] = factory
	return nil
}
