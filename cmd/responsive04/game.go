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

func NewGame() *Game {
	g := &Game{
		pages: make(map[string]Page),
	}

	switchPage := func(pageName string) {
		if page, exists := g.pages[pageName]; exists {
			log.Printf("Switching to page: %s\n", pageName)
			g.currentPage = page
		} else {
			log.Printf("Page %s does not exist!\n", pageName)
		}
	}

	mainMenu := pages.NewMainMenuPage(switchPage)
	settings := pages.NewSettingsPage(switchPage)

	g.pages["main"] = mainMenu
	g.pages["settings"] = settings

	g.currentPage = mainMenu

	return g
}

func (g *Game) Update() error {
	if err := g.currentPage.Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF})

	g.currentPage.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
