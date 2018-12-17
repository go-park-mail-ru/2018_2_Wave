package core

import (
	"Wave/internal/application/room"
	"Wave/internal/logger"

	"math/rand"
	"strconv"
	"sync"
)

type scene struct {
	LG         *logger.Logger
	fields     [][]field
	objects    []IObject
	objectMap  map[uint64]IObject
	size       Vec2i
	collisions []collision
	mu         sync.RWMutex
}

func newScene(size Vec2i) *scene {
	s := &scene{
		fields:    make([][]field, size.X),
		objectMap: make(map[uint64]IObject),
		size:      size,
	}
	for i := range s.fields {
		s.fields[i] = make([]field, size.Y)
	}
	return s
}

// ----------------|

func (s *scene) Tick() {
	for _, cl := range s.collisions {
		cl.collide()
	}
	s.collisions = nil
}

// assign the object th the scene
func (s *scene) AddObject(o IObject) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isPlaced(o) {
		return room.ErrorAlreadyExists
	}
	s.objects = append(s.objects, o)
	s.objectMap[o.GetID()] = o
	return nil
}

func (s *scene) RemoveObject(o IObject) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if o = s.actualiser(o); o == nil {
		return room.ErrorNotExists
	}
	for i, expectant := range s.objects {
		if expectant != o {
			continue
		}
		s.objects = append(s.objects[:i], s.objects[i+1:]...)
		delete(s.objectMap, o.GetID())

		currPosition, err := s.validatePosition(o.GetPos())
		if err != nil {
			return err
		}
		s.at(currPosition).remove(o)

		return nil
	}
	return room.ErrorNotExists
}

func (s *scene) FindGap(length int, dir Direction) (res []Vec2i, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delta := dir.GetDelta()
	iteration := 0
FIND_POSITION:
	{
		if iteration++; iteration > 300 {
			return nil, room.ErrorNotFound
		}

		position := Vec2i{
			X: rand.Intn(s.size.X),
			Y: rand.Intn(s.size.Y),
		}
		for i := 0; i < length; i++ {
			o := s.at(position)
			if o == nil || !o.isEmpty() {
				goto FIND_POSITION
			}
			res = append(res, position)
			position = position.Sum(delta)
		}
		return res, nil
	}
}

func (s *scene) PrintDebug() {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := ""
	for y := s.size.Y - 1; y >= 0; y-- {
		for x := 0; x < s.size.X; x++ {
			f := s.fields[x][y]
			if f.isEmpty() {
				res += "-"
			} else {
				res += strconv.Itoa(len(f))
			}
		}
		res += "\n"
	}
	println(res)
}

// ----------------|

// some functions calls by embedded structs
// but the structs have no acces to overrided
// methods of ther embeders, so we need to upgrade
// the incoming structs to their embedders to call the functions
func (s *scene) actualiser(o IObject) IObject {
	if o == nil {
		return nil
	}
	return s.objectMap[o.GetID()]
}

func (s *scene) onObjectMove(o IObject, expectedPosition Vec2i) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if o = s.actualiser(o); o != nil {
		currPosition, err := s.validatePosition(o.GetPos())
		if err != nil {
			return err
		}

		nextPosition, err := s.validatePosition(expectedPosition)
		if err != nil {
			return err
		}
		o.setPos(nextPosition)

		s.at(currPosition).remove(o)
		cs := s.at(nextPosition).collide(o)
		s.collisions = append(s.collisions, cs...)
		return nil
	}
	return room.ErrorNil
}

func (s *scene) isPlaced(o IObject) bool {
	if o == nil {
		return false
	}
	_, ok := s.objectMap[o.GetID()]
	return ok
}

func (s *scene) validatePosition(expectedPosition Vec2i) (validPosition Vec2i, err error) {
	validPosition.X = expectedPosition.X % s.size.X
	validPosition.Y = expectedPosition.Y % s.size.Y
	if validPosition.X < 0 {
		validPosition.X = s.size.X + validPosition.X
	}
	if validPosition.Y < 0 {
		validPosition.Y = s.size.Y + validPosition.Y
	}
	return validPosition, nil
}

func (s *scene) at(position Vec2i) *field {
	if position.X < s.size.X && position.X >= 0 &&
		position.Y < s.size.Y && position.Y >= 0 {
		return &s.fields[position.X][position.Y]
	}
	return nil
}

// ----------------| field

type field []IObject

func (f *field) remove(o IObject) {
	if o != nil {
		for i, elem := range *f {
			if elem != o {
				continue
			}
			*f = append((*f)[:i], (*f)[i+1:]...)
			return
		}
	}
}

func (f *field) collide(o IObject) (collisions []collision) {
	if o != nil {
		lastState := f.dump()
		*f = append(*f, o)
		for _, elem := range lastState {
			collisions = append(collisions, collision{
				o0: elem,
				o1: o,
			})
		}
	}
	return collisions
}

func (f *field) dump() []IObject {
	dump := []IObject{}
	for _, elem := range *f {
		dump = append(dump, elem)
	}
	return dump
}

func (f *field) isEmpty() bool {
	return len(*f) == 0
}

// ----------------|

type collision struct {
	o1 IObject
	o0 IObject
}

func (c *collision) collide() {
	c.o0.OnColided(c.o1)
	c.o1.OnColided(c.o0)
}
