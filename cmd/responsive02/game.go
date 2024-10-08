package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	UI           *UI
	ScreenWidth  int
	ScreenHeight int
}

func NewGame() *Game {
	// Define customizable breakpoints
	breakpoints := []Breakpoint{
		{
			Width:      1200,
			LayoutMode: "horizontal",
		},
		{
			Width:      800,
			LayoutMode: "grid",
		},
		{
			Width:      500,
			LayoutMode: "vertical",
		},
	}

	return &Game{
		UI: NewUI(breakpoints),
	}
}

// Update updates the game state.
func (g *Game) Update() error {
	g.UI.Update(g.ScreenWidth, g.ScreenHeight)

	// Handle mouse input for button clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		g.UI.HandleClick(mouseX, mouseY)
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x30, 0x30, 0x30, 0xFF}) // Dark background

	g.UI.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight
	g.UI.Update(g.ScreenWidth, g.ScreenHeight) // Ensure UI updates on layout change
	return outsideWidth, outsideHeight
}
