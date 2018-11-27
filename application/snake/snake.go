package snake

import (
	"time"

	"Wave/application/snake/core"
)

// ----------------| snake node

// snake body portion
type snakeNode struct {
	*core.Object                // base object
	letter       rune           // node letter
	snake        *snake         // parent snake
	direction    core.Direction // next node direction
	bHead        bool           // head node
}

const typeSnakeNode = "snake_node"

func newSnakeNode(letter rune, s *snake, position core.Vec2i) *snakeNode {
	n := &snakeNode{
		Object: core.NewObject(typeSnakeNode),
		letter: letter,
		snake:  s,
	}
	n.SetWorld(s.world)
	n.SetPos(position)
	return n
}

func (s *snakeNode) OnColided(o core.IObject) {
	if p, ok := o.(iItem); ok && p != nil {
		s.snake.pushBack(p.GetLetter())
		o.Destroy()
	}
}

// ----------------| snake

// snake representation
type snake struct {
	world    *core.World
	body     []*snakeNode   // body elements
	movement core.Direction // next step direction

	ticker *time.Ticker
	cancel chan interface{}
}

func newSnake(w *core.World, points []core.Vec2i, direction core.Direction) *snake {
	s := &snake{
		world:    w,
		ticker:   time.NewTicker(200 * time.Millisecond),
		cancel:   make(chan interface{}, 1),
		movement: direction,
	}
	for i := len(points) - 1; i >= 0; i-- {
		s.setHead('h', direction, points[i])
	}
	go s.tick()

	return s
}

func (s *snake) destroy() {
	s.cancel <- ""
	for _, elem := range s.body {
		elem.Destroy()
	}
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
		delta         = s.movement.GetDelta()
		nextPosition  = s.body[0].GetPos().Sum(delta)
		nextDirection = s.movement
		tmpPosition   = core.Vec2i{}
	)
	for i := range s.body {
		nextPosition, tmpPosition = s.body[i].GetPos(), nextPosition
		nextDirection, s.body[i].direction = s.body[i].direction, nextDirection
		s.body[i].SetPos(tmpPosition)
	}
}

func (s *snake) pushBack(letter rune) {
	if len(s.body) > 0 {
		var (
			curTail     = s.getTail()
			direction   = curTail.direction
			delta       = direction.GetDelta().Mult(-1)
			newPosition = curTail.GetPos().Sum(delta)
			newTail     = newSnakeNode(letter, s, newPosition)
		)
		newTail.direction = direction

		s.body = append(s.body, newTail)
	}
}

func (s *snake) setHead(letter rune, direction core.Direction, position core.Vec2i) {
	newHead := newSnakeNode(letter, s, position)
	newHead.direction = direction
	s.body = append([]*snakeNode{newHead}, s.body...)
}

func (s *snake) getTail() *snakeNode {
	if len(s.body) > 0 {
		return s.body[len(s.body)-1]
	}
	return nil
}
