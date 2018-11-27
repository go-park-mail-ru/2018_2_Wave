package snake

import (
	"Wave/application/room"
	"Wave/application/snake/core"
	"time"
)

type game struct {
	user2snake map[room.IUser]*snake
	world      *core.World
	walls      *walls
}

func newGame(worldSize core.Vec2i) *game {
	g := &game{
		user2snake: map[room.IUser]*snake{},
		world:      core.NewWorld(worldSize),
	}
	g.walls = newWalls(g.world)
	return g
}

// ----------------|

func (g *game) Tick(dt time.Duration) {
	g.world.Tick(dt)
}

// ----------------| controller interface

// get user bound snake
func (g *game) GetSnake(u room.IUser) (*snake, error) {
	if u == nil {
		return nil, room.ErrorNil
	}
	if snake, ok := g.user2snake[u]; ok {
		return snake, nil
	}
	return nil, room.ErrorNotExists
}

// create a new snake and place it into the world
func (g *game) CreateSnake(u room.IUser, length int) (*snake, error) {
	if u == nil {
		return nil, room.ErrorNil
	}
	if _, ok := g.user2snake[u]; ok {
		return nil, room.ErrorAlreadyExists
	}
	var ( // create a snake object and find a spwn area
		poss, dir = g.world.FindGap(length)
		snake     = newSnake(g.world, poss, dir)
	)
	g.user2snake[u] = snake
	snake.onDestoyed = func() {
		delete(g.user2snake, u)
	}
	return snake, nil
}

// delete a snake associated with the user from the world
func (g *game) DeleteSnake(u room.IUser) error {
	if u == nil {
		return room.ErrorNil
	}
	if s, ok := g.user2snake[u]; ok {
		s.destroy()
		delete(g.user2snake, u)
		return nil
	}
	return room.ErrorNotExists
}

func (g *game) GetGameInfo() *gameInfo {
	gi := &gameInfo{SceneSize: g.world.GetWorldInfo().SceneSize}
	// snakes
	for u, s := range g.user2snake {
		si := snakeInfo{UID: u.GetID()}
		for _, bn := range s.body {
			si.Snake = append(si.Snake, objectInfo{
				Letter:   bn.letter,
				Position: bn.GetPos(),
			})
		}
		gi.Snakes = append(gi.Snakes, si)
	}
	// food
	for _, o := range g.world.GetObjects() {
		if i, ok := o.(*food); ok {
			gi.Food = append(gi.Food, objectInfo{
				Letter:   i.GetLetter(),
				Position: i.GetPos(),
			})
		}
	}
	// walls
	for _, w := range g.walls.blocks {
		gi.Walls = append(gi.Walls, w.GetPos())
	}
	return gi
}
