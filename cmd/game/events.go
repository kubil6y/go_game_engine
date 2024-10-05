package main

import (
	"github.com/kubil6y/go_game_engine/pkg/ecs"
	"github.com/kubil6y/go_game_engine/pkg/eventbus"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	KEYDOWN_EVENT eventbus.EventID = iota
	COLLISION_EVENT
)

type KeydownEvent struct {
	Keysym sdl.Keysym
}

type CollisionEvent struct {
	a ecs.Entity
	b ecs.Entity
}
