package snake

import (
	"Wave/internal/application/snake/core"
	"time"
)

// anaconda's snack X)
type food struct {
	*core.Object      // base object
	letter       rune // object letter
}

const typeFood = "food"

func newFood(letter rune, world *core.World, position core.Vec2i) *food {
	f := &food{
		Object: core.NewObject(typeFood),
		letter: letter,
	}
	world.AddObject(f)
	f.SetPos(position)
	return f
}

func (p *food) SetLifetime(left time.Duration) {
	go func() {
		time.Sleep(left)
		p.Destroy()
	}()
}

func (p *food) GetLetter() rune { return p.letter }
