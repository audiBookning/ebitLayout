package main

import (
	"errors"
	"log"

	"example.com/menu/cmd02/more05/builder"
	"example.com/menu/cmd02/more05/navigator"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	navigator  *navigator.Navigator
	exit       bool
	prevWidth  int
	prevHeight int
}

func NewGame() *Game {
	screenWidth, screenHeight := 800, 600
	g := &Game{
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}

	// Define the onExit callback
	onExit := func() {
		log.Println("Exiting game...")
		g.exit = true
	}

	// Initialize navigator *without* setting the initial page yet
	g.navigator = navigator.NewNavigator(onExit)

	// Initialize pages with the navigator's SwitchTo method
	mainMenu := builder.NewMainMenuPage(g.navigator, screenWidth, screenHeight)
	settings := builder.NewSettingsPage(g.navigator, screenWidth, screenHeight)
	audio := builder.NewAudioPage(g.navigator, screenWidth, screenHeight)
	graphics := builder.NewGraphicsPage(g.navigator, screenWidth, screenHeight)
	startGame := builder.NewLevelGamePage(g.navigator, screenWidth, screenHeight, "start", "Start Game")

	// Add pages to navigator
	g.navigator.AddPage("main", mainMenu)
	g.navigator.AddPage("settings", settings)
	g.navigator.AddPage("start", startGame)
	g.navigator.AddPage("audio", audio)
	g.navigator.AddPage("graphics", graphics)

	g.navigator.Layout(g.prevWidth, g.prevHeight)

	// Set the initial page
	g.navigator.SwitchTo("main") // Start with the main menu

	return g
}

func (g *Game) Update() error {
	if g.exit {
		return errors.New("game exited by user")
	}

	if err := g.navigator.CurrentActivePage().Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a background color
	//screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF}) // Dark gray background

	// Draw the current page
	g.navigator.CurrentActivePage().Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != g.prevWidth || outsideHeight != g.prevHeight {
		g.prevWidth = outsideWidth
		g.prevHeight = outsideHeight
		g.navigator.Layout(g.prevWidth, g.prevHeight)
	}
	return outsideWidth, outsideHeight
}
