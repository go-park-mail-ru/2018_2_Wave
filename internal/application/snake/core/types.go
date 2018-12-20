package core

import (
	"time"
)

// ----------------| vector

type Vec2i struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (v Vec2i) Sum(o Vec2i) Vec2i {
	return Vec2i{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

func (v Vec2i) Diff(o Vec2i) Vec2i {
	return Vec2i{
		X: v.X - o.X,
		Y: v.Y - o.Y,
	}
}

func (v Vec2i) Mult(o int) Vec2i {
	return Vec2i{
		X: v.X * o,
		Y: v.Y * o,
	}
}

func (v Vec2i) IsZero() bool {
	return v.X == 0 && v.Y == 0
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

type Ticker struct {
	clb         func(time.Duration)
	accumulated time.Duration
	tickTime    time.Duration
}

func MakeTicker(clb func(time.Duration), tickTime time.Duration) Ticker {
	return Ticker{
		clb:      clb,
		tickTime: tickTime,
	}
}

func (t *Ticker) Tick(dt time.Duration) {
	t.accumulated += dt
	for {
		if t.accumulated > t.tickTime {
			if t.clb != nil {
				t.clb(t.tickTime)
			}
			t.accumulated -= t.tickTime
		} else {
			break
		}
	}
}

func (t *Ticker) SetTickTime(tickTime time.Duration) {
	t.tickTime = tickTime
}

// ----------------|

type WorldInfo struct {
	SceneSize Vec2i
}
