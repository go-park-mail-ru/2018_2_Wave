package proto

type IManager interface {
	CreateLobby(RoomToken, IUser, RoomType) (IRoom, error)
	RemoveLobby(RoomToken, IUser) error // if IUser==nil - force delete

	IRoom
}
