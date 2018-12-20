package snake

import (
	"Wave/internal/application/snake/core"
	"time"
)

type item struct {
	core.Object // base object
}

func newItem(Type core.ObjectType) *item {
	return &item{*core.NewObject(Type)}
}

func (i *item) SetLifetime(left time.Duration) {
	go func() {
		time.Sleep(left)
		i.Destroy()
	}()
}
