package tile

import (
	"encoding/json"
	"os"
)

type Tile struct {
	//ID 1-12 is Rivers
	//ID 1 is the starting river tile
	//ID 12 is end river tile
	ID      int
	Texture string

	// 0 = 0°
	// 1 = 90°
	// 2 = 180°
	// 3 = 270°
	Rotation int
	// Sides are stored in tile order: Top, Right, Bottom, Left.
	Sides [4]SideType
	// Features []Feature
}

type TileData struct {
	ID      int
	Texture string
	Sides   []string
}

func NewTile(id int, path string, sides [4]SideType) *Tile {
	return &Tile{
		ID:      id,
		Texture: path,
		Sides:   sides,
		// Features: features,
		Rotation: 0,
	}
}

func (t *Tile) Clone() *Tile {
	return &Tile{
		ID:      t.ID,
		Texture: t.Texture,
		Sides:   t.Sides,
		// Features: t.Features,
		Rotation: t.Rotation,
	}
}

func loadTileData(path string) ([]TileData, error) {
	var tileData []TileData
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(file).Decode(&tileData)
	if err != nil {
		return nil, err
	}

	return tileData, nil
}

func LoadTiles() []*Tile {
	var tiles []*Tile
	tileData, err := loadTileData("assets/data/tile_info.json")
	if err != nil {
		panic(err)
	}
	for _, data := range tileData {
		sides := [4]SideType{}
		for i, side := range data.Sides {
			switch side {
			case "field":
				sides[i] = Field
			case "road":
				sides[i] = Road
			case "city":
				sides[i] = City
			case "river":
				sides[i] = River
			}
		}
		tiles = append(tiles, NewTile(data.ID, data.Texture, sides))
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
	return t.Sides[(int(direction)-t.Rotation+4)%4]
}

func (t *Tile) GetSide(direction Direction) SideType {
	return t.SideAt(direction)
}

func (t *Tile) Rotate() {
	t.Rotation = (t.Rotation + 1) % 4
}

func (s SideType) String() string {
	switch s {
	case Field:
		return "Field"
	case Road:
		return "Road"
	case City:
		return "City"
	case River:
		return "River"
	default:
		return "Unknown"
	}
}
