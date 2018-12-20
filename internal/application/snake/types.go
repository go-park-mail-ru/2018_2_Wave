package snake

import (
	"Wave/internal/application/room"
	"Wave/internal/application/snake/core"
)

type objectInfo struct {
	Letter   rune       `json:"letter"`
	Position core.Vec2i `json:"position"`
}

type snakeInfo struct {
	UserToken room.UserID  `json:"user_token"`
	Score     int          `json:"score"`
	Serial    int64        `json:"user_serial"`
	Snake     []objectInfo `json:"body"`
}

type gameInfo struct {
	SceneSize core.Vec2i   `json:"scene_size"`
	Snakes    []snakeInfo  `json:"snakes"`
	Food      []objectInfo `json:"food"`
	Walls     []core.Vec2i `json:"walls"`
	Boosters  []core.Vec2i `json:"boosters"`
}
