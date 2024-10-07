package layout

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Container defines a struct to represent a container with padding and centering
type Container struct {
	PaddingX float64 // Padding on the left and right
	Width    float64 // Width of the container
	Height   float64 // Height of the container
	X        float64 // X position of the container
	Y        float64 // Y position of the container
}

func NewContainer(paddingX, width, height, x, y float64) *Container {
	return &Container{
		PaddingX: paddingX,
		Width:    width,
		Height:   height,
		X:        x,
		Y:        y,
	}
}

func (c *Container) Draw(screen *ebiten.Image) {
	// Create a rectangle with padding
	rectangleWidth := c.Width - 2*c.PaddingX
	rectangleHeight := c.Height

	// Define the color for the container
	color := color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue color

	// Draw the rectangle to represent the container
	ebitenutil.DrawRect(screen, c.X, c.Y, rectangleWidth, rectangleHeight, color)
}
