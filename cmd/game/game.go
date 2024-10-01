package main

import "github.com/veandco/go-sdl2/sdl"

type Game struct {
	debug   bool
	running bool
	window  *sdl.Window

	millisecondsPreviousFrame uint64 // check type
}

func (g *Game) Initialize() {
}

func (g *Game) Setup() {
}

func (g *Game) LoadLevel() {
}

func (g *Game) Run() {
}

func (g *Game) ProcessInput() {
}

func (g *Game) Update() {
}

func (g *Game) Render() {
}

func (g *Game) Destroy() {
}
