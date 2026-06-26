package board

import (
	"fmt"
)

type Board struct {
	tiles map[Coord]*Tile
}

func (b *Board) GetTiles() map[Coord]*Tile {
	return b.tiles
}

func NewBoard() *Board {
	return &Board{
		tiles: make(map[Coord]*Tile),
	}
}

func (b *Board) Clone() Board {
	clone := Board{
		tiles: make(map[Coord]*Tile, len(b.tiles)),
	}

	for coord, tile := range b.tiles {
		if tile != nil {
			clone.tiles[coord] = tile.Clone()
		}
	}

	return clone
}

func (b *Board) PlaceTile(c Coord, t *Tile) {
	b.tiles[c] = t
}

func (b *Board) GetTile(c Coord) (*Tile, bool) {
	t, exists := b.tiles[c]
	return t, exists
}

func (b *Board) PrintBoard() {
	for c, t := range b.tiles {
		fmt.Printf("Tile ID: %d at coord (%d, %d) [%v]\n", t.ID, c.X, c.Y, t.Sides)
	}
}

func (b *Board) CanPlaceTile(c Coord) bool {
	if len(b.tiles) == 0 {
		return true
	}
	if b.tiles[c] != nil {
		return false
	} else {
		// Check adjacent tiles
		if b.tiles[Coord{X: c.X, Y: c.Y - 1}] != nil {
			return true
		}
		if b.tiles[Coord{X: c.X + 1, Y: c.Y}] != nil {
			return true
		}
		if b.tiles[Coord{X: c.X, Y: c.Y + 1}] != nil {
			return true
		}
		if b.tiles[Coord{X: c.X - 1, Y: c.Y}] != nil {
			return true
		}
	}
	return false
}

func (b *Board) IsValidPlacement(c Coord, t Tile) bool {
	if len(b.tiles) == 0 {
		return true
	}
	if b.tiles[c] != nil {
		return false
	}

	validNeighbor := false

	if top, ok := b.tiles[Coord{X: c.X, Y: c.Y - 1}]; ok {
		if top.GetSide(Bottom) != t.GetSide(Top) {
			return false
		}
		validNeighbor = true
	}
	if right, ok := b.tiles[Coord{X: c.X + 1, Y: c.Y}]; ok {
		if right.GetSide(Left) != t.GetSide(Right) {
			return false
		}
		validNeighbor = true
	}
	if bottom, ok := b.tiles[Coord{X: c.X, Y: c.Y + 1}]; ok {
		if bottom.GetSide(Top) != t.GetSide(Bottom) {
			return false
		}
		validNeighbor = true
	}
	if left, ok := b.tiles[Coord{X: c.X - 1, Y: c.Y}]; ok {
		if left.GetSide(Right) != t.GetSide(Left) {
			return false
		}
		validNeighbor = true
	}

	return validNeighbor
}
