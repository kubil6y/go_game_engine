package main

import (
	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/logger"
	"github.com/veandco/go-sdl2/sdl"
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
	// }
}

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
	for _, entity := range s.GetSystemEntities() {
		sprite := s.Registry.GetComponent(entity, SpriteComponent{}).(SpriteComponent)
		tf := s.Registry.GetComponent(entity, TransformComponent{}).(TransformComponent)
		var dstRect sdl.Rect
		dstRect.X = int32(tf.Position.X)
		dstRect.Y = int32(tf.Position.Y)
		dstRect.W = int32(sprite.Width * int(tf.Scale.X))
		dstRect.H = int32(sprite.Height * int(tf.Scale.Y))
		s.renderer.CopyEx(s.assetStore.GetTexture(IMG_Tilemap), &sprite.SrcRect, &dstRect, 0, nil, sdl.FLIP_NONE)
	}
}
