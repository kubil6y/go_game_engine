package ecs

import (
	"fmt"

	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/logger"
)

type System interface {
	GetName() string
	AddEntityToSystem(entity Entity)
	RemoveEntityFromSystem(entity Entity)
	GetSystemEntities() []Entity
	GetSignature() *bitset.Bitset32
	RequireComponent(componentID int)
}

type BaseSystem struct {
	Name               string
	componentSignature *bitset.Bitset32
	entities           []Entity
	Logger             *logger.Logger
	Registry           *Registry
}

func NewBaseSystem(name string, logger *logger.Logger, registry *Registry, bitset *bitset.Bitset32) *BaseSystem {
	return &BaseSystem{
		Name:               name,
		componentSignature: bitset,
		entities:           make([]Entity, 0),
		Logger:             logger,
		Registry:           registry,
	}
}

func (s *BaseSystem) AddEntityToSystem(entity Entity) {
	s.entities = append(s.entities, entity)
	s.Logger.Debug(fmt.Sprintf("Entity{%d} added to %s", entity.GetID(), s.Name), nil)
}

func (s *BaseSystem) RemoveEntityFromSystem(entity Entity) {
	index := -1
	for i, e := range s.entities {
		if e.GetID() == entity.GetID() {
			index = i
			break
		}
	}
	if index != -1 {
		s.entities = append(s.entities[:index], s.entities[index+1:]...)
	}
}

func (s *BaseSystem) GetSystemEntities() []Entity {
	return s.entities
}

func (s *BaseSystem) GetSignature() *bitset.Bitset32 {
	return s.componentSignature
}

func (s *BaseSystem) RequireComponent(componentID int) {
	s.componentSignature.Set(componentID)
}
