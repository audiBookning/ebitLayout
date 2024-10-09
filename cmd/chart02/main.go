package main

import (
	"log"

	"example.com/menu/internals/charts"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	graph *charts.Chart02
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Initialize Graph if not already done
	if g.graph == nil {
		g.graph = charts.NewChart02(screen)
	}

	// Render the graph
	g.graph.Render()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 400 // Window size
}

func main() {
	game := &Game{}

	// Set the window size and title
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("2D Graph with Numbered Ticks")

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
