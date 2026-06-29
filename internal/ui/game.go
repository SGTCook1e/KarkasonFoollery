package ui

import (
	"KarkasonFoollery/internal/board"
	"KarkasonFoollery/internal/game"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type turnPhase string

const (
	AwaitingTile       turnPhase = "AwaitingTile"
	AwaitingMeeple     turnPhase = "AwaitingMeeple"
	ResolvingPlacement turnPhase = "ResolvingPlacement"
)

const tileSize = 256

const slotRadius = 35

type Game struct {
	state  *game.GameState
	phase  turnPhase
	assets *Assets

	cameraX, cameraY float64
	zoom             float64

	mousePressed   bool
	mouseX, mouseY int

	cameraSpeed float64

	hoverX, hoverY int

	rotPressed bool
}

func NewGame(state *game.GameState, assets *Assets) *Game {
	TRACKED = state //
	return &Game{
		state:        state,
		phase:        AwaitingTile,
		assets:       assets,
		mousePressed: false,
		cameraSpeed:  10,
		//1=100% zoom
		//0.5=50% zoom
		//2.0=200% zoom
		zoom: 1.0,
	}
}

func (g *Game) Update() error {
	g.updateHover()
	g.updateCamera()
	g.updateRotation()

	switch g.phase {
	case AwaitingTile:
		g.handleTilePlacementInput()
	case AwaitingMeeple:
		g.handleMeeplePlacementInput()
	case ResolvingPlacement:
		g.handlePlacementResolve()
	}

	return nil
}

func (g *Game) worldBounds(screenW, screenH int) (float64, float64, float64, float64) {
	invZoom := 1.0 / g.zoom

	left := g.cameraX
	top := g.cameraY

	right := g.cameraX + float64(screenW)*invZoom
	bottom := g.cameraY + float64(screenH)*invZoom

	return left, top, right, bottom
}

func (g *Game) calcFeatureCoords(worldX, worldY float64, f board.Feature, t board.Tile) (fx, fy float64) {
	half := float64(tileSize) / 2.0
	centerX := worldX + half
	centerY := worldY + half

	if f.Type == board.FeatureMonastery {
		return centerX, centerY
	}

	var dirX, dirY float64
	for _, side := range f.Sides {
		switch side.Direction.Rotate(t.Orientation) {
		case board.Top:
			dirY -= 1.0
		case board.Right:
			dirX += 1.0
		case board.Bottom:
			dirY += 1.0
		case board.Left:
			dirX -= 1.0
		}
	}

	shiftDistance := float64(tileSize) * 0.28
	featX := centerX + (dirX * shiftDistance)
	featY := centerY + (dirY * shiftDistance)

	return featX, featY
}

func (g *Game) updateHover() {
	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToGridFloor(mx, my)

	g.hoverX = worldX
	g.hoverY = worldY
}

func (g *Game) updateCamera() {
	speed := g.cameraSpeed / g.zoom
	_, wheelY := ebiten.Wheel()

	if wheelY != 0 {
		g.zoomAtCursor(wheelY)
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.cameraY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.cameraY += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.cameraX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.cameraX += speed
	}
}

func (g *Game) zoomAtCursor(wheelY float64) {
	mx, my := ebiten.CursorPosition()

	worldX := float64(mx)/g.zoom + g.cameraX
	worldY := float64(my)/g.zoom + g.cameraY

	if wheelY > 0 {
		g.zoom *= 1.1
	} else {
		g.zoom *= 0.9
	}

	if g.zoom < 0.2 {
		g.zoom = 0.2
	}
	if g.zoom > 4.0 {
		g.zoom = 4.0
	}

	g.cameraX = worldX - float64(mx)/g.zoom
	g.cameraY = worldY - float64(my)/g.zoom
}

func (g *Game) handleTilePlacementInput() {
	if !g.consumeLeftClick() {
		return
	}

	coord := g.cursorCoord()

	if _, exists := g.state.Board.GetTile(coord); exists {
		return
	}

	if !g.state.Board.IsValidPlacement(coord, g.state.TopTile) {
		return
	}

	tile := g.state.TopTile.Clone()

	g.state.Board.PlaceTile(coord, tile)
	g.state.CurrCoord = coord
	g.state.TopTile = g.state.Deck.Draw()

	g.phase = AwaitingMeeple
}

func (g *Game) handleMeeplePlacementInput() {
	tile, exists := g.state.Board.GetTile(g.state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", g.state.CurrCoord))
	}

	if !tile.HasFeatureTypes(board.FeatureCity, board.FeatureRoad, board.FeatureMonastery) {
		g.phase = ResolvingPlacement
		return
	}

	if !g.consumeLeftClick() {
		return
	}

	worldX := float64(g.state.CurrCoord.X * tileSize)
	worldY := float64(g.state.CurrCoord.Y * tileSize)

	featId, ok := g.getClickedFeatureId(worldX, worldY, *tile)
	if !ok {
		return
	}

	feature := tile.Features[featId]
	feature.Meeple = board.Peasant

	g.phase = ResolvingPlacement
}

func (g *Game) getClickedFeatureId(worldX, worldY float64, tile board.Tile) (int, bool) {
	mx, my := ebiten.CursorPosition()

	for index, feature := range tile.Features {
		if feature.Type != board.FeatureCity &&
			feature.Type != board.FeatureRoad &&
			feature.Type != board.FeatureMonastery {
			continue
		}
		fx, fy := g.calcFeatureCoords(worldX, worldY, feature, tile)
		sx, sy := g.worldToScreen(fx, fy)
		if inRadius(sx, sy, slotRadius, float64(mx), float64(my)) {
			return index, true
		}
	}
	return 0, false
}

func inRadius(centerX, centerY, radius, mouseX, mouseY float64) bool {
	dx := mouseX - centerX
	dy := mouseY - centerY

	sqDistance := (dx * dx) + (dy * dy)
	sqRadius := radius * radius

	if sqDistance <= sqRadius {
		return true
	}
	return false
}

func (g *Game) handlePlacementResolve() {
	result := game.ResolvePlacement(*g.state, 1)
	g.state.ApplyPlacement(result, 1)

	g.phase = AwaitingTile
}

func (g *Game) cursorCoord() board.Coord {
	mx, my := ebiten.CursorPosition()
	x, y := g.screenToGridFloor(mx, my)

	return board.Coord{X: x, Y: y}
}

func (g *Game) screenToGridFloor(x, y int) (int, int) {
	wx, wy := g.screenToWorld(x, y)

	return int(math.Floor(wx / tileSize)), int(math.Floor(wy / tileSize))
}

func (g *Game) consumeLeftClick() bool {
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if !pressed {
		g.mousePressed = false
		return false
	}

	if g.mousePressed {
		return false
	}

	g.mousePressed = true
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	return true
}

func (g *Game) updateRotation() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if !g.rotPressed {
			g.state.TopTile.Rotate()
			g.rotPressed = true
		}
	} else {
		g.rotPressed = false
	}
}

func (g *Game) Layout(outerWidth, outerHeight int) (int, int) {
	return 1280, 720
}

func (g *Game) worldToScreen(x, y float64) (float64, float64) {
	x -= g.cameraX
	y -= g.cameraY
	return x * g.zoom, y * g.zoom
}

func (g *Game) screenToWorld(x, y int) (float64, float64) {
	return float64(x)/g.zoom + g.cameraX,
		float64(y)/g.zoom + g.cameraY
}

func (g *Game) screenToGrid(x, y int) (int, int) {
	wx, wy := g.screenToWorld(x, y)
	return int(wx / tileSize), int(wy / tileSize)
}

func (g *Game) worldToDrawOpts(worldX, worldY float64, rotation int) ebiten.DrawImageOptions {
	opts := ebiten.DrawImageOptions{}

	// Apply camera translation first (same order as in Draw)
	opts.GeoM.Translate(-g.cameraX, -g.cameraY)

	// Then place to world position, rotate around center, and scale
	opts.GeoM.Translate(worldX, worldY)
	opts.GeoM.Translate(float64(tileSize)/2, float64(tileSize)/2)
	opts.GeoM.Rotate(float64(rotation) * math.Pi / 2)
	opts.GeoM.Translate(-float64(tileSize)/2, -float64(tileSize)/2)

	opts.GeoM.Scale(g.zoom, g.zoom)

	return opts
}
