package main

import (
	"fmt"

	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/logger"
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
	debug      bool
	running    bool
	window     *sdl.Window
	renderer   *sdl.Renderer
	logger     *logger.Logger
	assetStore *asset_store.AssetStore

	WindowWidth  int32
	WindowHeight int32
	registry     ecs.Registry

	millisecondsPreviousFrame uint32
}

func NewGame() *Game {
	logger := logger.New(logger.WithLogLevel(logger.LevelDebug))
	return &Game{
		WindowWidth:  WIDTH,
		WindowHeight: HEIGHT,
		logger:       logger,
		registry:     *ecs.NewRegistry(MAX_COMPONENTS_AMOUNT, logger, componentTypeRegistry, systemTypeRegistry),
		assetStore:   asset_store.New(),
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

	err = renderer.SetLogicalSize(g.WindowWidth, g.WindowHeight)
	if err != nil {
		g.logger.Error(err, "failed to set logical size", nil)
		return err
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
	// NOTE: all the components must be registered beforehand to entities
	// before they are consumed by systems.
	tank := g.registry.CreateEntity()
	tank2 := g.registry.CreateEntity()

	g.registry.AddComponent(tank, SpriteComponent{
		Name: "tank-sprite",
	})
	g.registry.AddComponent(tank, BoxColliderComponent{
		X: 10,
		Y: 20,
	})

	// g.registry.AddComponent(tank2, SpriteComponent{
	// 	Name: "tank-sprite",
	// })
	g.registry.AddComponent(tank2, BoxColliderComponent{
		X: 10,
		Y: 20,
	})

	// Create systems
	printSystem := NewPrintSystem(g.logger, &g.registry)
	anotherSystem := NewAnotherSystem(g.logger, &g.registry)

	// Register system
	g.registry.AddSystem(printSystem)
	g.registry.AddSystem(anotherSystem)
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
	waitDuration := MILLISECONDS_PER_FRAME - (sdl.GetTicks() - g.millisecondsPreviousFrame)
	if waitDuration > 0 && waitDuration <= MILLISECONDS_PER_FRAME {
		sdl.Delay(waitDuration)
	}
	dt := float32(sdl.GetTicks()-g.millisecondsPreviousFrame) / 1000.0
	g.millisecondsPreviousFrame = sdl.GetTicks()

	g.registry.Update()

	printSystemID, _ := systemTypeRegistry.Get(&PrintSystem{})
	printSystem := g.registry.GetSystem(printSystemID).(*PrintSystem)
	printSystem.Update(dt)
}

func (g *Game) Render() {
	g.renderer.SetDrawColor(0, 0, 0, 0)
	g.renderer.Clear()

	g.renderer.SetDrawColor(255, 255, 255, 255)
	points := []sdl.Point{
		{X: 400, Y: 100}, // Top vertex
		{X: 300, Y: 500}, // Bottom left vertex
		{X: 500, Y: 500}, // Bottom right vertex
		{X: 400, Y: 100}, // Closing the triangle
	}

	// Draw the triangle
	g.renderer.DrawLines(points)

	renderRect := sdl.Rect{
		X: 500,
		Y: 300,
		W: 32,
		H: 32,
	}
	g.renderer.CopyEx(g.assetStore.GetTexture(IMG_Chopper), nil, &renderRect, 0, nil, sdl.FLIP_NONE)

	g.renderer.Present()
}

func (g *Game) Destroy() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
