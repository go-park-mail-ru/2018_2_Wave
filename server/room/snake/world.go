package snake

import (
	"Wave/server/room"
	"math/rand"
	"time"
)

type worldInfo struct {
	sceneSize sceneSize
}

// game world
type world struct {
	user2snake map[room.IUser]*snake
	scene      [][]iObject
	objects    []iObject

	info worldInfo
}

// constructor
func newWorld(size sceneSize) *world {
	w := &world{
		user2snake: map[room.IUser]*snake{},
		scene:      make([][]iObject, size.X),
		info: worldInfo{
			sceneSize: size,
		},
	}
	for i := range w.scene {
		w.scene[i] = make([]iObject, size.Y)
	}
	return w
}

// ----------------| tick and lifecycle

// tick the world
func (w *world) Tick(dt time.Duration) {
	// NOTE: snakes tick itself
	w.cleanScene()
	// TODO: create food
	w.fillInscene(w.objects)
}

// ----------------| object handling

// assign the object to the world
func (w *world) AddObject(o iObject) error {
	_, err := w.getObjectIdx(o)
	if err == nil {
		return room.ErrorAlreadyExists
	}

	w.objects = append(w.objects, o)
	return nil
}

// remove the object from the forld
func (w *world) RemoveObject(o iObject) error {
	idx, err := w.getObjectIdx(o)
	if err != nil {
		return room.ErrorNotExists
	}

	w.objects = append(w.objects[:idx], w.objects[idx+1:]...)
	return nil
}

func (w *world) getObjectIdx(o iObject) (int, error) {
	for i, ob := range w.objects {
		if ob == o {
			return i, nil
		}
	}
	return 0, room.ErrorNotExists
}

// ----------------| controller interface

// get user bound snake
func (w *world) GetSnake(u room.IUser) (*snake, error) {
	if u == nil {
		return nil, room.ErrorNil
	}
	if snake, ok := w.user2snake[u]; ok {
		return snake, nil
	}
	return nil, room.ErrorNotExists
}

// create a new snake and place it into the world
func (w *world) CreateSnake(u room.IUser, length int) (*snake, error) {
	if u == nil {
		return nil, room.ErrorNil
	}
	if _, ok := w.user2snake[u]; ok {
		return nil, room.ErrorAlreadyExists
	}
	var ( // create a snake object and find a spwn area
		poss, dir = w.findSnakeCreationLocation(length)
		snake     = newSnake(w, poss, dir)
	)
	w.user2snake[u] = snake
	return snake, nil
}

// delete a snake associated with the user from the world
func (w *world) DeleteSnake(u room.IUser) error {
	if u == nil {
		return room.ErrorNil
	}
	if _, ok := w.user2snake[u]; !ok {
		return room.ErrorNotExists
	}
	delete(w.user2snake, u)
	return nil
}

func (w *world) GetGameInfo() *gameInfo {
	gf := &gameInfo{
		SceneSize: w.info.sceneSize,
	}

	for u, s := range w.user2snake {
		pf := playerInfo{
			UID: u.GetID(),
		}
		for _, bn := range s.body {
			pf.Snake = append(pf.Snake, sceneItemInfo{
				Letter:   bn.letter,
				Position: bn.GetPos(),
			})
		}
		gf.Playes = append(gf.Playes, pf)
	}
	for _, o := range w.objects {
		if i, ok := o.(iItem); ok {
			gf.Items = append(gf.Items, sceneItemInfo{
				Letter:   i.GetLetter(),
				Position: i.GetPos(),
			})
		}
	}
	return gf
}

// ----------------| find functions

// find positions for line nodes
func (w *world) findSnakeCreationLocation(length int) (res []vec2i, dir direction) {
	position := vec2i{
		X: rand.Intn(w.info.sceneSize.X),
		Y: rand.Intn(w.info.sceneSize.Y),
	}
	for i := 0; i < length; i++ {
		if w.getSceneAtPosition(position) != nil {
			return w.findSnakeCreationLocation(length)
		}
		res = append(res, position)
		position.X++
	}
	return res, right
}

// ----------------| scene functions

func (w *world) getSceneAtPosition(position vec2i) iObject {
	return w.scene[position.X][position.Y]
}

func (w *world) cleanScene() {
	for x, col := range w.scene {
		for y := range col {
			w.scene[x][y] = nil
		}
	}
}

func (w *world) fillInscene(objs []iObject) {
	max := w.info.sceneSize
	for _, o := range objs {
		if o == nil {
			continue
		}
		pos := o.GetPos()
		if pos.X > max.X || pos.Y > max.Y {
			// TODO:: handle the fuck
		}

		other := w.scene[pos.X][pos.Y]
		w.scene[pos.X][pos.Y] = o
		if other == nil {
			continue
		}

		other.OnColided(o)
		o.OnColided(other)
	}
}
