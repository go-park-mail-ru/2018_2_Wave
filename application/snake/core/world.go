package core

import (
	"Wave/application/room"
	"time"
)

// World - game world
type World struct {
	scene *scene

	info WorldInfo
}

// NewWorld - constructor
func NewWorld(size Vec2i) *World {
	return &World{
		scene: newScene(size),
		info: WorldInfo{
			SceneSize: size,
		},
	}
}

// ----------------| tick and lifecycle

func (w *World) Tick(dt time.Duration) {
	w.scene.Tick()
	// w.scene.PrintDebug()
}

// ----------------|

func (w *World) GetWorldInfo() WorldInfo { return w.info }
func (w *World) GetObjects() []IObject {
	if w.scene != nil {
		return w.scene.objects
	}
	return nil
}

// ----------------| object handling

// AddObject assigns the object to the scen but not places in.
func (w *World) AddObject(o IObject) error {
	if w.scene != nil || o == nil {
		if o.GetWorld() == w {
			return nil
		}

		if err := w.scene.AddObject(o); err != nil {
			return err
		}

		if o.GetWorld() != nil {
			o.GetWorld().RemoveObject(o)
		}
		o.setWorld(w)
		return nil
	}
	return room.ErrorNil
}

// RemoveObject removes the object from the world.
func (w *World) RemoveObject(o IObject) error {
	if w.scene != nil {
		return w.scene.RemoveObject(o)
	}
	return room.ErrorNil
}

// ----------------| scene functions

// FindGap of @length
func (w *World) FindGap(length int) (res []Vec2i, dir Direction, err error) {
	if w.scene != nil {
		return w.scene.FindGap(length)
	}
	// TODO:: log
	return nil, NoDirection, nil
}

func (w *World) onObjectMove(o IObject, expectedPosition Vec2i) (err error) {
	if w.scene != nil {
		return w.scene.onObjectMove(o, expectedPosition)
	}
	// TODO:: log
	return room.ErrorNil
}
