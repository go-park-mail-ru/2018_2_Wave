package snake

// letter could be mealed by a snake
type item struct {
	*object      // base object
	letter  rune // object letter
}

const typePickup = "item"

func newPickup(letter rune) *item {
	return &item{
		object: newObject(typePickup),
		letter: letter,
	}
}

func (p *item) GetLetter() rune { return p.letter }
