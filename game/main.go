package main

import (
	"KarkasonFoollery/internal/board"
	"KarkasonFoollery/internal/ui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	board1 := board.NewBoard()
	assets, err := ui.LoadAssets()
	if err != nil {
		log.Fatal(err)
	}
	tiles := board.LoadTiles()
	deck := board.NewDeck(tiles)

	game := ui.NewGame(board1, assets, deck)
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Karkason Foollery")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
