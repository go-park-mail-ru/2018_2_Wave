package snake

import (
	"Wave/application/snake/core"
	"math/rand"
)

// ----------------| wall

const wallType = "wall"

type wall struct {
	*core.Object // base object
}

func newWall(world *core.World, position core.Vec2i) *wall {
	w := &wall{
		Object: core.NewObject(wallType),
	}
	world.AddObject(w) // TODO:: change order. HINT: look at @SetWorld function problems
	w.SetPos(position)
	return w
}

// ----------------| wall

type walls struct {
	world  *core.World
	blocks []*wall
}

func newWalls(world *core.World) *walls {
	ws := &walls{
		world: world,
	}
	const skeep = 0.1

	size := world.GetWorldInfo().SceneSize
	for i := 0; i < size.X; i++ {
		if rand.Float32() < skeep {
			i += 9
			continue
		}
		var (
			upWall   = newWall(world, core.Vec2i{X: i, Y: size.Y - 1})
			downWall = newWall(world, core.Vec2i{X: i, Y: 0})
		)
		ws.blocks = append(ws.blocks, upWall, downWall)
	}

	for i := 1; i < size.Y-1; i++ {
		if rand.Float32() < skeep {
			i += 5
			continue
		}
		var (
			lWall = newWall(world, core.Vec2i{X: 0, Y: i})
			rWall = newWall(world, core.Vec2i{X: size.X - 1, Y: i})
		)
		ws.blocks = append(ws.blocks, lWall, rWall)
	}

	wallLen := 5
	numWalls := 6
	for i := 0; i < numWalls; i++ {
		points, _ := world.FindGap(wallLen, func() core.Direction {
			if rand.Float32() < 0.5 {
				return core.Right | core.Up
			}
			return core.Right | core.Down
		}())
		for _, point := range points {
			ws.blocks = append(ws.blocks, newWall(world, point))
		}
	}

	return ws
}

func (ws *walls) Destroy() {
	for _, w := range ws.blocks {
		w.Destroy()
	}
}
