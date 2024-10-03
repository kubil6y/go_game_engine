package asset_store

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type AssetID int

type AssetStore struct {
	textures map[AssetID]*sdl.Texture
}

func New() *AssetStore {
	return &AssetStore{
		textures: make(map[AssetID]*sdl.Texture),
	}
}

func (s *AssetStore) AddTexture(renderer *sdl.Renderer, assetID AssetID, filepath string) error {
	surface, err := img.Load(filepath)
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}

	s.textures[assetID] = texture
	return nil
}

func (s *AssetStore) GetTexture(assetID AssetID) *sdl.Texture {
	return s.textures[assetID]
}

func (s *AssetStore) GetOrLoadTexture(renderer *sdl.Renderer, assetID AssetID, filepath string) (*sdl.Texture, error) {
	texture, exists := s.textures[assetID]
	if exists {
		return texture, nil
	}
	err := s.AddTexture(renderer, assetID, filepath)
	if err != nil {
		return nil, err
	}
	return s.textures[assetID], nil
}

func (s *AssetStore) Clear() {
	for assetID, texture := range s.textures {
		texture.Destroy()
		delete(s.textures, assetID)
	}
}
