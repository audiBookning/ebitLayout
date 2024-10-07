package main

import (
	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	container := layout.NewContainer(10, 300, 200, 50, 50)

	container.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Container Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
