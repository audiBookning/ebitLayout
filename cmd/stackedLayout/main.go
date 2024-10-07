package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a solid color
	screen.Fill(color.RGBA{0x80, 0x80, 0x80, 0xff})

	// Draw stacked rectangles (as an example)
	drawStackedRectangles(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func drawStackedRectangles(screen *ebiten.Image) {
	// Positions and sizes for the stacked rectangles
	positions := []struct {
		X, Y, Width, Height int
		Color               color.Color
	}{
		{X: 50, Y: 50, Width: 100, Height: 50, Color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{X: 50, Y: 110, Width: 100, Height: 50, Color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{X: 50, Y: 170, Width: 100, Height: 50, Color: color.RGBA{0x00, 0x00, 0xff, 0xff}},
	}

	for _, rect := range positions {
		// Create a rectangle image
		img := ebiten.NewImage(rect.Width, rect.Height)
		img.Fill(rect.Color)

		// Draw the rectangle at the specified position
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(rect.X), float64(rect.Y))
		screen.DrawImage(img, op)
	}
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Stacked Layout Example")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
