package core

import (
	"Wave/application/room"
	"time"
)

// World - game world
type World struct {
	scene   *scene
	objects []IObject

	info worldInfo
}

// NewWorld - constructor
func NewWorld(size Vec2i) *World {
	return &World{
		scene: newScene(size),
		info: worldInfo{
			sceneSize: size,
		},
	}
}

// ----------------| tick and lifecycle

// tick the world
func (w *World) Tick(dt time.Duration) {

}

// ----------------| object handling

// assign the object to the world
func (w *World) AddObject(o IObject) error {
	_, err := w.getObjectIdx(o)
	if err == nil {
		return room.ErrorAlreadyExists
	}

	w.objects = append(w.objects, o)
	return nil
}

// remove the object from the forld
func (w *World) RemoveObject(o IObject) error {
	idx, err := w.getObjectIdx(o)
	if err != nil {
		return room.ErrorNotExists
	}

	w.objects = append(w.objects[:idx], w.objects[idx+1:]...)
	return nil
}

func (w *World) getObjectIdx(o IObject) (int, error) {
	for i, ob := range w.objects {
		if ob == o {
			return i, nil
		}
	}
	return 0, room.ErrorNotExists
}

// ----------------| scene functions

// FindGap of @length
func (w *World) FindGap(length int) (res []Vec2i, dir Direction) {
	if w.scene != nil {
		return w.scene.FindGap(length)
	}
	// TODO:: log
	return nil, NoDirection
}

func (w *World) onObjectMove(o IObject, expectedPosition Vec2i) (nextPosition Vec2i, err error) {
	if w.scene != nil {
		return w.scene.onObjectMove(o, expectedPosition)
	}
	// TODO:: log
	return Vec2i{}, room.ErrorNil
}
