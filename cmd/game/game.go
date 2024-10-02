package main

import (
	"log"

	"github.com/kubil6y/go_game_engine/pkg/logger"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	TITLE     = "README"
	WIDTH     = 800
	HEIGHT    = 600
	FRAMERATE = 60
)

var (
	playerX, playerY   = int32(WIDTH / 2), int32(HEIGHT / 2)
	playerVX, playerVY = int32(0), int32(0)
)

type Game struct {
	debug    bool
	running  bool
	window   *sdl.Window
	renderer *sdl.Renderer
	logger   *logger.Logger

	WindowWidth  int32
	WindowHeight int32

	millisecondsPreviousFrame uint64 // check type
}

func NewGame() *Game {
	logger := logger.New(logger.WithLogLevel(logger.LevelDebug))
	return &Game{
		debug:                     false,
		running:                   false,
		WindowWidth:               WIDTH,
		WindowHeight:              HEIGHT,
		millisecondsPreviousFrame: 0,
		logger:                    logger,
	}
}

func (g *Game) Initialize() error {
	g.logger.Debug("game initialize called", nil)

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
}

func (g *Game) LoadLevel() {
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
			if t.State == sdl.RELEASED {
				if t.Keysym.Sym == sdl.K_LEFT {
					playerVX -= 1
				} else if t.Keysym.Sym == sdl.K_RIGHT {
					playerVX += 1
				}
				if t.Keysym.Sym == sdl.K_UP {
					playerVY -= 1
				} else if t.Keysym.Sym == sdl.K_DOWN {
					playerVY += 1
				}
			}
			break
		}
	}
}

func (g *Game) Update() {
	// loopTime := loop(surface)
	// window.UpdateSurface()

	// delay := (1000 / FRAMERATE) - loopTime
	// sdl.Delay(delay)
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
	if err := g.renderer.DrawLines(points); err != nil {
		log.Fatal(err)
	}
	g.renderer.Present()
}

func (g *Game) Destroy() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}
