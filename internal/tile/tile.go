package tile

import (
	"encoding/json"
	"os"
)

type Tile struct {
	//ID 1-12 is Rivers
	//ID 1 is the starting river tile
	//ID 12 is end river tile
	ID      int    `json:"id"`
	Texture string `json:"texture"`

	// 0 = 0°
	// 1 = 90°
	// 2 = 180°
	// 3 = 270°
	Dir Direction
	// Sides are stored in tile order: Top, Right, Bottom, Left.
	Sides [4]SideType `json:"sides"`
	// Features []Feature
}

func NewTile(id int, path string, sides [4]SideType) *Tile {
	return &Tile{
		ID:      id,
		Texture: path,
		Sides:   sides,
		// Features: features,
		Dir: Top,
	}
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		ID:      t.ID,
		Texture: t.Texture,
		Sides:   t.Sides,
		// Features: t.Features,
		Dir: t.Dir,
	}
}

func LoadTiles() []*Tile {
	var tiles []*Tile
	data, err := os.ReadFile("assets/data/tile_info.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &tiles)
	if err != nil {
		panic(err)
	}
	return tiles
}

// func (t *Tile) FeatureAt(direction Direction) *Feature {
// 	rotatedDir := Direction((int(direction) - t.Rotation + 4) % 4)
// 	for i := range t.Features {
// 		if t.Features[i].HasSide(rotatedDir) {
// 			return &t.Features[i]
// 		}
// 	}
// 	return nil
// }

func (t *Tile) SideAt(direction Direction) SideType {
	return t.Sides[(int(direction)-int(t.Dir)+4)%4]
}

func (t *Tile) GetSide(direction Direction) SideType {
	return t.SideAt(direction)
}

func (t *Tile) Rotate() {
	t.Dir = (t.Dir + 1) % 4
}
