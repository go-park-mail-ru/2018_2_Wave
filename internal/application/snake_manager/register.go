package manager

import (
	"Wave/internal/application/proto"
	"Wave/internal/application/snake"
)

// type2Factory - singletone map converting room type to the room factory
// To add a new factory user @RegisterRoomType function
var type2Factory = map[proto.RoomType]proto.RoomFactory{}

// RegisterRoomType factory
func RegisterRoomType(roomType proto.RoomType, factory proto.RoomFactory) error {
	if IsRegisteredType(roomType) {
		return proto.ErrorAlreadyExists
	}
	type2Factory[roomType] = factory
	return nil
}

// IsRegisteredType - weather the type factory was gegistered
func IsRegisteredType(roomType proto.RoomType) bool {
	_, ok := type2Factory[roomType]
	return ok
}

func init() {
	for Type, Factory := range map[proto.RoomType]proto.RoomFactory{
		snake.RoomType: snake.New,
	} {
		if err := RegisterRoomType(Type, Factory); err != nil {
			panic(err)
		}
	}
}
