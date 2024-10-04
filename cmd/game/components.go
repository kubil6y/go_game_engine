package main

import (
	"math"

	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/vector"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	MAX_COMPONENTS_AMOUNT = 32
)

var (
	componentTypeRegistry = ecs.NewTypeRegistry(MAX_COMPONENTS_AMOUNT)
	systemTypeRegistry    = ecs.NewTypeRegistry(math.MaxInt)
)

type SpriteComponent struct {
	Name    string
	AssetID asset_store.AssetID
	Width   int
	Height  int
	ZIndex  int
	IsFixed bool
	SrcRect sdl.Rect
}

func NewSpriteComponent(assetID asset_store.AssetID, width, height, zIndex int, isFixed bool, srcRectX, srcRectY int) SpriteComponent {
	return SpriteComponent{
		Name:    "SpriteComponent",
		AssetID: assetID,
		Width:   width,
		Height:  height,
		ZIndex:  zIndex,
		IsFixed: isFixed,
		SrcRect: sdl.Rect{
			X: int32(srcRectX),
			Y: int32(srcRectY),
			W: int32(width),
			H: int32(height),
		},
	}
}

func (c SpriteComponent) GetID() (int, error) {
	return componentTypeRegistry.Get(c)
}

func (c SpriteComponent) String() string {
	return "SpriteComponent"
}

type TransformComponent struct {
	Position vector.Vec2
	Scale    vector.Vec2
	Rotation float32
}

func (c TransformComponent) GetID() (int, error) {
	return componentTypeRegistry.Get(c)
}

func (c TransformComponent) String() string {
	return "TransformComponent"
}

type BoxColliderComponent struct {
	Width  int
	Height int
	Offset vector.Vec2
}

func (c BoxColliderComponent) GetID() (int, error) {
	return componentTypeRegistry.Get(c)
}

func (c BoxColliderComponent) String() string {
	return "BoxColliderComponent"
}

type RigidbodyComponent struct {
	Velocity vector.Vec2
}

func (c RigidbodyComponent) GetID() (int, error) {
	return componentTypeRegistry.Get(c)
}

func (c RigidbodyComponent) String() string {
	return "RigidBodyComponent"
}
