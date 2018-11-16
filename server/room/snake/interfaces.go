package snake

type objectType string

// basic world object
type iObject interface {
	GetType() objectType // get object type literal
	GetPos() vec2i       // get objetc position
	GetWorld() *world    // get object world
	SetWorld(w *world)   // set objetc world
	OnColided(iObject)   // on colided callback
	Destroy()            // remove the object from it's world and 'destroy'
}

// pickable object
type iItem interface {
	iObject
	GetLetter() rune
}
