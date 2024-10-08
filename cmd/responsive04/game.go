package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/cmd/responsive05/pages"
)

type Game struct {
	currentPage Page
	pages       map[string]Page
}

// NewGame initializes the game with multiple pages and handles page switching.
func NewGame() *Game {
	g := &Game{
		pages: make(map[string]Page),
	}

	// Define a function to switch pages
	switchPage := func(pageName string) {
		if page, exists := g.pages[pageName]; exists {
			log.Printf("Switching to page: %s\n", pageName)
			g.currentPage = page
		} else {
			log.Printf("Page %s does not exist!\n", pageName)
		}
	}

	// Initialize pages with the switchPage function
	mainMenu := pages.NewMainMenuPage(switchPage)
	settings := pages.NewSettingsPage(switchPage)

	g.pages["main"] = mainMenu
	g.pages["settings"] = settings

	g.currentPage = mainMenu // Start with the main menu

	return g
}

// Update updates the current page.
func (g *Game) Update() error {
	if err := g.currentPage.Update(); err != nil {
		return err
	}
	return nil
}

// Draw renders the current page.
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a background color
	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF}) // Dark gray background

	// Draw the current page
	g.currentPage.Draw(screen)
}

// Layout handles the layout of the game window.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
