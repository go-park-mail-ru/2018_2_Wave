package snake

import (
	"Wave/internal/application/snake/core"
	"time"
)

type booster struct {
	item
	Factor   float64
	Duration time.Duration
}

const typeBooster = "booster"

func newBooster(world *core.World, position core.Vec2i) *booster {
	return &booster{
		item:     *newItem(typeBooster, world, position),
		Factor:   1.5,
		Duration: 3,
	}
}
