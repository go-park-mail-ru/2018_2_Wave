package core

// ----------------| vector

type Vec2i struct {
	X, Y int
}

func (v Vec2i) Sum(o Vec2i) Vec2i {
	return Vec2i{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

func (v Vec2i) Mult(o int) Vec2i {
	return Vec2i{
		X: v.X * o,
		Y: v.Y * o,
	}
}

// ----------------| direction

// direction
type Direction int

const ( // dirrection enum
	NoDirection Direction = 0
	Up          Direction = 1 << iota
	Down
	Left
	Right
)

func (d Direction) Is(o Direction) bool {
	return d&o != 0
}

func (d Direction) GetDelta() (res Vec2i) {
	if d == NoDirection {
		return Vec2i{}
	}
	if d.Is(Up) {
		res.Y++
	}
	if d.Is(Down) {
		res.Y--
	}
	if d.Is(Left) {
		res.X--
	}
	if d.Is(Right) {
		res.X++
	}
	return res
}

// ----------------|

type worldInfo struct {
	sceneSize Vec2i
}
