package snake

import (
	"time"

	"Wave/internal/application/snake/core"
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
	s.world.AddObject(n)
	n.SetPos(position)
	return n
}

func (s *snakeNode) OnColided(o core.IObject) {
	if f, ok := o.(*food); ok {
		s.snake.pushBack(f.GetLetter())
		s.snake.score++
		f.Destroy()
	}
	if _, ok := o.(*snakeNode); ok && s.bHead {
		s.snake.destroy()
	}
	if _, ok := o.(*wall); ok {
		s.snake.destroy()
	}
	if b, ok := o.(*booster); ok {
		s.snake.tickFactor *= b.Factor
		go func() {
			time.Sleep(b.Duration)
			s.snake.tickFactor /= b.Factor
		}()
		b.Destroy()
	}
}

// ----------------| snake

// snake representation
type snake struct {
	world    *core.World    // game world
	body     []*snakeNode   // body elements
	movement core.Direction // next step direction
	score    int            // game score

	ticker     core.Ticker   // movement ticker
	baseTick   time.Duration // base tick time
	tickFactor float64       // length factor

	onDestoyed func() // to remove the snake from the game
}

func newSnake(w *core.World, points []core.Vec2i, direction core.Direction) *snake {
	s := &snake{
		world:      w,
		movement:   direction,
		baseTick:   100 * time.Millisecond,
		tickFactor: 0.8,
	}
	s.ticker = core.MakeTicker(s.moveNext, s.baseTick)
	l := 'a'
	for i := range points {
		s.setHead(l, direction, points[i])
		l++
	}
	return s
}

func (s *snake) Tick(dt time.Duration) {
	s.ticker.Tick(dt)
}

func (s *snake) SetDirection(d core.Direction) {
	curr := s.movement.GetDelta()
	next := d.GetDelta()
	if curr.Diff(next).IsZero() {
		return
	}
	s.movement = d
}

func (s *snake) destroy() {
	for _, elem := range s.body {
		elem.Destroy()
	}
	if s.onDestoyed != nil {
		s.onDestoyed()
	}
}

func (s *snake) moveNext(dt time.Duration) {
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
		s.onLengthChanged()
	}
}

func (s *snake) setHead(letter rune, direction core.Direction, position core.Vec2i) {
	newHead := newSnakeNode(letter, s, position)
	newHead.bHead = true
	newHead.direction = direction
	if len(s.body) > 0 {
		s.body[0].bHead = true
	}
	s.body = append([]*snakeNode{newHead}, s.body...)
	s.onLengthChanged()
}

func (s *snake) getTail() *snakeNode {
	if len(s.body) > 0 {
		return s.body[len(s.body)-1]
	}
	return nil
}

func (s *snake) onLengthChanged() {
	if len(s.body) > 0 {
		factor := s.tickFactor * float64(len(s.body))
		time := s.baseTick * time.Duration(factor)
		s.ticker.SetTickTime(time)
	}
}
