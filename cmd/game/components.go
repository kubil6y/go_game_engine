package main

import (
	"math"

	"github.com/kubil6y/go_game_engine/internal/type_registry"
)

const (
	MAX_COMPONENTS_AMOUNT = 32
)

var (
	componentTypeRegistry = type_registry.New(MAX_COMPONENTS_AMOUNT)
	systemTypeRegistry    = type_registry.New(math.MaxInt)
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
