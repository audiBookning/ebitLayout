package layout

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Container struct {
	PaddingX float64
	Width    float64
	Height   float64
	X        float64
	Y        float64
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

	rectangleWidth := c.Width - 2*c.PaddingX
	rectangleHeight := c.Height

	color := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	ebitenutil.DrawRect(screen, c.X, c.Y, rectangleWidth, rectangleHeight, color)
}
