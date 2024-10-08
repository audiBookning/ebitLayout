package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Responsive Layout with Ebitengine")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
