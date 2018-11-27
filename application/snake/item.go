package snake

import (
	"Wave/application/snake/core"
)

// ----------------|

// pickable object
type iItem interface {
	core.IObject
	GetLetter() rune
}

// ----------------|

// letter could be mealed by a snake
type item struct {
	*core.Object      // base object
	letter       rune // object letter
}

const typePickup = "item"

func newPickup(letter rune) *item {
	return &item{
		Object: core.NewObject(typePickup),
		letter: letter,
	}
}

func (p *item) GetLetter() rune { return p.letter }
