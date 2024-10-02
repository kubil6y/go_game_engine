package ecs

import (
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/set"
)

type System struct {
	componentSignature bitset.Bitset32
	entities           set.Set[Entity]
}

func (s *System) AddEntityToSystem(entity Entity) {
	s.entities.Add(entity)
}

func (s *System) RemoveEntityFromSystem(entity Entity) {
	s.entities.Remove(entity)
}

func (s *System) GetSignature() bitset.Bitset32 {
	return s.componentSignature
}

func (s *System) RequireComponent(componentID int) {
	s.componentSignature.Set(componentID)
}
