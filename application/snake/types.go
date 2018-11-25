package snake

import (
	"Wave/server/room"
)

// ----------------| vector

type vec2i struct {
	X, Y int
}

func (v vec2i) Sum(o vec2i) vec2i {
	return vec2i{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

func (v vec2i) Mult(o int) vec2i {
	return vec2i{
		X: v.X * o,
		Y: v.Y * o,
	}
}

// ----------------| direction

// direction
type direction int

const ( // dirrection enum
	noDirection direction = 0
	up          direction = 1 << iota
	down
	left
	right
)

func (d direction) have(o direction) bool {
	return d&o != 0
}

func (d direction) getDelta() (res vec2i) {
	if d == noDirection {
		return vec2i{}
	}
	if d.have(up) {
		res.Y++
	}
	if d.have(down) {
		res.Y--
	}
	if d.have(left) {
		res.X--
	}
	if d.have(right) {
		res.X++
	}
	return res
}

// ----------------| scene size

type sceneSize struct {
	X, Y int
}

type sceneItemInfo struct {
	Letter   rune
	Position vec2i
}

type playerInfo struct {
	UID   room.UserID
	Snake []sceneItemInfo
}

type worldInfo struct {
	sceneSize sceneSize
}

type sceneInfo struct {
	Playes []playerInfo
	Items  []sceneItemInfo
}

type gameInfo struct {
	SceneSize sceneSize
	sceneInfo
}
