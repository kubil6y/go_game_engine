package ecs

import (
	"github.com/kubil6y/go_game_engine/pkg/bitset"
)

type System struct {
	componentSignature bitset.Bitset32
	entities           []Entity
}

func (s *System) AddEntityToSystem(entity Entity) {
	s.entities = append(s.entities, entity)
}

func (s *System) RemoveEntityFromSystem(entity Entity) {
	index := -1
	for i := 0; i < len(s.entities); i++ {
		if s.entities[i].GetID() == entity.GetID() {
			index = i
			break
		}
	}
	if index != -1 {
		s.entities = append(s.entities[:index], s.entities[index+1:]...)
	}
}

func (s *System) GetSystemEntities() []Entity {
	return s.entities
}

func (s *System) GetSignature() bitset.Bitset32 {
	return s.componentSignature
}

func (s *System) RequireComponent(componentID int) {
	s.componentSignature.Set(componentID)
}
