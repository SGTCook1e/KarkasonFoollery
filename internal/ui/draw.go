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
		g.drawRegionMarkers(screen, *t, worldX, worldY)  //
		g.drawMeeples(screen, *t, worldX, worldY)
	}

	// Draw preview last so it's on top of placed tiles
	if g.phase == AwaitingTile || g.phase == AwaitingMeeple {
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
	if g.state.CurrTile.Texture == "" {
		return
	}

	img, ok := g.assets.Tiles[g.state.CurrTile.Texture[:len(g.state.CurrTile.Texture)-4]]
	if !ok {
		return
	}

	var coord board.Coord
	if g.phase == AwaitingTile {
		coord = g.cursorCoord()
	} else if g.phase == AwaitingMeeple {
		coord = g.state.CurrCoord
	}

	worldX := float64(coord.X * tileSize)
	worldY := float64(coord.Y * tileSize)

	// Rotate around the tile center in local space, then scale and position
	half := float64(tileSize) / 2
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-half, -half)
	opts.GeoM.Rotate(float64(g.state.CurrTile.Orientation) * math.Pi / 2)
	opts.GeoM.Translate(half, half)

	opts.GeoM.Scale(g.zoom, g.zoom)
	sxScreen, syScreen := g.worldToScreen(worldX, worldY)
	opts.GeoM.Translate(sxScreen, syScreen)

	opts.ColorScale.ScaleAlpha(0.5)
	screen.DrawImage(img, &opts)
	//g.drawTileSideLabels(screen, g.curentTile, worldX, worldY)
}

func (g *Game) drawMeeples(screen *ebiten.Image, t board.Tile, worldX, worldY float64) {
	for _, feature := range t.Features {
		if feature.Meeple.Type == board.NoUnit {
			continue
		}

		fx, fy := g.calcFeatureCoords(worldX, worldY, feature, t)
		sx, sy := g.worldToScreen(fx, fy)

		label := fmt.Sprintf("Meeple owner: %d", feature.Meeple.Owner)
		ebitenutil.DebugPrintAt(screen, label, int(math.Round(sx)), int(math.Round(sy))) //

		// img := g.assets.Meeples[string(feature.Meeple)]
		// screen.DrawImage(img, &opts)
	}
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
		resetDir := dir.Reset(t.Orientation)
		side, exist := t.GetFeatureSide(resetDir)
		if !exist {
			panic(fmt.Sprintf("Feature Side at %s in Tile %d does not exist!", resetDir, t.ID))
		}
		sideType := string(t.SideAt(dir))
		label := fmt.Sprintf("%s%d%t", sideType, t.ID, side.Complete)

		sx, sy := g.worldToScreen(worldX+offsets[dir].x, worldY+offsets[dir].y)
		ebitenutil.DebugPrintAt(screen, label, int(math.Round(sx)), int(math.Round(sy)))
	}
}

func getColorFromId(id int) color.RGBA {
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
		reg := g.state.Regions.ByID[feature.RegionID]
		if feature.RegionID == board.NoRegion {
			continue
		}

		fx, fy := g.calcFeatureCoords(worldX, worldY, feature, t)
		sx, sy := g.worldToScreen(fx, fy)

		// regionColor := getColorFromId(int(feature.RegionID))
		regionColor := getColorFromId(int(reg.Owner))
		if reg.Owner == board.NoOwner {
			regionColor = color.RGBA{0, 0, 0, 100}
		}
		if reg.Contested {
			regionColor = color.RGBA{255, 0, 0, 100}
		}

		scaledSize := float32(40 * g.zoom)
		sx = sx - 0.5*float64(scaledSize)
		sy = sy - 0.5*float64(scaledSize)

		vector.FillRect(screen, float32(sx), float32(sy), scaledSize, scaledSize, regionColor, true)
		label := fmt.Sprintf("%d:%d", reg.Owner, reg.ID)
		ebitenutil.DebugPrintAt(screen, label, int(math.Round(sx)), int(math.Round(sy)))
	}
}

func (g *Game) drawMeepleSlots(screen *ebiten.Image) {
	tile := g.state.CurrTile

	worldX := float64(g.state.CurrCoord.X * tileSize)
	worldY := float64(g.state.CurrCoord.Y * tileSize)
	slotColor := color.RGBA{R: 120, G: 120, B: 120, A: 128}
	backlightColor := color.RGBA{R: 140, G: 140, B: 140, A: 128}
	scaledRad := float32(slotRadius * g.zoom)

	for index, feature := range tile.Features {
		if !g.state.IsValidMeeplePlacement(index) {
			continue
		}
		fx, fy := g.calcFeatureCoords(worldX, worldY, feature, *tile)
		sx, sy := g.worldToScreen(fx, fy)

		mx, my := ebiten.CursorPosition()
		if inRadius(sx, sy, slotRadius, float64(mx), float64(my)) {
			vector.FillCircle(screen, float32(sx), float32(sy), scaledRad, backlightColor, true)
		}

		vector.FillCircle(screen, float32(sx), float32(sy), scaledRad, slotColor, true)
	}
}

func (g *Game) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	sx1, sy1 := g.worldToScreen(x1, y1)
	sx2, sy2 := g.worldToScreen(x2, y2)

	ebitenutil.DrawLine(screen, sx1, sy1, sx2, sy2, color.RGBA{R: 51, G: 51, B: 51, A: 172})
}
