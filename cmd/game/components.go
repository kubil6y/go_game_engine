package main

import "github.com/kubil6y/go_game_engine/internal/type_registry"

const (
	MAX_COMPONENTS_AMOUNT = 32
)

var componentTypeRegistry = type_registry.New(MAX_COMPONENTS_AMOUNT)

type SpriteComponent struct {
	Name string
}

func (c SpriteComponent) GetID() int {
	componentID, _ := componentTypeRegistry.Register(SpriteComponent{})
	return componentID
}

func (c SpriteComponent) String() string {
	return "SpriteComponent"
}

type BoxColliderComponent struct {
	X int
	Y int
}

func (c BoxColliderComponent) GetID() int {
	componentID, _ := componentTypeRegistry.Register(BoxColliderComponent{})
	return componentID
}

func (c BoxColliderComponent) String() string {
	return "BoxColliderComponent"
}
