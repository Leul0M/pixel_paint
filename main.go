package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	iconSize      = 32
	iconMargin    = 8
	toolbarHeight = iconSize + 2*iconMargin
)

type Game struct {
	width, height int
	mode          string // "brush" or "idle"
	isDrawing     bool
	leftPressedLast bool
	prevX, prevY  int
	background    *ebiten.Image
	drawing       *ebiten.Image

	brushIcon image.Rectangle
	clearIcon image.Rectangle
}

func NewGame(w, h int) *Game {
	bg := ebiten.NewImage(w, h)
	bg.Fill(color.White)

	drw := ebiten.NewImage(w, h)
	drw.Fill(color.White)

	g := &Game{
		width:      w,
		height:     h,
		mode:       "idle",
		background: bg,
		drawing:    drw,
	}

	g.brushIcon = image.Rect(iconMargin, iconMargin,
		iconMargin+iconSize, iconMargin+iconSize)
	g.clearIcon = image.Rect(2*iconMargin+iconSize, iconMargin,
		2*iconMargin+2*iconSize, iconMargin+iconSize)

	return g
}

func (g *Game) Update() error {
	// --- toolbar / icon clicks ---
	leftPressedNow := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
justPressed := leftPressedNow && !g.leftPressedLast
g.leftPressedLast = leftPressedNow

	if justPressed {
		if g.mode == "idle" {
			g.mode = "brush"
			return nil
}
		x, y := ebiten.CursorPosition()

		// toolbar strip
		if y < toolbarHeight {
			if image.Pt(x, y).In(g.brushIcon) {
				g.mode = "brush"
				return nil
			}
			if image.Pt(x, y).In(g.clearIcon) {
				g.drawing.Fill(color.White)
				return nil
			}
		}

		// start stroke
		if g.mode == "brush" {
			g.isDrawing = true
			g.prevX, g.prevY = x, y
		}
	}

	// --- stroke continuation ---
	if g.mode == "brush" && g.isDrawing &&
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.drawStroke()
	} else {
		g.isDrawing = false
	}
	return nil
}

func (g *Game) drawStroke() {
	x, y := ebiten.CursorPosition()
	if x == g.prevX && y == g.prevY {
		return
	}
	ebitenutil.DrawLine(g.drawing,
		float64(g.prevX), float64(g.prevY),
		float64(x), float64(y), color.Black)
	g.prevX, g.prevY = x, y
}

func (g *Game) Draw(screen *ebiten.Image) {
	// canvas
	screen.DrawImage(g.background, nil)
	screen.DrawImage(g.drawing, nil)

	// --- toolbar background ---
	tb := ebiten.NewImage(g.width, toolbarHeight)
	tb.Fill(color.RGBA{220, 220, 220, 255})
	screen.DrawImage(tb, &ebiten.DrawImageOptions{})

	// --- brush icon ---
	brushImg := ebiten.NewImage(iconSize, iconSize)
	if g.mode == "brush" {
		brushImg.Fill(color.RGBA{150, 150, 255, 255}) // highlight
	} else {
		brushImg.Fill(color.White)
	}
	ebitenutil.DrawRect(brushImg, 2, 2, iconSize-4, iconSize-4, color.Black)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.brushIcon.Min.X), float64(g.brushIcon.Min.Y))
	screen.DrawImage(brushImg, op)

	// --- clear icon ---
	clearImg := ebiten.NewImage(iconSize, iconSize)
	clearImg.Fill(color.White)
	ebitenutil.DebugPrintAt(clearImg, "C", iconSize/2-4, iconSize/2-6)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.clearIcon.Min.X), float64(g.clearIcon.Min.Y))
	screen.DrawImage(clearImg, op)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return g.width, g.height
}

func main() {
	ebiten.SetWindowSize(1020, 668)
	ebiten.SetWindowTitle("Simple painting app")
	if err := ebiten.RunGame(NewGame(640, 480)); err != nil {
		panic(err)
	}
}