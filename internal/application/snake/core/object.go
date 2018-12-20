package core

import (
	"sync/atomic"
)

// ----------------|

type ObjectType string

// IObject - base world object interface
type IObject interface {
	GetID() uint64       // get object unique id
	GetType() ObjectType // get object type literal
	GetPos() Vec2i       // get objetc position
	SetPos(Vec2i)        //set object position
	GetWorld() *World    // get object world
	OnColided(IObject)   // on colided callback
	Destroy()            // remove the object from it's world and 'destroy'

	setWorld(w *World) // set world variable
	setPos(pos Vec2i)  // set position variable
}

// ----------------|

var objectIDCounter uint64

// Object- base game object
type Object struct {
	id       uint64     // object unique id
	position Vec2i      // object position
	world    *World     // object world
	_type    ObjectType // object type literal
}

func NewObject(_type ObjectType) *Object {
	return &Object{
		_type: _type,
		id:    atomic.AddUint64(&objectIDCounter, 1),
	}
}

func (o *Object) GetID() uint64       { return o.id }
func (o *Object) GetType() ObjectType { return o._type }
func (o *Object) GetPos() Vec2i       { return o.position }
func (o *Object) GetWorld() *World    { return o.world }
func (o *Object) OnColided(IObject)   {}

func (o *Object) SetPos(expected Vec2i) {
	if o.world != nil {
		o.world.onObjectMove(o, expected)
	}
}

func (o *Object) setWorld(w *World) { o.world = w }
func (o *Object) setPos(pos Vec2i)  { o.position = pos }

func (o *Object) Destroy() {
	if o.world != nil {
		o.world.RemoveObject(o)
	}
}
