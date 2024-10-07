package main

import (
	"log"

	"example.com/menu/internals/charts"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game struct manages the game state.
type Game struct {
	graph *charts.Chart02
}

// Update is called every frame (before the Draw function).
func (g *Game) Update() error {
	return nil
}

// Draw draws the current game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Initialize Graph if not already done
	if g.graph == nil {
		g.graph = charts.NewChart02(screen)
	}

	// Render the graph
	g.graph.Render()
}

// Layout sets the size of the screen.
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
