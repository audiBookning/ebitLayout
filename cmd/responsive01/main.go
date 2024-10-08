package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

// Update updates the game state.
func (g *Game) Update() error {
	g.UI.Update(g.ScreenWidth, g.ScreenHeight)

	// Handle mouse input for button click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		if g.UI.IsButtonClicked(mouseX, mouseY) {
			g.UI.OnButtonClick()
		}
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff}) // Background color

	g.UI.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight
	g.UI.Update(g.ScreenWidth, g.ScreenHeight) // Ensure UI updates on layout change
	return outsideWidth, outsideHeight
}

type Game struct {
	UI           *UI
	ScreenWidth  int
	ScreenHeight int
}

func NewGame() *Game {
	return &Game{
		UI: NewUI(),
	}
}
