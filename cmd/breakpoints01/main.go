package main

import (
	"log"

	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	layoutSystem *layout.BreakpointLayoutSystem
	currentBP    layout.Breakpoint
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	layout := g.layoutSystem.GetLayout(g.currentBP)
	layout.DrawLayout(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.currentBP = g.layoutSystem.DetermineBreakpoint(outsideWidth, outsideHeight)
	currentLayout := g.layoutSystem.GetLayout(g.currentBP)
	return currentLayout.Width, currentLayout.Height
}

func main() {
	layoutSystem := layout.NewLayoutSystem(nil)
	game := &Game{layoutSystem: layoutSystem}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
