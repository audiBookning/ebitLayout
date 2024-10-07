package main

import (
	"log"

	"example.com/menu/internals/charts"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game handles the game state.
type Game struct {
	graph *charts.Chart01
}

// NewGame initializes and returns a new Game instance.
func NewGame(screenWidth, screenHeight float64) *Game {
	graph := charts.NewChart01(screenWidth, screenHeight)
	return &Game{graph: graph}
}

// Update is called every frame (before the Draw function).
func (g *Game) Update() error {
	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Delegate drawing to the Graph struct
	g.graph.DrawAxes(screen)
	g.graph.PlotSineWave(screen)
}

// Layout sets the size of the screen.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(g.graph.ScreenWidth), int(g.graph.ScreenHeight) // Window size
}

func main() {
	screenWidth, screenHeight := 400.0, 400.0
	game := NewGame(screenWidth, screenHeight)

	// Set the window size and title
	ebiten.SetWindowSize(int(screenWidth), int(screenHeight))
	ebiten.SetWindowTitle("2D Graph with Axis Labels")

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
