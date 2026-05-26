package tile

import "math/rand/v2"

type Deck struct {
	tiles []*Tile
}

func NewDeck(tiles []*Tile) *Deck {
	var startRiver *Tile
	var endRiver *Tile

	var riverTiles []*Tile
	var otherTiles []*Tile

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

	var returnTiles []*Tile

	returnTiles = append(returnTiles, startRiver)
	returnTiles = append(returnTiles, riverTiles...)
	returnTiles = append(returnTiles, endRiver)
	returnTiles = append(returnTiles, otherTiles...)

	return &Deck{
		tiles: returnTiles,
	}
}

func shuffleTiles(tiles []*Tile) {
	// Implement a shuffling algorithm, such as Fisher-Yates
	for i := len(tiles) - 1; i > 0; i-- {
		j := rand.IntN(i)
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}
}

func (d *Deck) Draw() *Tile {
	if len(d.tiles) == 0 {
		return &Tile{}
	}
	t := d.tiles[0]
	d.tiles = d.tiles[1:]
	return t.Clone()
}
