package snake

import (
	"Wave/application/room"
	"Wave/application/snake/core"
	"fmt"
	"time"
)

type game struct {
	user2snake map[room.IUser]*snake
	world      *core.World
}

func newGame(worldSize core.Vec2i) *game {
	return &game{
		user2snake: map[room.IUser]*snake{},
		world:      core.NewWorld(worldSize),
	}
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
	fmt.Printf("%v\n", poss) // TODO:: remove
	g.user2snake[u] = snake
	return snake, nil
}

// delete a snake associated with the user from the world
func (g *game) DeleteSnake(u room.IUser) error {
	if u == nil {
		return room.ErrorNil
	}
	if _, ok := g.user2snake[u]; !ok {
		return room.ErrorNotExists
	}
	delete(g.user2snake, u)
	return nil
}

func (g *game) GetSceneInfo() *sceneInfo {
	si := &sceneInfo{}
	for u, s := range g.user2snake {
		pf := playerInfo{
			UID: u.GetID(),
		}
		for _, bn := range s.body {
			pf.Snake = append(pf.Snake, sceneItemInfo{
				Letter:   bn.letter,
				Position: bn.GetPos(),
			})
		}
		si.Playes = append(si.Playes, pf)
	}
	// for _, o := range g.objects {
	// 	if i, ok := o.(iItem); ok {
	// 		si.Items = append(si.Items, sceneItemInfo{
	// 			Letter:   i.GetLetter(),
	// 			Position: i.GetPos(),
	// 		})
	// 	}
	// }
	return si
}

func (g *game) GetGameInfo() *gameInfo {
	return &gameInfo{
		// SceneSize: w.info.sceneSize,
		// sceneInfo: *w.GetSceneInfo(),
	}
}
