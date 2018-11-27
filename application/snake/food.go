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

func newFood(letter rune) *food {
	return &food{
		Object: core.NewObject(typeFood),
		letter: letter,
	}
}

func (p *food) GetLetter() rune { return p.letter }
