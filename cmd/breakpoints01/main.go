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
	// Determine breakpoint based on window size
	switch {
	case outsideWidth < 576:
		g.currentBP = layout.Bp_ExtraSmall
	case outsideWidth < 768:
		g.currentBP = layout.Bp_Small
	case outsideWidth < 992:
		g.currentBP = layout.Bp_Medium
	case outsideWidth < 1200:
		g.currentBP = layout.Bp_Large
	case outsideWidth < 1400:
		g.currentBP = layout.Bp_ExtraLarge
	default:
		g.currentBP = layout.Bp_ExtraExtraLarge
	}
	layout := g.layoutSystem.GetLayout(g.currentBP)
	return layout.Width, layout.Height
}

func main() {
	layoutSystem := layout.NewLayoutSystem()
	game := &Game{layoutSystem: layoutSystem}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
