package main

import (
	"fmt"

	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/logger"
)

type PrintSystem struct {
	*ecs.BaseSystem
	fooState int
}

func NewPrintSystem(logger *logger.Logger, registry *ecs.Registry) *PrintSystem {
	bs := bitset.NewBitset32()
	bs.Set(componentTypeRegistry.Getx(SpriteComponent{}))
	bs.Set(componentTypeRegistry.Getx(BoxColliderComponent{}))
	return &PrintSystem{
		BaseSystem: ecs.NewBaseSystem("PrintSystem", logger, registry, bs),
		fooState:   88,
	}
}

func (s PrintSystem) GetName() string {
	return s.Name
}

func (s *PrintSystem) Update(dt float32) {
	// for _, entity := range s.GetSystemEntities() {
	// 	sprite := s.Registry.GetComponent(entity, SpriteComponent{}).(SpriteComponent)
	// 	s.fooState++
	// 	fmt.Printf("entity{%d} sprite name: %s fooState: %d\n", entity.GetID(), sprite.Name, s.fooState)
	// }
}

type AnotherSystem struct {
	*ecs.BaseSystem
	fooState int
}

func NewAnotherSystem(logger *logger.Logger, registry *ecs.Registry) *AnotherSystem {
	bs := bitset.NewBitset32()
	bs.Set(componentTypeRegistry.Getx(SpriteComponent{}))
	return &AnotherSystem{
		BaseSystem: ecs.NewBaseSystem("AnotherSystem", logger, registry, bs),
		fooState:   88,
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
