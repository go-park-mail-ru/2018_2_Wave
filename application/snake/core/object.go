package core

// ----------------|

type objectType string

// IObject - base world object interface
type IObject interface {
	GetType() objectType // get object type literal
	GetPos() Vec2i       // get objetc position
	SetPos(Vec2i)        //set object position
	GetWorld() *World    // get object world
	SetWorld(w *World)   // set objetc world
	OnColided(IObject)   // on colided callback
	Destroy()            // remove the object from it's world and 'destroy'
}

// ----------------|

// Object- base game object
type Object struct {
	position Vec2i      // object position
	world    *World     // object world
	_type    objectType // object type literal
}

func NewObject(_type objectType) *Object {
	return &Object{
		_type: _type,
	}
}

func (o *Object) GetType() objectType { return o._type }
func (o *Object) GetPos() Vec2i       { return o.position }
func (o *Object) GetWorld() *World    { return o.world }
func (o *Object) OnColided(IObject)   {}

func (o *Object) SetPos(expected Vec2i) {
	if o.world != nil {
		if next, err := o.world.onObjectMove(o, expected); err == nil {
			o.position = next
		}
	}
}

func (o *Object) SetWorld(w *World) {
	if o.world == w {
		return
	}
	if o.world != nil {
		o.world.RemoveObject(o)
	}
	o.world = w
	if o.world != nil {
		o.world.AddObject(o)
	}
}

func (o *Object) Destroy() {
	if o.world != nil {
		o.world.RemoveObject(o)
	}
}
