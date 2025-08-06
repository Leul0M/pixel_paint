package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/* ---------- constants ---------- */

const (
	iconSize      = 32
	iconMargin    = 8
	paletteCols   = 8
	paletteRows   = 2
	swatchSize    = 24
	palettePad    = 4

	// Top bar height (toolbar + palette)
	topBarHeight = iconSize + 2*iconMargin + paletteRows*(swatchSize+palettePad) + palettePad
	paletteTop   = iconSize + 2*iconMargin + palettePad
)

/* ---------- game state ---------- */

type Game struct {
	width, height int

	mode            string // "" | "brush" | "eraser"
	currentColor    color.Color
	isDrawing       bool
	leftPressedLast bool
	prevX, prevY    int

	background *ebiten.Image // solid white underlay
	drawing    *ebiten.Image // user strokes

	// toolbar rectangles
	brushIcon  image.Rectangle
	eraserIcon image.Rectangle
	clearIcon  image.Rectangle
	saveIcon   image.Rectangle

	// palette
	palette []color.Color
}

/* ---------- constructor ---------- */

func NewGame(w, h int) *Game {
	bg := ebiten.NewImage(w, h)
	bg.Fill(color.White)
	drw := ebiten.NewImage(w, h)
	drw.Fill(color.White)

	g := &Game{
		width:         w,
		height:        h,
		background:    bg,
		drawing:       drw,
		currentColor:  color.Black,
	}

	// fixed 16-colour palette
	g.palette = []color.Color{
		color.Black,
		color.RGBA{127, 0, 0, 255},
		color.RGBA{0, 127, 0, 255},
		color.RGBA{0, 0, 127, 255},
		color.RGBA{127, 127, 0, 255},
		color.RGBA{127, 0, 127, 255},
		color.RGBA{0, 127, 127, 255},
		color.RGBA{127, 127, 127, 255},
		color.White,
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 255, 255},
		color.RGBA{0, 255, 255, 255},
		color.RGBA{255, 255, 255, 255},
	}

	// toolbar icon positions
	g.brushIcon = image.Rect(iconMargin, iconMargin,
		iconMargin+iconSize, iconMargin+iconSize)
	g.eraserIcon = image.Rect(2*iconMargin+iconSize, iconMargin,
		2*iconMargin+2*iconSize, iconMargin+iconSize)
	g.clearIcon = image.Rect(3*iconMargin+2*iconSize, iconMargin,
		3*iconMargin+3*iconSize, iconMargin+iconSize)
	g.saveIcon = image.Rect(4*iconMargin+3*iconSize, iconMargin,
		4*iconMargin+4*iconSize, iconMargin+iconSize)

	return g
}

/* ---------- game loop ---------- */

func (g *Game) Update() error {
	leftPressedNow := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	justPressed := leftPressedNow && !g.leftPressedLast
	g.leftPressedLast = leftPressedNow

	if justPressed {
		g.handleClick()
	}

	if g.isDrawing && leftPressedNow {
		g.continueStroke()
	} else {
		g.isDrawing = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, nil)
	screen.DrawImage(g.drawing, nil)

	// Draw top bar background (toolbar + palette)
	topBar := ebiten.NewImage(g.width, topBarHeight)
	topBar.Fill(color.RGBA{220, 220, 220, 255})
	screen.DrawImage(topBar, nil)

	g.drawToolbar(screen)
	g.drawPalette(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth == g.width && outsideHeight == g.height {
		return g.width, g.height
	}

	newBg := ebiten.NewImage(outsideWidth, outsideHeight)
	newBg.Fill(color.White)
	newBg.DrawImage(g.background, nil)

	newDrw := ebiten.NewImage(outsideWidth, outsideHeight)
	newDrw.Fill(color.White)
	newDrw.DrawImage(g.drawing, nil)

	g.background, g.drawing = newBg, newDrw
	g.width, g.height = outsideWidth, outsideHeight
	return g.width, g.height
}

/* ---------- private helpers ---------- */

func (g *Game) handleClick() {
	x, y := ebiten.CursorPosition()

	// toolbar clicks
	if y < iconSize+2*iconMargin {
		switch {
		case image.Pt(x, y).In(g.brushIcon):
			g.mode = "brush"
			return
		case image.Pt(x, y).In(g.eraserIcon):
			g.mode = "eraser"
			return
		case image.Pt(x, y).In(g.clearIcon):
			g.drawing.Fill(color.White)
			return
		case image.Pt(x, y).In(g.saveIcon):
			g.save()
			return
		}
	}

	// palette colour pick
	if y >= paletteTop && y < paletteTop+paletteRows*(swatchSize+palettePad) {
		idx := (y-paletteTop)/(swatchSize+palettePad)*paletteCols +
			(x - palettePad) / (swatchSize + palettePad)
		if idx >= 0 && idx < len(g.palette) {
			g.currentColor = g.palette[idx]
			return
		}
	}

	// start stroke on canvas when a tool is active (only below top bar)
	if g.mode == "brush" || g.mode == "eraser" {
		if y > topBarHeight {
			g.isDrawing = true
			g.prevX, g.prevY = x, y
		}
	}
}

func (g *Game) continueStroke() {
	x, y := ebiten.CursorPosition()

	// prevent drawing inside top bar
	if y <= topBarHeight {
		return
	}

	if x == g.prevX && y == g.prevY {
		return
	}

	col := g.currentColor
	if g.mode == "eraser" {
		col = color.White
	}

	ebitenutil.DrawLine(g.drawing,
		float64(g.prevX), float64(g.prevY),
		float64(x), float64(y), col)
	g.prevX, g.prevY = x, y
}

func (g *Game) save() {
	full := ebiten.NewImage(g.width, g.height)
	full.DrawImage(g.background, nil)
	full.DrawImage(g.drawing, nil)

	f, err := os.Create("output.png")
	if err != nil {
		println("save failed:", err.Error())
		return
	}
	defer f.Close()
	if err := png.Encode(f, full); err != nil {
		println("save failed:", err.Error())
	}
}

func (g *Game) drawToolbar(screen *ebiten.Image) {
	drawIcon := func(rect image.Rectangle, label string, active bool, activeCol color.Color) {
		img := ebiten.NewImage(iconSize, iconSize)
		if active {
			img.Fill(activeCol)
		} else {
			img.Fill(color.White)
		}
		ebitenutil.DrawRect(img, 2, 2, iconSize-4, iconSize-4, color.Black)
		ebitenutil.DebugPrintAt(img, label, iconSize/2-4, iconSize/2-6)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		screen.DrawImage(img, op)
	}

	drawIcon(g.brushIcon, "B", g.mode == "brush", color.RGBA{150, 150, 255, 255})
	drawIcon(g.eraserIcon, "E", g.mode == "eraser", color.RGBA{255, 150, 150, 255})
	drawIcon(g.clearIcon, "C", false, color.White)
	drawIcon(g.saveIcon, "S", false, color.White)
}

func (g *Game) drawPalette(screen *ebiten.Image) {
	// swatches
	for i, col := range g.palette {
		x := palettePad + (i%paletteCols)*(swatchSize+palettePad)
		y := paletteTop + (i/paletteCols)*(swatchSize+palettePad)
		sq := ebiten.NewImage(swatchSize, swatchSize)
		sq.Fill(col)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(sq, op)
	}

	// preview square
	preview := ebiten.NewImage(swatchSize, swatchSize)
	preview.Fill(g.currentColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(palettePad+(paletteCols+1)*(swatchSize+palettePad)), float64(paletteTop))
	screen.DrawImage(preview, op)
}

/* ---------- entry point ---------- */

func main() {
	ebiten.SetWindowSize(1020, 668)
	ebiten.SetWindowTitle("Simple paint – B=brush  E=eraser  C=clear  S=save")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(NewGame(1020, 668)); err != nil {
		panic(err)
	}
}
