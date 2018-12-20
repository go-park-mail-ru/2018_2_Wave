package core

// import (
// 	"sync/atomic"
// )

// type IActor interface {
// 	GetID() uint64
// 	Destory()
// 	GetWorld() *World
// 	setWorld(*World)
// }

// // ----------------|

// var actorIDCounter uint64

// type Actor struct {
// 	id      uint64
// 	world   *World
// 	objects []IObject
// }

// func NewActor() IActor {
// 	return &Actor{
// 		id: atomic.AddUint64(&actorIDCounter, 1),
// 	}
// }

// func (a *Actor) GetID() uint64    { return a.id }
// func (a *Actor) GetWorld() *World { return a.world }

// func (a *Actor) Destory() {

// }

// func (a *Actor) setWorld(*World) {

// }
