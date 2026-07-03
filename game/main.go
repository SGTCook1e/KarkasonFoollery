package main

import (
	"KarkasonFoollery/internal/board"
	"KarkasonFoollery/internal/game"
	"KarkasonFoollery/internal/ui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	assets, err := ui.LoadAssets()
	if err != nil {
		log.Fatal(err)
	}
	tiles := board.LoadTiles()

	state := game.NewState(tiles)
	state.NextPlayerId = state.AddPlayer()
	state.NextPlayerId = state.AddPlayer()
	game := ui.NewGame(state, assets)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Karkason Foollery")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
