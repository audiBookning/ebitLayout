package main

import (
	"errors"
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/customPages"
	"example.com/menu/cmd02/more03/navigator"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	navigator *navigator.Navigator
	exit      bool
}

// NewGame initializes the game with multiple pages and handles page switching.
func NewGame() *Game {
	g := &Game{}

	// Define the onExit callback
	onExit := func() {
		log.Println("Exiting game...")
		g.exit = true
	}

	// Initialize navigator *without* setting the initial page yet
	g.navigator = navigator.NewNavigator(onExit)

	// Initialize pages with the navigator's SwitchTo method
	mainMenu := customPages.NewMainMenuPage(g.navigator)
	settings := customPages.NewSettingsPage(g.navigator)
	startGame := customPages.NewLevelGamePage(g.navigator)
	audio := customPages.NewAudioPage(g.navigator)
	graphics := customPages.NewGraphicsPage(g.navigator)
	level01 := customPages.NewLevel01Page(g.navigator)
	level02 := customPages.NewLevel02Page(g.navigator)

	// Add pages to navigator
	g.navigator.AddPage("main", mainMenu)
	g.navigator.AddPage("settings", settings)
	g.navigator.AddPage("start", startGame)
	g.navigator.AddPage("audio", audio)
	g.navigator.AddPage("graphics", graphics)
	g.navigator.AddPage("level01", level01)
	g.navigator.AddPage("level02", level02)

	// Set the initial page
	//g.navigator.current = mainMenu
	g.navigator.SwitchTo("main") // Start with the main menu

	return g
}

// Update updates the current page.
func (g *Game) Update() error {
	if g.exit {
		return errors.New("game exited by user")
	}

	if err := g.navigator.CurrentPage().Update(); err != nil {
		return err
	}
	return nil
}

// Draw renders the current page.
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a background color
	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF}) // Dark gray background

	// Draw the current page
	g.navigator.CurrentPage().Draw(screen)
}

// Layout handles the layout of the game window.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}