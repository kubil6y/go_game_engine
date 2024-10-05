package main

import (
	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/vector"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SPRITE_COMPONENT ecs.ComponentTypeID = iota
	BOX_COLLIDER_COMPONENT
	TRANSFORM_COMPONENT
	RIGIDBODY_COMPONENT
	ANIMATION_COMPONENT
	KEYBOARD_CONTROLLED_COMPONENT
)

const (
	MAX_COMPONENTS_AMOUNT = 32
)

// ////////////////////////////////////////////////
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

func (c SpriteComponent) GetID() int {
	return int(SPRITE_COMPONENT)
}

func (c SpriteComponent) String() string {
	return "SpriteComponent"
}

// ////////////////////////////////////////////////
type TransformComponent struct {
	Position vector.Vec2
	Scale    vector.Vec2
	Rotation float32
}

func (c TransformComponent) GetID() int {
	return int(TRANSFORM_COMPONENT)
}

func (c TransformComponent) String() string {
	return "TransformComponent"
}

// ////////////////////////////////////////////////
type BoxColliderComponent struct {
	Width  float32
	Height float32
	Offset vector.Vec2
}

func (c BoxColliderComponent) GetID() int {
	return int(BOX_COLLIDER_COMPONENT)
}

func (c BoxColliderComponent) String() string {
	return "BoxColliderComponent"
}

// ////////////////////////////////////////////////
type RigidbodyComponent struct {
	Velocity vector.Vec2
}

func (c RigidbodyComponent) GetID() int {
	return int(RIGIDBODY_COMPONENT)
}

func (c RigidbodyComponent) String() string {
	return "RigidBodyComponent"
}

// ////////////////////////////////////////////////
type AnimationComponent struct {
	numFrames      int
	currentFrame   int
	frameRateSpeed int // ms
	loop           bool
	startTime      uint32
}

func NewAnimationComponent(numFrames, frameRateSpeed int, loop bool) AnimationComponent {
	if numFrames < 1 || frameRateSpeed < 0 {
		panic("invalid parameter")
	}
	return AnimationComponent{
		numFrames:      numFrames,
		currentFrame:   0,
		frameRateSpeed: frameRateSpeed,
		loop:           loop,
		startTime:      sdl.GetTicks(),
	}
}

func (c AnimationComponent) GetID() int {
	return int(ANIMATION_COMPONENT)
}

func (c AnimationComponent) String() string {
	return "AnimationComponent"
}

// ////////////////////////////////////////////////
type KeyboardControlledComponent struct {
	upVelocity    vector.Vec2
	downVelocity  vector.Vec2
	leftVelocity  vector.Vec2
	rightVelocity vector.Vec2
}

func (c KeyboardControlledComponent) GetID() int {
	return int(KEYBOARD_CONTROLLED_COMPONENT)
}

func (c KeyboardControlledComponent) String() string {
	return "KeyboardControlledComponent"
}
