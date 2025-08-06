package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game struct to hold game state and window dimensions
type Game struct {
	width     int
	height    int
	isDrawing bool // Flag to indicate if drawing is in progress
	prevX     int  // Previous mouse X position
	prevY     int  // Previous mouse Y position
	background *ebiten.Image // Background image
	drawing    *ebiten.Image // Drawing image
	
}

// NewGame creates a new Game instance with specified dimensions
func NewGame(width, height int) *Game {
	background := ebiten.NewImage(width, height)
	background.Fill(color.White) // White background

	drawing := ebiten.NewImage(width, height)
	drawing.Fill(color.White) // White drawing image

	return &Game{
		width:    width,
		height:   height,
		background: background,
		drawing:   drawing,
	}
}

// Update is called once per frame. It is used to advance the game state.
func (g *Game) Update() error {
	g.drawBrush()
	return nil
}

// Draw is called once per frame to render the game. It draws the background and the drawing image.
func (g *Game) Draw(screen *ebiten.Image) {
    screen.DrawImage(g.background, &ebiten.DrawImageOptions{})
    screen.DrawImage(g.drawing, &ebiten.DrawImageOptions{})
}

func (g *Game) drawBrush() {
	// Get the current mouse position
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.isDrawing {
			// Start drawing, set the previous position to the current position
			g.prevX, g.prevY = x, y
			g.isDrawing = true
		}  else {
            // only draw once per mouse-move
            if x != g.prevX || y != g.prevY {
                ebitenutil.DrawLine(g.drawing,
                    float64(g.prevX), float64(g.prevY),
                    float64(x), float64(y),
                    color.Black)
                g.prevX, g.prevY = x, y
            }
        }
    } else {
        g.isDrawing = false
    }
}

// Layout takes the outside dimensions of the window and returns the dimensions
// of the game. This function is used to specify the logical size of the game
// screen.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

// main is the entry point for the executable.
//
// It creates a new Game with a resolution of 1024x768, and then passes it to
// ebiten.RunGame to start the game loop. If there is an error initializing the
// game, it will panic with that error.
func main() {
	ebiten.SetWindowSize(1020, 668)
	ebiten.SetWindowTitle("Simple painting app")
	if err := ebiten.RunGame(NewGame(640, 480)); err != nil {
		panic(err)
	}
}