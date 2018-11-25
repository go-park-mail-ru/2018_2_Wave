package snake

import (
	"time"
)

// ----------------| snake node

// snake body portion
type snakeNode struct {
	*object             // base object
	letter    rune      // node letter
	snake     *snake    // parent snake
	direction direction // next node direction
	bHead     bool      // head node
}

const typeSnakeNode = "snake_node"

func newSnakeNode(letter rune, s *snake) *snakeNode {
	n := &snakeNode{
		object: newObject(typeSnakeNode),
		letter: letter,
		snake:  s,
	}
	n.SetWorld(s.world)
	return n
}

func (s *snakeNode) OnColided(o iObject) {
	if p, ok := o.(iItem); ok && p != nil {
		s.snake.pushBack(p.GetLetter())
		o.Destroy()
	}
}

// ----------------| snake

// snake representation
type snake struct {
	world    *world
	body     []*snakeNode // body elements
	movement direction    // next step direction

	ticker *time.Ticker
	cancel chan interface{}
}

func newSnake(w *world, points []vec2i, direction direction) *snake {
	s := &snake{
		world:    w,
		ticker:   time.NewTicker(300 * time.Microsecond),
		movement: direction,
	}
	for i := len(points); i >= 0; i-- {
		s.setHead('h', direction, points[i])
	}
	go s.tick()

	return s
}

func (s *snake) destroy() {
	s.cancel <- ""
}

func (s *snake) tick() {
	for { //
		select {
		case <-s.ticker.C:
			s.moveNext()
		case <-s.cancel:
			return
		}
	}
}

func (s *snake) moveNext() {
	var (
		delta         = s.movement.getDelta()
		nextPosition  = s.body[0].Position.Sum(delta)
		nextDirection = s.movement
	)
	for i := range s.body {
		nextPosition, s.body[i].Position = s.body[i].Position, nextPosition
		nextDirection, s.body[i].direction = s.body[i].direction, nextDirection
	}
}

func (s *snake) pushBack(letter rune) {
	if len(s.body) > 0 {
		var (
			curTail     = s.getTail()
			newTail     = newSnakeNode(letter, s)
			direction   = curTail.direction
			delta       = direction.getDelta().Mult(-1)
			newPosition = curTail.Position.Sum(delta)
		)
		newTail.Position = newPosition
		newTail.direction = direction

		s.body = append(s.body, newTail)
	}
}

func (s *snake) setHead(letter rune, direction direction, position vec2i) {
	newHead := newSnakeNode(letter, s)
	newHead.Position = position
	newHead.direction = direction
	s.body = append([]*snakeNode{newHead}, s.body...)
}

func (s *snake) getTail() *snakeNode {
	if len(s.body) > 0 {
		return s.body[len(s.body)-1]
	}
	return nil
}
