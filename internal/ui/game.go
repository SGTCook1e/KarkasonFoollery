package ui

import (
	"KarkasonFoollery/internal/board"
	"KarkasonFoollery/internal/tile"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const tileSize = 256

type Game struct {
	board  *board.Board
	assets *Assets

	cameraX, cameraY float64
	zoom             float64

	mousePressed   bool
	mouseX, mouseY int

	cameraSpeed float64

	curentTile *tile.Tile
	// curentTileID int
	rotPressed bool

	hoverX, hoverY int

	// tiles []*tile.Tile
	deck *tile.Deck
}

func NewGame(b *board.Board, assets *Assets, deck *tile.Deck) *Game {
	firstTile := deck.Draw()

	return &Game{
		board:      b,
		assets:     assets,
		deck:       deck,
		curentTile: firstTile,
		// curentTileID: 0,
		mousePressed: false,
		cameraSpeed:  10,
		//1=100% zoom
		//0.5=50% zoom
		//2.0=200% zoom
		zoom: 1.0,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawGrid(screen)

	for coords, t := range g.board.GetTiles() {
		img, ok := g.assets.Tiles[t.Texture[:len(t.Texture)-4]]
		if !ok {
			continue
		}

		worldX := float64(coords.X * tileSize)
		worldY := float64(coords.Y * tileSize)

		opts := &ebiten.DrawImageOptions{}
		// Rotate around the tile center in local space
		half := float64(tileSize) / 2
		opts.GeoM.Translate(-half, -half)
		opts.GeoM.Rotate(float64(t.Dir) * math.Pi / 2)
		opts.GeoM.Translate(half, half)

		// Scale and then translate to screen coordinates
		opts.GeoM.Scale(g.zoom, g.zoom)
		sxScreen, syScreen := g.worldToScreen(worldX, worldY)
		opts.GeoM.Translate(sxScreen, syScreen)

		screen.DrawImage(img, opts)
		//g.drawTileSideLabels(screen, t, worldX, worldY)
	}

	// Draw preview last so it's on top of placed tiles
	g.drawPreview(screen)
}

func (g *Game) worldBounds(screenW, screenH int) (float64, float64, float64, float64) {
	invZoom := 1.0 / g.zoom

	left := g.cameraX
	top := g.cameraY

	right := g.cameraX + float64(screenW)*invZoom
	bottom := g.cameraY + float64(screenH)*invZoom

	return left, top, right, bottom
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	left, top, right, bottom := g.worldBounds(1280, 720)

	startX := int(left) / tileSize * tileSize
	endX := int(right)/tileSize*tileSize + tileSize

	startY := int(top) / tileSize * tileSize
	endY := int(bottom)/tileSize*tileSize + tileSize

	for x := startX; x <= endX; x += tileSize {
		g.drawLine(screen,
			float64(x), top,
			float64(x), bottom,
		)
	}

	for y := startY; y <= endY; y += tileSize {
		g.drawLine(screen,
			left, float64(y),
			right, float64(y),
		)
	}
}

func (g *Game) drawPreview(screen *ebiten.Image) {
	if g.curentTile == nil || g.curentTile.Texture == "" {
		return
	}

	img, ok := g.assets.Tiles[g.curentTile.Texture[:len(g.curentTile.Texture)-4]]
	if !ok {
		return
	}

	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToWorld(mx, my)

	cellX := math.Floor(worldX / tileSize)
	cellY := math.Floor(worldY / tileSize)
	// левый верхний угол клетки в WORLD координатах (привязка к сетке)
	worldX = cellX * tileSize
	worldY = cellY * tileSize

	// Rotate around the tile center in local space, then scale and position
	half := float64(tileSize) / 2
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-half, -half)
	opts.GeoM.Rotate(float64(g.curentTile.Dir) * math.Pi / 2)
	opts.GeoM.Translate(half, half)

	opts.GeoM.Scale(g.zoom, g.zoom)
	sxScreen, syScreen := g.worldToScreen(worldX, worldY)
	opts.GeoM.Translate(sxScreen, syScreen)

	opts.ColorScale.ScaleAlpha(0.5)
	screen.DrawImage(img, &opts)
	//g.drawTileSideLabels(screen, g.curentTile, worldX, worldY)

}

func (g *Game) drawTileSideLabels(screen *ebiten.Image, t *tile.Tile, worldX, worldY float64) {
	if g.zoom < 0.25 {
		return
	}

	offsets := map[tile.Direction]struct{ x, y float64 }{
		tile.Top:    {x: tileSize*0.5 - 28, y: 8},
		tile.Right:  {x: tileSize - 92, y: tileSize*0.5 - 8},
		tile.Bottom: {x: tileSize*0.5 - 28, y: tileSize - 24},
		tile.Left:   {x: 8, y: tileSize*0.5 - 8},
	}

	for _, dir := range []tile.Direction{tile.Top, tile.Right, tile.Bottom, tile.Left} {
		label := string(t.SideAt(dir))
		sx, sy := g.worldToScreen(worldX+offsets[dir].x, worldY+offsets[dir].y)
		ebitenutil.DebugPrintAt(screen, label, int(math.Round(sx)), int(math.Round(sy)))
	}
}

func (g *Game) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	sx1, sy1 := g.worldToScreen(x1, y1)
	sx2, sy2 := g.worldToScreen(x2, y2)

	ebitenutil.DrawLine(screen, sx1, sy1, sx2, sy2, color.RGBA{R: 51, G: 51, B: 51, A: 172})
}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	worldX, worldY := g.screenToWorld(mx, my)

	boardX := int(math.Floor(worldX / tileSize))
	boardY := int(math.Floor(worldY / tileSize))

	g.hoverX = boardX
	g.hoverY = boardY
	speed := g.cameraSpeed / g.zoom
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if pressed && !g.mousePressed {
		g.mousePressed = true
		g.mouseX, g.mouseY = ebiten.CursorPosition()

		worldX, worldY := g.screenToWorld(g.mouseX, g.mouseY)
		boardX := int(math.Floor(worldX / tileSize))
		boardY := int(math.Floor(worldY / tileSize))
		coord := board.Coord{
			X: boardX,
			Y: boardY,
		}
		_, exists := g.board.GetTile(coord)
		if !exists {
			if g.board.IsValidPlacement(coord, g.curentTile) {
				g.board.PlaceTile(coord, g.curentTile.Clone())
				g.curentTile = g.deck.Draw()
			}
		}
	}
	if !pressed {
		g.mousePressed = false
	}

	_, wheelY := ebiten.Wheel()

	if wheelY != 0 {
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
		return nil
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

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if !g.rotPressed {
			g.curentTile.Rotate()
			g.rotPressed = true
		}
	} else {
		g.rotPressed = false
	}

	return nil
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
