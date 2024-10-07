package main

import (
	"image/color"
	"log"

	"example.com/menu/internals/widgets"

	"github.com/hajimehoshi/ebiten/v2"
)

// ***** GAME *****
type Game struct {
	InputManager *widgets.InputManager
	centerButton *widgets.Button01
}

func NewGame() *Game {
	game := &Game{
		InputManager: &widgets.InputManager{},
	}

	// Create a button with normal, hover, and click colors
	centerButton := widgets.NewButton(200, 200, 80, 25, "CENTER",
		color.RGBA{200, 0, 0, 255},                      // Normal color
		color.RGBA{150, 0, 0, 255},                      // Hover color
		color.RGBA{100, 0, 0, 255},                      // Click color
		func() { log.Println("Center button clicked") }) // OnClick handler
	game.InputManager.Register(centerButton)
	game.centerButton = centerButton

	return game
}

func (g *Game) Update() error {
	g.InputManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.centerButton.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// ***** MAIN *****
func main() {
	game := NewGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Button with Hover and Click Effects")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
