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
	window       *sdl.Window
	renderer     *sdl.Renderer
	logger       *logger.Logger
	assetStore   *asset_store.AssetStore
	registry     ecs.Registry
	events       *eventbus.EventBus
}

func NewGame() *Game {
	logger := logger.New(logger.WithLogLevel(logger.LevelInfo))
	return &Game{
		windowWidth:  WIDTH,
		windowHeight: HEIGHT,
		logger:       logger,
		registry:     *ecs.NewRegistry(MAX_COMPONENTS_AMOUNT, logger, componentTypeRegistry, systemTypeRegistry),
		assetStore:   asset_store.New(),
		events:       eventbus.NewEventBus(),
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
	g.running = true
	return nil
}

func (g *Game) Setup() {
	g.RegisterComponents()
	g.LoadLevel()
}

func (g *Game) RegisterComponents() {
	g.registry.RegisterComponent(SpriteComponent{})
	g.registry.RegisterComponent(TransformComponent{})
	g.registry.RegisterComponent(BoxColliderComponent{})
	g.registry.RegisterComponent(RigidBodyComponent{})
}

func (g *Game) LoadLevel() {
	if err := g.LoadAssets(); err != nil {
		g.logger.Fatal(err, fmt.Sprintf("failed to load assets"), nil)
	}
	tank := g.registry.CreateEntity()

	g.registry.AddComponent(tank, NewSpriteComponent(IMG_Tank, 32, 32, 1, false, 0, 0))
	g.registry.AddComponent(tank, TransformComponent{
		Position: vector.Vec2{X: 300, Y: 300},
		Scale:    vector.Vec2{X: 1, Y: 1},
		Rotation: 0,
	})

	// Create systems
	printSystem := NewPrintSystem(g.logger, &g.registry)
	renderSystem := NewRenderSystem(g.logger, &g.registry, g.renderer, g.assetStore)

	// Register systems
	g.registry.AddSystem(printSystem)
	g.registry.AddSystem(renderSystem)
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
			if t.State == sdl.RELEASED {
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					g.logger.Debug("Quitting with escape", nil)
					g.running = false
					break
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

	printSystemID, _ := systemTypeRegistry.Get(&PrintSystem{})
	printSystem := g.registry.GetSystem(printSystemID).(*PrintSystem)
	printSystem.Update(dt)
}

func (g *Game) Render() {
	g.renderer.SetDrawColor(0, 0, 0, 0)
	g.renderer.Clear()

	renderSystemID, _ := systemTypeRegistry.Get(&RenderSystem{})
	renderSystem := g.registry.GetSystem(renderSystemID).(*RenderSystem)
	renderSystem.Update()

	g.renderer.Present()
}

func (g *Game) Destroy() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
