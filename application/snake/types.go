package snake

import (
	"Wave/application/room"
	"Wave/application/snake/core"
)

type sceneItemInfo struct {
	Letter   rune
	Position core.Vec2i
}

type playerInfo struct {
	UID   room.UserID
	Snake []sceneItemInfo
}

type sceneInfo struct {
	Playes []playerInfo
	Items  []sceneItemInfo
}

type gameInfo struct {
	SceneSize core.Vec2i
	sceneInfo
}
