package snake

import (
	"Wave/application/snake/core"
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

func (p *food) GetLetter() rune { return p.letter }
