package snake

// base game object
type object struct {
	Position vec2i      // object position
	World    *world     // object world
	Type     objectType // object type literal
}

func newObject(Type objectType) *object {
	return &object{
		Type: Type,
	}
}

func (o *object) GetType() objectType { return o.Type }
func (o *object) GetPos() vec2i       { return o.Position }
func (o *object) GetWorld() *world    { return o.World }
func (o *object) OnColided(iObject)   {}

func (o *object) SetWorld(w *world) {
	if o.World == w {
		return
	}
	if o.World != nil {
		o.World.RemoveObject(o)
	}
	o.World = w
	if o.World != nil {
		o.World.AddObject(o)
	}
}

func (o *object) Destroy() {
	if o.World != nil {
		o.World.RemoveObject(o)
	}
}
