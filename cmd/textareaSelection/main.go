package main

import (
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	textAreaX    = 50
	textAreaY    = 50
	textAreaW    = 540
	textAreaH    = 300
)

type Game struct {
	textarea *widgets.TextAreaSelection
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.textarea.Draw(screen)
}

func (g *Game) Update() error {
	return g.textarea.Update()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text Input with Selection Example")
	game := &Game{
		textarea: widgets.NewTextAreaSelection(textAreaX, textAreaY, textAreaW, textAreaH, 10),
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
