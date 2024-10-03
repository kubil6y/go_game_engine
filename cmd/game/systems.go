package main

import (
	"fmt"

	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/logger"
)

type PrintSystem struct {
	ecs.System
	fooState int
}

func NewPrintSystem(logger *logger.Logger, registry *ecs.Registry) *PrintSystem {
	return &PrintSystem{
		System:   ecs.NewSystem("PrintSystem", logger, registry),
		fooState: 88,
	}
}

func (s PrintSystem) GetName() string {
	return s.Name
}

func (s *PrintSystem) Update(dt float32) {
	for _, entity := range s.GetSystemEntities() {
		sprite := s.Registry.GetComponent(entity, SpriteComponent{}).(SpriteComponent)
		s.fooState++
		fmt.Printf("entity id: %d sprite name: %s fooState: %d\n", entity.GetID(), sprite.Name, s.fooState)
	}
}

type AnotherSystem struct {
	ecs.System
	fooState int
}

func NewAnotherSystem(logger *logger.Logger, registry *ecs.Registry) *AnotherSystem {
	return &AnotherSystem{
		System:   ecs.NewSystem("AnotherSystem", logger, registry),
		fooState: 88,
	}
}

func (s AnotherSystem) GetName() string {
	return s.Name
}

func (s *AnotherSystem) Update(dt float32) {
	for _, entity := range s.GetSystemEntities() {
		sprite := s.Registry.GetComponent(entity, SpriteComponent{}).(SpriteComponent)
		s.fooState++
		fmt.Printf("entity id: %d sprite name: %s fooState: %d\n", entity.GetID(), sprite.Name, s.fooState)
	}
}
