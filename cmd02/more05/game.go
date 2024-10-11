package main

import (
	"errors"
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/builder"
	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	navigator   *navigator.Navigator
	prevWidth   int
	prevHeight  int
	exit        bool
	textWrapper *textwrapper.TextWrapper
}

func NewGame() *Game {
	fontPath := GetFilePath("assets/fonts/roboto_regularTTF.ttf")
	fontSize := 44.0
	textWrapper, err := textwrapper.NewTextWrapper(fontPath, fontSize, false)
	if err != nil {
		log.Fatalf("Failed to create TextWrapper: %v", err)
	}
	textWrapper.Color = color.RGBA{255, 255, 255, 255}

	screenWidth, screenHeight := 800, 600
	g := &Game{
		prevWidth:   screenWidth,
		prevHeight:  screenHeight,
		textWrapper: textWrapper,
	}

	onExit := func() {
		log.Println("Exiting game...")
		g.exit = true
	}

	g.navigator = navigator.NewNavigator(onExit)

	mainMenu := builder.NewMainMenuPage(g.navigator, textWrapper, screenWidth, screenHeight)
	settings := builder.NewSettingsPage(g.navigator, textWrapper, screenWidth, screenHeight)
	audio := builder.NewAudioPage(g.navigator, textWrapper, screenWidth, screenHeight)
	graphics := builder.NewGraphicsPage(g.navigator, textWrapper, screenWidth, screenHeight)
	startGame := builder.NewLevelGamePage(g.navigator, textWrapper, screenWidth, screenHeight, "start", "Start Game")

	g.navigator.AddPage("main", mainMenu)
	g.navigator.AddPage("settings", settings)
	g.navigator.AddPage("start", startGame)
	g.navigator.AddPage("audio", audio)
	g.navigator.AddPage("graphics", graphics)

	g.navigator.Layout(g.prevWidth, g.prevHeight)

	g.navigator.SwitchTo("main")

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

	//screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF})

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
