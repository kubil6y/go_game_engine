package main

import (
	"fmt"

	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/eventbus"
	"github.com/kubil6y/go_game_engine/pkg/logger"
	"github.com/kubil6y/go_game_engine/pkg/vector"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	TITLE  = "README"
	WIDTH  = 800
	HEIGHT = 600
	FPS    = 60

	MILLISECONDS_PER_FRAME = 1000 / FPS
)

type Game struct {
	debug        bool
	running      bool
	msPrevFrame  uint32
	windowWidth  int32
	windowHeight int32
	mapWidth     float32
	mapHeight    float32
	camera       sdl.Rect
	window       *sdl.Window
	renderer     *sdl.Renderer
	logger       *logger.Logger
	assetStore   *asset_store.AssetStore
	registry     ecs.Registry
	events       *eventbus.EventBus
}

func NewGame() *Game {
	logger := logger.New(logger.WithLogLevel(logger.LEVEL_DEBUG))
	return &Game{
		windowWidth:  WIDTH,
		windowHeight: HEIGHT,
		logger:       logger,
		registry:     *ecs.NewRegistry(MAX_COMPONENTS_AMOUNT, logger),
		assetStore:   asset_store.New(),
		events:       eventbus.NewEventBus(),
		debug:        true,
	}
}

func (g *Game) Initialize() error {
	g.logger.Debug("Game Initialize called", nil)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		g.logger.Fatal(err, "failed to initialize sdl", nil)
		return err
	}

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WIDTH, HEIGHT, sdl.WINDOW_BORDERLESS)
	if err != nil {
		g.logger.Fatal(err, "failed to create window", nil)
		return err
	}
	g.window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		g.logger.Fatal(err, "failed to create renderer", nil)
		return err
	}
	g.renderer = renderer

	err = window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		g.logger.Error(err, "failed to set fullscreen", nil)
		return err
	}

	err = renderer.SetLogicalSize(g.windowWidth, g.windowHeight)
	if err != nil {
		g.logger.Error(err, "failed to set logical size", nil)
		return err
	}
	// init camera size
	g.camera = sdl.Rect{
		X: 0,
		Y: 0,
		W: WIDTH,
		H: HEIGHT,
	}
	g.running = true
	return nil
}

func (g *Game) Setup() {
	g.LoadLevel()
}

func (g *Game) LoadLevel() {
	if err := g.LoadAssets(); err != nil {
		g.logger.Fatal(err, fmt.Sprintf("failed to load assets"), nil)
	}

	chopper := g.registry.CreateEntity()
	g.registry.AddComponent(chopper, CAMERA_FOLLOW_COMPONENT, CameraFollowComponent{})
	g.registry.AddComponent(chopper, SPRITE_COMPONENT, NewSpriteComponent(IMG_Chopper, 32, 32, 1, false, 0, 0))
	g.registry.AddComponent(chopper, ANIMATION_COMPONENT, NewAnimationComponent(2, 10, true))
	g.registry.AddComponent(chopper, TRANSFORM_COMPONENT, TransformComponent{
		Position: vector.Vec2{X: 50, Y: 50},
		Scale:    vector.Vec2{X: 1, Y: 1},
		Rotation: 0,
	})
	g.registry.AddComponent(chopper, BOX_COLLIDER_COMPONENT, BoxColliderComponent{
		Width:  32,
		Height: 32,
		Offset: vector.NewZeroVec2(),
	})
	g.registry.AddComponent(chopper, RIGIDBODY_COMPONENT, RigidbodyComponent{
		Velocity: vector.NewZeroVec2(),
	})
	g.registry.AddComponent(chopper, KEYBOARD_CONTROLLED_COMPONENT, KeyboardControlledComponent{
		upVelocity:    vector.Vec2{X: 0, Y: -120},
		downVelocity:  vector.Vec2{X: 0, Y: 120},
		leftVelocity:  vector.Vec2{X: -120, Y: 0},
		rightVelocity: vector.Vec2{X: 120, Y: 0},
	})

	tankSpawner := g.registry.CreateEntity()
	g.registry.AddComponent(tankSpawner, TANK_SPAWNER_COMPONENT, TankSpawnerComponent{})

	tank := g.registry.CreateEntity()
	g.registry.AddComponent(tank, SPRITE_COMPONENT, NewSpriteComponent(IMG_Tank, 32, 32, 1, false, 0, 0))
	g.registry.AddComponent(tank, TRANSFORM_COMPONENT, TransformComponent{
		Position: vector.Vec2{X: 100, Y: 200},
		Scale:    vector.Vec2{X: 1, Y: 1},
		Rotation: 0,
	})
	g.registry.AddComponent(tank, RIGIDBODY_COMPONENT, RigidbodyComponent{
		Velocity: vector.Vec2{X: 30, Y: 0},
	})
	g.registry.AddComponent(tank, BOX_COLLIDER_COMPONENT, BoxColliderComponent{
		Width:  32,
		Height: 32,
		Offset: vector.NewZeroVec2(),
	})

	tank2 := g.registry.CreateEntity()
	g.registry.AddComponent(tank2, SPRITE_COMPONENT, NewSpriteComponent(IMG_Tank, 32, 32, 1, false, 0, 0))
	g.registry.AddComponent(tank2, TRANSFORM_COMPONENT, TransformComponent{
		Position: vector.Vec2{X: 400, Y: 200},
		Scale:    vector.Vec2{X: 1, Y: 1},
		Rotation: 0,
	})
	g.registry.AddComponent(tank2, RIGIDBODY_COMPONENT, RigidbodyComponent{
		Velocity: vector.Vec2{X: -30, Y: 0},
	})
	g.registry.AddComponent(tank2, BOX_COLLIDER_COMPONENT, BoxColliderComponent{
		Width:  32,
		Height: 32,
		Offset: vector.NewZeroVec2(),
	})

	// Create systems
	renderSystem := NewRenderSystem(g.logger, &g.registry, g.renderer, g.assetStore, &g.camera)
	movementSystem := NewMovementSystem(g.logger, &g.registry)
	animationSystem := NewAnimationSystem(g.logger, &g.registry)
	collisionSystem := NewCollisionSystem(g.logger, &g.registry, g.events)
	renderCollisionSystem := NewRenderCollisionSystem(g.logger, &g.registry, g.renderer, &g.camera)
	damageSystem := NewDamageSystem(g.logger, &g.registry, g.events)
	keyboardControlSystem := NewKeyboardControlSystem(g.logger, &g.registry, g.events)
	cameraMovementSystem := NewCameraMovementSystem(g.logger, &g.registry, &g.camera, &g.mapWidth, &g.mapHeight)
	tankSpawnerSystem := NewTankSpawnerSystem(g.logger, &g.registry, &g.camera)

	// Register systems
	g.registry.AddSystem(RENDER_SYSTEM, renderSystem)
	g.registry.AddSystem(MOVEMENT_SYSTEM, movementSystem)
	g.registry.AddSystem(ANIMATION_SYSTEM, animationSystem)
	g.registry.AddSystem(COLLISION_SYSTEM, collisionSystem)
	g.registry.AddSystem(RENDER_COLLISION_SYSTEM, renderCollisionSystem)
	g.registry.AddSystem(DAMAGE_SYSTEM, damageSystem)
	g.registry.AddSystem(KEYBOARD_CONTROL_SYSTEM, keyboardControlSystem)
	g.registry.AddSystem(CAMERA_MOVEMENT_SYSTEM, cameraMovementSystem)
	g.registry.AddSystem(TANK_SPAWNER_SYSTEM, tankSpawnerSystem)

	// Subscribe to events
	g.registry.GetSystem(DAMAGE_SYSTEM).SubscribeToEvents()
	g.registry.GetSystem(KEYBOARD_CONTROL_SYSTEM).SubscribeToEvents()
}

func (g *Game) Run() {
	g.Setup()
	for g.running {
		g.ProcessInput()
		g.Update()
		g.Render()
	}
}

func (g *Game) ProcessInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			g.logger.Debug("Quitting with escape", nil)
			g.running = false
			break
		case *sdl.KeyboardEvent:
			if t.State == sdl.PRESSED {
				g.events.Emit(KEYDOWN_EVENT, KeydownEvent{
					Keysym: t.Keysym,
				})

				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					g.logger.Debug("Quitting with escape", nil)
					g.running = false
					break
				case sdl.K_o:
					g.debug = !g.debug
				}
			}
			break
		}
	}
}

func (g *Game) Update() {
	waitDuration := MILLISECONDS_PER_FRAME - (sdl.GetTicks() - g.msPrevFrame)
	if waitDuration > 0 && waitDuration <= MILLISECONDS_PER_FRAME {
		sdl.Delay(waitDuration)
	}
	dt := float32(sdl.GetTicks()-g.msPrevFrame) / 1000.0
	g.msPrevFrame = sdl.GetTicks()

	g.registry.Update()

	movementSystem := g.registry.GetSystem(MOVEMENT_SYSTEM).(*MovementSystem)
	animationSystem := g.registry.GetSystem(ANIMATION_SYSTEM).(*AnimationSystem)
	collisionSystem := g.registry.GetSystem(COLLISION_SYSTEM).(*CollisionSystem)
	cameraMovementSystem := g.registry.GetSystem(CAMERA_MOVEMENT_SYSTEM).(*CameraMovementSystem)
	tankSpawnerSystem := g.registry.GetSystem(TANK_SPAWNER_SYSTEM).(*TankSpawnerSystem)

	movementSystem.Update(dt)
	animationSystem.Update(dt)
	collisionSystem.Update(dt)
	cameraMovementSystem.Update(dt)
	tankSpawnerSystem.Update(dt)
}

func (g *Game) Render() {
	g.renderer.SetDrawColor(0, 0, 0, 0)
	g.renderer.Clear()

	renderSystem := g.registry.GetSystem(RENDER_SYSTEM).(*RenderSystem)
	renderCollisionSystem := g.registry.GetSystem(RENDER_COLLISION_SYSTEM).(*RenderCollisionSystem)

	renderSystem.Update(0)
	if g.debug {
		renderCollisionSystem.Update(0)
	}

	g.renderer.Present()
}

func (g *Game) Destroy() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
