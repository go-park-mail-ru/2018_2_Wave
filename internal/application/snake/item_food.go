package snake

import (
	"Wave/internal/application/snake/core"
)

// anaconda's snack X)
type food struct {
	item        // base object
	letter rune // object letter
}

const typeFood = "food"

func newFood(letter rune, world *core.World, position core.Vec2i) *food {
	f := &food{
		item:   *newItem(typeFood),
		letter: letter,
	}
	world.AddObject(f)
	f.SetPos(position)
	return f
}

func (p *food) GetLetter() rune { return p.letter }
