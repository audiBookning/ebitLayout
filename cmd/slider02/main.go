package main

import (
	"image/color"

	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	slider *widgets.Slider
}

func (g *Game) Update() error {
	g.slider.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	g.slider.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Slider Example")

	slider := &widgets.Slider{
		X:         100,
		Y:         200,
		Width:     400,
		Height:    20,
		HandlePos: 200,
	}

	game := &Game{slider: slider}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
