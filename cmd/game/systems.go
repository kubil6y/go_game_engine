package main

import (
	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/logger"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderSystem struct {
	*ecs.BaseSystem
	renderer   *sdl.Renderer
	assetStore *asset_store.AssetStore
}

func NewRenderSystem(logger *logger.Logger, registry *ecs.Registry, renderer *sdl.Renderer, assetStore *asset_store.AssetStore) *RenderSystem {
	bs := bitset.NewBitset32()
	bs.Set(componentTypeRegistry.Getx(SpriteComponent{}))
	bs.Set(componentTypeRegistry.Getx(TransformComponent{}))
	return &RenderSystem{
		BaseSystem: ecs.NewBaseSystem("RenderSystem", logger, registry, bs),
		renderer:   renderer,
		assetStore: assetStore,
	}
}

func (s RenderSystem) GetName() string {
	return s.Name
}

func (s *RenderSystem) Update() {
	var currZIndex int
	var maxZIndex int

	for currZIndex <= maxZIndex {
		for _, entity := range s.GetSystemEntities() {
			sprite := s.Registry.GetComponent(entity, SpriteComponent{}).(SpriteComponent)
			if maxZIndex < sprite.ZIndex {
				maxZIndex = sprite.ZIndex
			}
			if currZIndex != sprite.ZIndex {
				continue
			}
			tf := s.Registry.GetComponent(entity, TransformComponent{}).(TransformComponent)
			var dstRect sdl.Rect
			dstRect.X = int32(tf.Position.X)
			dstRect.Y = int32(tf.Position.Y)
			dstRect.W = int32(sprite.Width * int(tf.Scale.X))
			dstRect.H = int32(sprite.Height * int(tf.Scale.Y))
			s.renderer.CopyEx(s.assetStore.GetTexture(sprite.AssetID), &sprite.SrcRect, &dstRect, 0, nil, sdl.FLIP_NONE)
		}
		currZIndex++
	}
}

type MovementSystem struct {
	*ecs.BaseSystem
	renderer   *sdl.Renderer
	assetStore *asset_store.AssetStore
}

func NewMovementSystem(logger *logger.Logger, registry *ecs.Registry) *MovementSystem {
	bs := bitset.NewBitset32()
	bs.Set(componentTypeRegistry.Getx(RigidbodyComponent{}))
	bs.Set(componentTypeRegistry.Getx(TransformComponent{}))
	return &MovementSystem{
		BaseSystem: ecs.NewBaseSystem("MovementSystem", logger, registry, bs),
	}
}

func (s MovementSystem) GetName() string {
	return s.Name
}

func (s *MovementSystem) Update(dt float32) {
	for _, entity := range s.GetSystemEntities() {
		tf := s.Registry.GetComponent(entity, TransformComponent{}).(TransformComponent)
		rb := s.Registry.GetComponent(entity, RigidbodyComponent{}).(RigidbodyComponent)
		tf.Position.X += rb.Velocity.X * dt
		tf.Position.Y += rb.Velocity.Y * dt
        fmt.Printf(
	}
}
