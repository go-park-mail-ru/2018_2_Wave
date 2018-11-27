package core

import (
	"Wave/application/room"
	"Wave/utiles/logger"

	"math/rand"
)

type scene struct {
	LG      *logger.Logger
	fields  [][]field
	objects []IObject
	size    Vec2i
}

func newScene(size Vec2i) *scene {
	s := &scene{
		fields: make([][]field, size.X),
		size:   size,
	}
	for i := range s.fields {
		s.fields[i] = make([]field, size.Y)
	}
	return s
}

// ----------------|

func (s *scene) AddObject(o IObject) error {
	if s.isPlaced(o) {
		return room.ErrorAlreadyExists
	}
	s.objects = append(s.objects, o)
	return nil
}

func (s *scene) RemoveObject(o IObject) error {
	if o == nil {
		return room.ErrorNotExists
	}
	for i, expectant := range s.objects {
		if expectant != o {
			continue
		}
		s.objects = append(s.objects[:i], s.objects[i+1:]...)
		return nil
	}
	return room.ErrorNotExists
}

func (s *scene) FindGap(length int) (res []Vec2i, dir Direction) {
	position := Vec2i{
		X: rand.Intn(s.size.X),
		Y: rand.Intn(s.size.Y),
	}
	for i := 0; i < length; i++ {
		if o := s.at(position); o == nil || !o.isEmpty() {
			return s.FindGap(length)
		}
		res = append(res, position)
		position.X++
	}
	return res, Right
}

// ----------------|

func (s *scene) onObjectMove(o IObject, expectedPosition Vec2i) (nextPosition Vec2i, err error) {
	if o != nil {
		currPosition := o.GetPos()
		nextPosition, err := s.validatePosition(expectedPosition)
		if err != nil {
			return Vec2i{}, err
		}

		s.at(currPosition).remove(o)
		s.at(nextPosition).colide(o)
		return nextPosition, nil
	}
	return Vec2i{}, room.ErrorNil
}

func (s *scene) isPlaced(o IObject) bool {
	if o == nil {
		return false
	}
	for _, expectant := range s.objects {
		if expectant == o {
			return true
		}
	}
	return false
}

func (s *scene) validatePosition(expectedPosition Vec2i) (validPosition Vec2i, err error) {
	validPosition.X = expectedPosition.X % s.size.X
	validPosition.Y = expectedPosition.Y % s.size.Y
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

// NOTE:: nil safe
func (f *field) remove(o IObject) {
	if f != nil && o != nil {
		for i, elem := range *f {
			if elem != o {
				continue
			}
			*f = append((*f)[:i], (*f)[i+1:]...)
			return
		}
	}
}

// NOTE:: nil safe
func (f *field) colide(o IObject) {
	if f != nil && o != nil {
		for _, elem := range *f {
			elem.OnColided(o)
			o.OnColided(elem)
		}
	}
}

func (f *field) isEmpty() bool {
	return len(*f) == 0
}
