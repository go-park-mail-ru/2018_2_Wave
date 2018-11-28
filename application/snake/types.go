package snake

import (
	"Wave/application/room"
	"Wave/application/snake/core"
)

type objectInfo struct {
	Letter   rune       `json:"letter"`
	Position core.Vec2i `json:"position"`
}

type snakeInfo struct {
	UID   room.UserID  `json:"user_id"`
	Snake []objectInfo `json:"body"`
}

type gameInfo struct {
	SceneSize core.Vec2i   `json:"scene_size"`
	Snakes    []snakeInfo  `json:"snakes"`
	Food      []objectInfo `json:"food"`
	Walls     []core.Vec2i `json:"walls"`
}
