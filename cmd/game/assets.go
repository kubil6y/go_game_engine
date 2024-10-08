package main

import (
	"bufio"
	"os"
	"strconv"

	"github.com/kubil6y/go_game_engine/pkg/asset_store"
	"github.com/kubil6y/go_game_engine/pkg/vector"
)

const (
	IMG_Chopper asset_store.AssetID = iota
	IMG_Tank
	IMG_Tilemap
)

const (
	tileSize   = 32
	tileScale  = 2.0
	mapNumCols = 25
	mapNumRows = 20
)

func (g *Game) LoadAssets() error {
	g.assetStore.AddTexture(g.renderer, IMG_Tilemap, "./assets/tilemaps/jungle.png")
	g.assetStore.AddTexture(g.renderer, IMG_Tank, "./assets/images/tank-panther-left.png")
	g.assetStore.AddTexture(g.renderer, IMG_Chopper, "./assets/images/chopper-spritesheet.png")

	// render the map
	mapFile, err := os.Open("./assets/tilemaps/jungle.map")
	if err != nil {
		g.logger.Fatal(err, "failed to read map file", nil)
	}
	defer mapFile.Close()

	reader := bufio.NewReader(mapFile)
	for y := 0; y < mapNumRows; y++ {
		for x := 0; x < mapNumCols; x++ {
			// Read first character
			ch, err := reader.ReadByte()
			if err != nil {
				g.logger.Fatal(err, "Error reading map file", nil)
				return err
			}
			srcRectY, _ := strconv.Atoi(string(ch))
			srcRectY *= tileSize

			ch, err = reader.ReadByte()
			if err != nil {
				g.logger.Fatal(err, "Error reading map file", nil)
				return err
			}
			srcRectX, _ := strconv.Atoi(string(ch))
			srcRectX *= tileSize

			reader.Discard(1)

			tile := g.registry.CreateEntity()
			g.registry.AddComponent(tile,
				TRANSFORM_COMPONENT,
				TransformComponent{
					Position: vector.Vec2{
						X: float32(x) * (tileScale * tileSize),
						Y: float32(y) * (tileScale * tileSize),
					},
					Scale:    vector.Vec2{X: tileScale, Y: tileScale},
					Rotation: 0.0,
				},
			)

			g.registry.AddComponent(tile, SPRITE_COMPONENT, NewSpriteComponent(IMG_Tilemap, tileSize, tileSize, 0, false, srcRectX, srcRectY))
		}
	}

	g.mapWidth = mapNumCols * tileSize * tileScale
	g.mapHeight = mapNumRows * tileSize * tileScale

	return nil
}
