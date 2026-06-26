package game

import (
	b "KarkasonFoollery/internal/board"
	"math/rand/v2"
)

type Deck struct {
	tiles []*b.Tile
}

func NewDeck(tiles []*b.Tile) *Deck {
	var startRiver *b.Tile
	var endRiver *b.Tile

	var riverTiles []*b.Tile
	var otherTiles []*b.Tile

	startRiver = tiles[0].Clone()
	endRiver = tiles[11].Clone()
	for i := 1; i < 11; i++ {
		riverTiles = append(riverTiles, tiles[i].Clone())
	}
	for i := 12; i < len(tiles); i++ {
		otherTiles = append(otherTiles, tiles[i].Clone())
	}
	shuffleTiles(riverTiles)
	shuffleTiles(otherTiles)

	var returnTiles []*b.Tile

	returnTiles = append(returnTiles, startRiver)
	returnTiles = append(returnTiles, riverTiles...)
	returnTiles = append(returnTiles, endRiver)
	returnTiles = append(returnTiles, otherTiles...)

	return &Deck{
		tiles: returnTiles,
	}
}

func shuffleTiles(tiles []*b.Tile) {
	// Implement a shuffling algorithm, such as Fisher-Yates
	for i := len(tiles) - 1; i > 0; i-- {
		j := rand.IntN(i)
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}
}

func (d *Deck) Draw() b.Tile {
	if len(d.tiles) == 0 {
		return b.Tile{}
	}
	t := d.tiles[0]
	d.tiles = d.tiles[1:]
	return *t
}
