package board

import (
	"encoding/json"
	"fmt"
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
	Orientation Direction
	// Sides are stored in tile order: Top, Right, Bottom, Left.
	Sides     [4]SideType `json:"sides"`
	Features  []Feature   `json:"features"`
	Monastery bool        `json:"monastery"`
}

func NewTile(id int, path string, sides [4]SideType) *Tile {
	return &Tile{
		ID:      id,
		Texture: path,
		Sides:   sides,
		// Features: features,
		Orientation: Top,
	}
}

func (t *Tile) Clone() *Tile {
	clone := &Tile{
		ID:          t.ID,
		Texture:     t.Texture,
		Sides:       t.Sides,
		Orientation: t.Orientation,
		Monastery:   t.Monastery,
	}
	if t.Features != nil {
		clone.Features = make([]Feature, len(t.Features))

		for i, feature := range t.Features {
			clone.Features[i] = feature
			clone.Features[i].Region = nil

			if feature.Sides != nil {
				clone.Features[i].Sides = make([]Side, len(feature.Sides))
				copy(clone.Features[i].Sides, feature.Sides)
			}
		}
	}
	return clone
}

func LoadTiles() []*Tile {
	file, err := os.Open("assets/data/tile_info.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var tiles []*Tile

	err = json.NewDecoder(file).Decode(&tiles)
	if err != nil {
		fmt.Println("Unmarshal error= &err", err)
		panic(err)
	}
	return tiles
}

func (t *Tile) SideAt(direction Direction) SideType {
	return t.Sides[(int(direction)-int(t.Orientation)+4)%4]
}

func (t *Tile) GetSide(direction Direction) SideType {
	return t.SideAt(direction)
}

func (t *Tile) Rotate() {
	t.Orientation = (t.Orientation + 1) % 4
}

func (t *Tile) FeatureByDirection(direction Direction) *Feature {
	rotatedDir := Direction((int(direction) - int(t.Orientation) + 4) % 4)
	for i := range t.Features {
		if t.Features[i].HasSide(rotatedDir) {
			return &t.Features[i]
		}
	}
	return nil
}
