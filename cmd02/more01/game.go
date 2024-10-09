package main

import (
	"errors"
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/pages"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	currentPage Page
	pages       map[string]Page
	exit        bool
}

func NewGame() *Game {
	g := &Game{
		pages: make(map[string]Page),
	}

	// Define a function to switch pages
	switchPage := func(pageName string) {
		if pageName == "exit" {
			log.Println("Exit requested")
			g.exit = true
			return
		}
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
	startGame := pages.NewStartGamePage(switchPage)
	audio := pages.NewAudioPage(switchPage)
	graphics := pages.NewGraphicsPage(switchPage)

	g.pages["main"] = mainMenu
	g.pages["settings"] = settings
	g.pages["start"] = startGame
	g.pages["audio"] = audio
	g.pages["graphics"] = graphics

	g.currentPage = mainMenu // Start with the main menu

	return g
}

func (g *Game) Update() error {
	if g.exit {
		return errors.New("game exited by user")
	}

	if err := g.currentPage.Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a background color
	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF}) // Dark gray background

	// Draw the current page
	g.currentPage.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
