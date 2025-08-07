// Modern single-file paint with fixed top bar for palette and toolbar
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
	iconSize    = 32
	iconMargin  = 8
	paletteCols = 8
	paletteRows = 2
	swatchSize  = 24
	palettePad  = 4

	minWidth = 1
	maxWidth = 5
	sliderH  = 4
	sliderW  = 120
	handleW  = 12

	topBarHeight = iconSize + 2*iconMargin + paletteRows*(swatchSize+palettePad) + palettePad
)

/* ---------- game state ---------- */
type Game struct {
	width, height int

	mode            string // "brush" | "eraser"
	currentColor    color.Color
	brushSize       int
	isDrawing       bool
	leftPressedLast bool
	prevX, prevY    int

	background *ebiten.Image
	drawing    *ebiten.Image

	// rectangles
	brushIcon   image.Rectangle
	eraserIcon  image.Rectangle
	clearIcon   image.Rectangle
	saveIcon    image.Rectangle
	sliderTrack image.Rectangle
	sliderKnob  image.Rectangle
	palette     []color.Color
}

/* ---------- constructor ---------- */
func NewGame(w, h int) *Game {
	bg := ebiten.NewImage(w, h)
	bg.Fill(color.White)
	drw := ebiten.NewImage(w, h)
	drw.Fill(color.White)

	g := &Game{
		width:        w,
		height:       h,
		background:   bg,
		drawing:      drw,
		currentColor: color.Black,
		brushSize:    3,
	}

	g.palette = []color.Color{
		color.Black, color.RGBA{127, 0, 0, 255}, color.RGBA{0, 127, 0, 255},
		color.RGBA{0, 0, 127, 255}, color.RGBA{127, 127, 0, 255},
		color.RGBA{127, 0, 127, 255}, color.RGBA{0, 127, 127, 255},
		color.RGBA{127, 127, 127, 255}, color.White,
		color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 0, 255, 255}, color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 255, 255}, color.RGBA{0, 255, 255, 255},
		color.RGBA{255, 255, 255, 255},
	}

	g.brushIcon = image.Rect(iconMargin, iconMargin, iconMargin+iconSize, iconMargin+iconSize)
	g.eraserIcon = image.Rect(2*iconMargin+iconSize, iconMargin, 2*iconMargin+2*iconSize, iconMargin+iconSize)
	g.clearIcon = image.Rect(3*iconMargin+2*iconSize, iconMargin, 3*iconMargin+3*iconSize, iconMargin+iconSize)
	g.saveIcon = image.Rect(4*iconMargin+3*iconSize, iconMargin, 4*iconMargin+4*iconSize, iconMargin+iconSize)

	g.sliderTrack = image.Rect(iconMargin, topBarHeight-10, iconMargin+sliderW, topBarHeight-10+sliderH)
	g.updateSliderKnob()
	return g
}

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

	topBar := ebiten.NewImage(g.width, topBarHeight)
	topBar.Fill(color.RGBA{240, 240, 240, 255})
	screen.DrawImage(topBar, nil)

	g.drawUI(screen)
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

func (g *Game) updateSliderKnob() {
	ratio := float64(g.brushSize-minWidth) / float64(maxWidth-minWidth)
	x := g.sliderTrack.Min.X + int(ratio*float64(sliderW-handleW))
	g.sliderKnob = image.Rect(x, g.sliderTrack.Min.Y-4, x+handleW, g.sliderTrack.Min.Y+sliderH+4)
}

func (g *Game) handleClick() {
	x, y := ebiten.CursorPosition()
	if y < topBarHeight {
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
		for i, col := range g.palette {
			px := palettePad + (i%paletteCols)*(swatchSize+palettePad)
			py := iconSize + 2*iconMargin + palettePad + (i/paletteCols)*(swatchSize+palettePad)
			if x >= px && x < px+swatchSize && y >= py && y < py+swatchSize {
				g.currentColor = col
				return
			}
		}
		if image.Pt(x, y).In(g.sliderTrack) || image.Pt(x, y).In(g.sliderKnob) {
			ratio := float64(x-g.sliderTrack.Min.X) / float64(sliderW)
			val := int(ratio*(maxWidth-minWidth)+0.5) + minWidth
			if val < minWidth {
				val = minWidth
			}
			if val > maxWidth {
				val = maxWidth
			}
			g.brushSize = val
			g.updateSliderKnob()
			return
		}
	} else if g.mode == "brush" || g.mode == "eraser" {
		g.isDrawing = true
		g.prevX, g.prevY = x, y
	}
}

func (g *Game) continueStroke() {
	x, y := ebiten.CursorPosition()
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
	ebitenutil.DrawRect(g.drawing,
		float64(x)-float64(g.brushSize)/2,
		float64(y)-float64(g.brushSize)/2,
		float64(g.brushSize), float64(g.brushSize), col)
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

func (g *Game) drawUI(screen *ebiten.Image) {
	// draw icons
	drawRect := func(rect image.Rectangle, label string, active bool, activeCol color.Color) {
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
	drawRect(g.brushIcon, "B", g.mode == "brush", color.RGBA{150, 150, 255, 255})
	drawRect(g.eraserIcon, "E", g.mode == "eraser", color.RGBA{255, 150, 150, 255})
	drawRect(g.clearIcon, "C", false, color.White)
	drawRect(g.saveIcon, "S", false, color.White)

	// palette
	for i, col := range g.palette {
		x := palettePad + (i%paletteCols)*(swatchSize+palettePad)
		y := iconSize + 2*iconMargin + palettePad + (i/paletteCols)*(swatchSize+palettePad)
		sq := ebiten.NewImage(swatchSize, swatchSize)
		sq.Fill(col)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(sq, op)
	}

	// slider
	track := ebiten.NewImage(sliderW, sliderH)
	track.Fill(color.Gray{160})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.sliderTrack.Min.X), float64(g.sliderTrack.Min.Y))
	screen.DrawImage(track, op)
	knob := ebiten.NewImage(handleW, handleW)
	knob.Fill(color.Gray{80})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.sliderKnob.Min.X), float64(g.sliderKnob.Min.Y))
	screen.DrawImage(knob, op)
}

/* ---------- entry point ---------- */
func main() {
	ebiten.SetWindowSize(1020, 668)
	ebiten.SetWindowTitle("Modern Paint – B/E/C/S + palette + slider")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(NewGame(1020, 668)); err != nil {
		panic(err)
	}
}
