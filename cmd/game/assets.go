package main

import "github.com/kubil6y/go_game_engine/pkg/asset_store"

const (
	IMG_Chopper asset_store.AssetID = iota
	IMG_Tank
)

func (g *Game) LoadAssets() error {
	g.assetStore.AddTexture(g.renderer, IMG_Chopper, "./assets/images/tank-panther-right.png")
    return nil
}
