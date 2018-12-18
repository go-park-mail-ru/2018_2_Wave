package snake

import (
	"Wave/internal/application/snake/core"
	"time"
)

type item struct {
	core.Object // base object
}

func newItem(Type core.ObjectType, world *core.World, position core.Vec2i) *item {
	i := &item{*core.NewObject(Type)}
	world.AddObject(i)
	i.SetPos(position)
	return i
}

func (i *item) SetLifetime(left time.Duration) {
	go func() {
		time.Sleep(left)
		i.Destroy()
	}()
}
