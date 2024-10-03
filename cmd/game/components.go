package main

import (
	"math"

	"github.com/kubil6y/go_game_engine/pkg/ecs"
)

const (
	MAX_COMPONENTS_AMOUNT = 32
)

var (
	componentTypeRegistry = ecs.NewTypeRegistry(MAX_COMPONENTS_AMOUNT)
	systemTypeRegistry    = ecs.NewTypeRegistry(math.MaxInt)
)

type SpriteComponent struct {
	Name string
}

func (c SpriteComponent) GetID() (int, error) {
	return componentTypeRegistry.Get(c)
}

func (c SpriteComponent) String() string {
	return "SpriteComponent"
}

type BoxColliderComponent struct {
	X int
	Y int
}

func (c BoxColliderComponent) GetID() (int, error) {
	return componentTypeRegistry.Get(c)
}

func (c BoxColliderComponent) String() string {
	return "BoxColliderComponent"
}
