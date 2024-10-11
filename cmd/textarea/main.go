package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/internals/widgets"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	textarea *widgets.TextArea
}

func (g *Game) Update() error {
	return g.textarea.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.textarea.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Enhanced Text Input Example")
	game := &Game{
		textarea: widgets.NewTextArea(50, 50, 540, 300, 10),
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
