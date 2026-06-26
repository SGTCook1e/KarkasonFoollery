package ui

import (
	"KarkasonFoollery/internal/board"
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawGrid(screen)

	for coords, t := range g.state.Board.GetTiles() {
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
		opts.GeoM.Rotate(float64(t.Orientation) * math.Pi / 2)
		opts.GeoM.Translate(half, half)

		// Scale and then translate to screen coordinates
		opts.GeoM.Scale(g.zoom, g.zoom)
		sxScreen, syScreen := g.worldToScreen(worldX, worldY)
		opts.GeoM.Translate(sxScreen, syScreen)

		screen.DrawImage(img, opts)
		g.drawTileSideLabels(screen, *t, worldX, worldY) //
		g.drawRegionMarkers(screen, *t, worldX, worldY)  //}
	}

	// Draw preview last so it's on top of placed tiles
	if g.phase == AwaitingTile {
		g.drawTilePreview(screen)
	}
	if g.phase == AwaitingMeeple {
		g.drawMeepleSlots(screen)
	}
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

func (g *Game) drawTilePreview(screen *ebiten.Image) {
	if g.state.TopTile.Texture == "" {
		return
	}

	img, ok := g.assets.Tiles[g.state.TopTile.Texture[:len(g.state.TopTile.Texture)-4]]
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
	opts.GeoM.Rotate(float64(g.state.TopTile.Orientation) * math.Pi / 2)
	opts.GeoM.Translate(half, half)

	opts.GeoM.Scale(g.zoom, g.zoom)
	sxScreen, syScreen := g.worldToScreen(worldX, worldY)
	opts.GeoM.Translate(sxScreen, syScreen)

	opts.ColorScale.ScaleAlpha(0.5)
	screen.DrawImage(img, &opts)
	//g.drawTileSideLabels(screen, g.curentTile, worldX, worldY)
}

func (g *Game) drawTileSideLabels(screen *ebiten.Image, t board.Tile, worldX, worldY float64) {
	if g.zoom < 0.25 {
		return
	}

	offsets := map[board.Direction]struct{ x, y float64 }{
		board.Top:    {x: tileSize*0.5 - 28, y: 8},
		board.Right:  {x: tileSize - 92, y: tileSize*0.5 - 8},
		board.Bottom: {x: tileSize*0.5 - 28, y: tileSize - 24},
		board.Left:   {x: 8, y: tileSize*0.5 - 8},
	}

	for _, dir := range []board.Direction{board.Top, board.Right, board.Bottom, board.Left} {
		label := string(t.SideAt(dir))
		label = fmt.Sprintf("%s%d", label, t.ID)
		sx, sy := g.worldToScreen(worldX+offsets[dir].x, worldY+offsets[dir].y)
		ebitenutil.DebugPrintAt(screen, label, int(math.Round(sx)), int(math.Round(sy)))
	}
}

func getColorFromId(id int) color.Color {
	src := rand.NewSource(int64(id))
	rnd := rand.New(src)

	return color.RGBA{
		R: uint8(rnd.Intn(156) + 100),
		G: uint8(rnd.Intn(156) + 100),
		B: uint8(rnd.Intn(156) + 100),
		A: 240,
	}
}

func (g *Game) drawRegionMarkers(screen *ebiten.Image, t board.Tile, worldX, worldY float64) {
	for _, feature := range t.Features {
		if feature.RegionID == board.NoRegion {
			continue
		}

		fx, fy := g.calcFeatureCoords(worldX, worldY, feature, t)
		sx, sy := g.worldToScreen(fx, fy)

		regionColor := getColorFromId(int(feature.RegionID))
		scaledSize := float32(40 * g.zoom)

		sx = sx - 0.5*float64(scaledSize)
		sy = sy - 0.5*float64(scaledSize)

		vector.FillRect(screen, float32(sx), float32(sy), scaledSize, scaledSize, regionColor, true)
	}
}

func (g *Game) drawMeepleSlots(screen *ebiten.Image) {
	tile, exists := g.state.Board.GetTile(g.state.CurrCoord)
	if !exists {
		panic(fmt.Sprintf("Tile at %+v does not exist!", g.state.CurrCoord))
	}

	worldX := float64(g.state.CurrCoord.X * tileSize)
	worldY := float64(g.state.CurrCoord.Y * tileSize)
	color := color.RGBA{R: 120, G: 120, B: 120, A: 128}
	scaledRad := float32(40 * g.zoom)

	if tile.Monastery {
		mx := worldX + tileSize*0.5
		my := worldY + tileSize*0.5

		sx, sy := g.worldToScreen(mx, my)
		vector.FillCircle(screen, float32(sx), float32(sy), scaledRad, color, true)
	}

	for _, feature := range tile.Features {
		if feature.Type == board.FeatureField || feature.Type == board.FeatureRiver {
			continue
		}
		fx, fy := g.calcFeatureCoords(worldX, worldY, feature, *tile)
		sx, sy := g.worldToScreen(fx, fy)

		vector.FillCircle(screen, float32(sx), float32(sy), scaledRad, color, true)
	}
}

func (g *Game) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	sx1, sy1 := g.worldToScreen(x1, y1)
	sx2, sy2 := g.worldToScreen(x2, y2)

	ebitenutil.DrawLine(screen, sx1, sy1, sx2, sy2, color.RGBA{R: 51, G: 51, B: 51, A: 172})
}
