// A minimal single-file paint program with:
//   - resizable window
//   - toolbar: Brush (black)  Eraser (white)  Clear (C)
//   - left-drag to draw/erase when tool active
package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

/* ---------- constants ---------- */

const (
	iconSize      = 32
	iconMargin    = 8
	toolbarHeight = iconSize + 2*iconMargin
)

/* ---------- game state ---------- */

type Game struct {
	width, height int

	mode            string // "brush" | "eraser"
	isDrawing       bool
	leftPressedLast bool
	prevX, prevY    int

	background *ebiten.Image // solid white underlay
	drawing    *ebiten.Image // user strokes

	// toolbar rectangles (screen space)
	brushIcon  image.Rectangle
	eraserIcon image.Rectangle
	clearIcon  image.Rectangle
}

/* ---------- constructor ---------- */

func NewGame(w, h int) *Game {
	bg := ebiten.NewImage(w, h)
	bg.Fill(color.White)

	drw := ebiten.NewImage(w, h)
	drw.Fill(color.White)

	g := &Game{
		width:      w,
		height:     h,
		background: bg,
		drawing:    drw,
	}

	// icon positions
	g.brushIcon = image.Rect(iconMargin, iconMargin,
		iconMargin+iconSize, iconMargin+iconSize)
	g.eraserIcon = image.Rect(2*iconMargin+iconSize, iconMargin,
		2*iconMargin+2*iconSize, iconMargin+iconSize)
	g.clearIcon = image.Rect(3*iconMargin+2*iconSize, iconMargin,
		3*iconMargin+3*iconSize, iconMargin+iconSize)

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
	g.drawToolbar(screen)
}

/* ---------- layout & resize ---------- */

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
	if y < toolbarHeight {
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
		}
	}

	// start stroke on canvas when a tool is active
	if g.mode == "brush" || g.mode == "eraser" {
		g.isDrawing = true
		g.prevX, g.prevY = x, y
	}
}

func (g *Game) continueStroke() {
	x, y := ebiten.CursorPosition()
	if x == g.prevX && y == g.prevY {
		return
	}

	col := color.Black
	if g.mode == "eraser" {
		col = color.White
	}

	ebitenutil.DrawLine(g.drawing,
		float64(g.prevX), float64(g.prevY),
		float64(x), float64(y), col)
	g.prevX, g.prevY = x, y
}

func (g *Game) drawToolbar(screen *ebiten.Image) {
	// background
	tb := ebiten.NewImage(g.width, toolbarHeight)
	tb.Fill(color.RGBA{220, 220, 220, 255})
	screen.DrawImage(tb, nil)

	// helper to draw one icon
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
}

/* ---------- entry point ---------- */

func main() {
	ebiten.SetWindowSize(1020, 668)
	ebiten.SetWindowTitle("Simple paint – B=brush  E=eraser  C=clear")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(NewGame(1020, 668)); err != nil {
		panic(err)
	}
}
