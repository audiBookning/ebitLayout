package main

import (
	"image/color"
	"log"

	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	toggleButton *widgets.ToggleButton04
	inputManager *widgets.InputManager
}

func NewGame() *Game {
	inputManager := &widgets.InputManager{}

	tx, err := textwrapper.NewTextWrapper(
		"cmd/buttonToggle03/fonts/roboto_regularTTF.ttf",
		14,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	tx.Color = color.RGBA{255, 255, 255, 255} // white
	tx.SetFontSize(12)

	toggleButton := widgets.NewToggleButton04(
		screenWidth/2-50, screenHeight/2-25, 100, 50,
		"ON", "OFF",
		color.RGBA{200, 0, 0, 255},
		color.RGBA{0, 200, 0, 255},
		tx,
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
