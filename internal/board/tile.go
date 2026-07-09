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
	Sides    [4]SideType `json:"sides"`
	Features []Feature   `json:"features"`
}

func (t *Tile) Clone() *Tile {
	clone := &Tile{
		ID:          t.ID,
		Texture:     t.Texture,
		Sides:       t.Sides,
		Orientation: t.Orientation,
	}
	if t.Features != nil {
		clone.Features = make([]Feature, len(t.Features))

		for i, feature := range t.Features {
			clone.Features[i] = feature

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
	return t.Sides[direction.Reset(t.Orientation)]
}

func (t *Tile) GetSide(direction Direction) SideType {
	return t.SideAt(direction)
}

func (t *Tile) Rotate() {
	t.Orientation = (t.Orientation + 1) % 4
}

func (t *Tile) FeatureByDirection(dir Direction) (*Feature, int) {
	rotatedDir := dir.Reset(t.Orientation)
	for i := range t.Features {
		if _, exists := t.Features[i].GetSide(rotatedDir); exists {
			return &t.Features[i], i
		}
	}
	panic(fmt.Sprintf("Feature not present in Tile %d by Direction %d", t.ID, dir))
}

func (t *Tile) GetFeatureSide(dir Direction) (*Side, bool) {
	for i := range t.Features {
		if side, exists := t.Features[i].GetSide(dir); exists {
			return side, true
		}
	}
	return nil, false
}

func (t *Tile) CompleteSide(dir Direction) {
	resetDir := dir.Reset(t.Orientation)
	side, exist := t.GetFeatureSide(resetDir)
	if !exist {
		panic(fmt.Sprintf("Feature Side at %s in Tile %d does not exist!", resetDir, t.ID))
	}
	side.Complete = true
}

func (t *Tile) UpdateRegionId(featId int, regId RegionID) {
	t.Features[featId].RegionID = regId
}
