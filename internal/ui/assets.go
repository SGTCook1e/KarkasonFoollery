package ui

import (
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	Tiles map[string]*ebiten.Image
}

func LoadAssets() (*Assets, error) {
	assets := &Assets{
		Tiles: make(map[string]*ebiten.Image),
	}
	path := filepath.Join("assets", "tiles")
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		img, err := loadImage(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		assets.Tiles[file.Name()[:len(file.Name())-4]] = img
	}
	return assets, nil
}

func loadImage(path string) (*ebiten.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(image), nil
}
