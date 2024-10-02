package ecs

import (
	"github.com/kubil6y/go_game_engine/internal/type_registry"
)

var componentTypeRegistry = type_registry.New(MAX_COMPONENTS_AMOUNT)

type Component interface {
    GetID() int
}

type SpriteComponent struct {}
type BoxColliderComponent struct {}

func (c SpriteComponent) GetID() int {
	componentID, _ := componentTypeRegistry.Register(SpriteComponent{})
    return componentID
}

func (c BoxColliderComponent) GetID() int {
	componentID, _ := componentTypeRegistry.Register(BoxColliderComponent{})
    return componentID
}
