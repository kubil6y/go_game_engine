package main

import (
	"fmt"

	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/bitset"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/eventbus"
	"github.com/kubil6y/go_game_engine/pkg/logger"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	RENDER_SYSTEM ecs.SystemTypeID = iota
	RENDER_COLLISION_SYSTEM
	MOVEMENT_SYSTEM
	ANIMATION_SYSTEM
	COLLISION_SYSTEM
	DAMAGE_SYSTEM
	KEYBOARD_CONTROL_SYSTEM
)

// RENDER SYSTEM ////////////////////////////////////////////////
type RenderSystem struct {
	*ecs.BaseSystem
	renderer   *sdl.Renderer
	assetStore *asset_store.AssetStore
}

func NewRenderSystem(logger *logger.Logger, registry *ecs.Registry, renderer *sdl.Renderer, assetStore *asset_store.AssetStore) *RenderSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(SPRITE_COMPONENT))
	bs.Set(int(TRANSFORM_COMPONENT))

	return &RenderSystem{
		BaseSystem: ecs.NewBaseSystem("RenderSystem", logger, registry, bs),
		renderer:   renderer,
		assetStore: assetStore,
	}
}

func (s RenderSystem) GetName() string {
	return s.Name
}

func (s *RenderSystem) Update(dt float32) {
	var currZIndex int
	var maxZIndex int

	for currZIndex <= maxZIndex {
		for _, entity := range s.GetSystemEntities() {
			sprite := s.Registry.GetComponentPtr(entity, SPRITE_COMPONENT).(*SpriteComponent)
			if maxZIndex < sprite.ZIndex {
				maxZIndex = sprite.ZIndex
			}
			if currZIndex != sprite.ZIndex {
				continue
			}
			tf := s.Registry.GetComponentPtr(entity, TRANSFORM_COMPONENT).(*TransformComponent)
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

// MOVEMENT SYSTEM ////////////////////////////////////////////////
type MovementSystem struct {
	*ecs.BaseSystem
	renderer   *sdl.Renderer
	assetStore *asset_store.AssetStore
}

func NewMovementSystem(logger *logger.Logger, registry *ecs.Registry) *MovementSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(RIGIDBODY_COMPONENT))
	bs.Set(int(TRANSFORM_COMPONENT))
	return &MovementSystem{
		BaseSystem: ecs.NewBaseSystem("MovementSystem", logger, registry, bs),
	}
}

func (s MovementSystem) GetName() string {
	return s.Name
}

func (s *MovementSystem) Update(dt float32) {
	for _, entity := range s.GetSystemEntities() {
		tf := s.Registry.GetComponentPtr(entity, TRANSFORM_COMPONENT).(*TransformComponent)
		rb := s.Registry.GetComponentPtr(entity, RIGIDBODY_COMPONENT).(*RigidbodyComponent)
		tf.Position.X += rb.Velocity.X * dt
		tf.Position.Y += rb.Velocity.Y * dt
	}
}

// ANIMATION SYSTEM ////////////////////////////////////////////////
type AnimationSystem struct {
	*ecs.BaseSystem
	renderer   *sdl.Renderer
	assetStore *asset_store.AssetStore
}

func NewAnimationSystem(logger *logger.Logger, registry *ecs.Registry) *AnimationSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(SPRITE_COMPONENT))
	bs.Set(int(ANIMATION_COMPONENT))
	return &AnimationSystem{
		BaseSystem: ecs.NewBaseSystem("AnimationSystem", logger, registry, bs),
	}
}

func (s AnimationSystem) GetName() string {
	return s.Name
}

func (s *AnimationSystem) Update(dt float32) {
	for _, entity := range s.GetSystemEntities() {
		sprite := s.Registry.GetComponentPtr(entity, SPRITE_COMPONENT).(*SpriteComponent)
		animation := s.Registry.GetComponentPtr(entity, ANIMATION_COMPONENT).(*AnimationComponent)

		// TODO support loop
		animation.currentFrame = int((sdl.GetTicks() - animation.startTime)) *
			animation.frameRateSpeed / 1000 %
			animation.numFrames
		sprite.SrcRect.X = int32(animation.currentFrame * sprite.Width)
	}
}

// COLLISION SYSTEM ////////////////////////////////////////////////
type CollisionSystem struct {
	*ecs.BaseSystem
	events *eventbus.EventBus
}

func NewCollisionSystem(logger *logger.Logger, registry *ecs.Registry, events *eventbus.EventBus) *CollisionSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(TRANSFORM_COMPONENT))
	bs.Set(int(BOX_COLLIDER_COMPONENT))
	return &CollisionSystem{
		BaseSystem: ecs.NewBaseSystem("CollisionSystem", logger, registry, bs),
		events:     events,
	}
}

func (s CollisionSystem) GetName() string {
	return s.Name
}

func (s *CollisionSystem) Update(dt float32) {
	entities := s.GetSystemEntities()
	for i := 0; i < len(entities); i++ {
		a := entities[i]
		atf := s.Registry.GetComponentPtr(a, TRANSFORM_COMPONENT).(*TransformComponent)
		acol := s.Registry.GetComponentPtr(a, BOX_COLLIDER_COMPONENT).(*BoxColliderComponent)
		for j := 0; j < len(entities); j++ {
			b := entities[j]
			if a.GetID() == b.GetID() {
				continue
			}
			btf := s.Registry.GetComponentPtr(b, TRANSFORM_COMPONENT).(*TransformComponent)
			bcol := s.Registry.GetComponentPtr(b, BOX_COLLIDER_COMPONENT).(*BoxColliderComponent)
			if CheckAABB(atf, btf, acol, bcol) {
				s.events.Emit(COLLISION_EVENT, CollisionEvent{
					a: a,
					b: b,
				})
				s.Logger.Debug(fmt.Sprintf("COLLISION_EVENT fired entity=%d and entity=%d", a.GetID(), b.GetID()), nil)
			}
		}
	}
}

func CheckAABB(atf, btf *TransformComponent, acol, bcol *BoxColliderComponent) bool {
	aMinX := atf.Position.X + acol.Offset.X*atf.Scale.X
	aMaxX := aMinX + acol.Width*atf.Scale.X
	aMinY := atf.Position.Y + acol.Offset.Y*atf.Scale.Y
	aMaxY := aMinY + acol.Height*atf.Scale.Y
	bMinX := btf.Position.X + bcol.Offset.X*btf.Scale.X
	bMaxX := bMinX + bcol.Width*btf.Scale.X
	bMinY := btf.Position.Y + bcol.Offset.Y*btf.Scale.Y
	bMaxY := bMinY + bcol.Height*btf.Scale.Y
	return aMinX < bMaxX && aMaxX > bMinX && aMinY < bMaxY && aMaxY > bMinY
}

// RENDER COLLISION SYSTEM ////////////////////////////////////////////////
type RenderCollisionSystem struct {
	*ecs.BaseSystem
	renderer *sdl.Renderer
}

func NewRenderCollisionSystem(logger *logger.Logger, registry *ecs.Registry, renderer *sdl.Renderer) *RenderCollisionSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(TRANSFORM_COMPONENT))
	bs.Set(int(BOX_COLLIDER_COMPONENT))
	return &RenderCollisionSystem{
		BaseSystem: ecs.NewBaseSystem("RenderCollisionSystem", logger, registry, bs),
		renderer:   renderer,
	}
}

func (s RenderCollisionSystem) GetName() string {
	return s.Name
}

func (s *RenderCollisionSystem) Update(dt float32) {
	for _, entity := range s.GetSystemEntities() {
		tf := s.Registry.GetComponentPtr(entity, TRANSFORM_COMPONENT).(*TransformComponent)
		col := s.Registry.GetComponentPtr(entity, BOX_COLLIDER_COMPONENT).(*BoxColliderComponent)
		rect := sdl.Rect{
			X: int32(tf.Position.X + col.Offset.X),
			Y: int32(tf.Position.Y + col.Offset.Y),
			W: int32(tf.Scale.X * col.Width),
			H: int32(tf.Scale.Y * col.Height),
		}
		s.renderer.SetDrawColor(255, 0, 0, 255)
		s.renderer.DrawRect(&rect)
	}
}

// DAMAGE SYSTEM ////////////////////////////////////////////////
type DamageSystem struct {
	*ecs.BaseSystem
	events *eventbus.EventBus
}

func NewDamageSystem(logger *logger.Logger, registry *ecs.Registry, events *eventbus.EventBus) *DamageSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(RIGIDBODY_COMPONENT))
	bs.Set(int(TRANSFORM_COMPONENT))
	return &DamageSystem{
		BaseSystem: ecs.NewBaseSystem("DamageSystem", logger, registry, bs),
		events:     events,
	}
}

func (s DamageSystem) GetName() string {
	return s.Name
}

func (s *DamageSystem) SubscribeToEvents() {
	s.events.On(COLLISION_EVENT, s.OnCollision)
}

func (s *DamageSystem) OnCollision(payload any) {
	p, ok := payload.(CollisionEvent)
	if !ok {
		s.Logger.Debug(fmt.Sprintf("DamageSystem:OnCollision invalid payload: %+v", payload), nil)
	}
	s.Registry.KillEntity(p.a)
	s.Registry.KillEntity(p.b)
	s.Logger.Debug(fmt.Sprintf("COLLISION_EVENT captured entity=%d and entity=%d", p.a.GetID(), p.b.GetID()), nil)
}

// KeyboardControl SYSTEM ////////////////////////////////////////////////
type KeyboardControlSystem struct {
	*ecs.BaseSystem
	events *eventbus.EventBus
}

func NewKeyboardControlSystem(logger *logger.Logger, registry *ecs.Registry, events *eventbus.EventBus) *KeyboardControlSystem {
	bs := bitset.NewBitset32()
	bs.Set(int(SPRITE_COMPONENT))
	bs.Set(int(RIGIDBODY_COMPONENT))
	bs.Set(int(KEYBOARD_CONTROLLED_COMPONENT))
	return &KeyboardControlSystem{
		BaseSystem: ecs.NewBaseSystem("KeyboardControlSystem", logger, registry, bs),
		events:     events,
	}
}

func (s KeyboardControlSystem) GetName() string {
	return s.Name
}

func (s *KeyboardControlSystem) SubscribeToEvents() {
	s.events.On(KEYDOWN_EVENT, s.OnKeydown)
}

func (s *KeyboardControlSystem) Update(dt float32) {
}

func (s *KeyboardControlSystem) OnKeydown(payload any) {
	p, ok := payload.(KeydownEvent)
	if !ok {
		return
	}

	for _, entity := range s.GetSystemEntities() {
		keyboard := s.Registry.GetComponentPtr(entity, KEYBOARD_CONTROLLED_COMPONENT).(*KeyboardControlledComponent)
		sprite := s.Registry.GetComponentPtr(entity, SPRITE_COMPONENT).(*SpriteComponent)
		rb := s.Registry.GetComponentPtr(entity, RIGIDBODY_COMPONENT).(*RigidbodyComponent)

		switch p.Keysym.Sym {
		case sdl.K_UP:
			rb.Velocity = keyboard.upVelocity
			sprite.SrcRect.Y = int32(sprite.Height * 0)
		case sdl.K_RIGHT:
			rb.Velocity = keyboard.rightVelocity
			sprite.SrcRect.Y = int32(sprite.Height * 1)
		case sdl.K_DOWN:
			rb.Velocity = keyboard.downVelocity
			sprite.SrcRect.Y = int32(sprite.Height * 2)
		case sdl.K_LEFT:
			rb.Velocity = keyboard.leftVelocity
			sprite.SrcRect.Y = int32(sprite.Height * 3)
		}
	}
}
