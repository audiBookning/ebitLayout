package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/internals/widgets"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	toggleButton *widgets.ToggleButton03
	inputManager *widgets.InputManager
}

func NewGame() *Game {
	inputManager := &widgets.InputManager{}
	toggleButton := widgets.NewToggleButton03(
		screenWidth/2-50, screenHeight/2-25, 100, 50,
		"Toggle",
		color.RGBA{200, 0, 0, 255},
		color.RGBA{0, 200, 0, 255},
		func() {
			log.Println("Button toggled")
		},
	)
	inputManager.Register(toggleButton)

	return &Game{
		toggleButton: toggleButton,
		inputManager: inputManager,
	}
}

func (g *Game) Update() error {
	g.inputManager.Update()
	g.toggleButton.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.toggleButton.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Toggle Button Example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
