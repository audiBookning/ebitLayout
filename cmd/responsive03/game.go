package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/cmd/responsive03/responsive"
)

type Game struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
}

func NewGame() *Game {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Button1", func() { log.Println("Button1 clicked") }),
		responsive.NewButton("Button2", func() { log.Println("Button2 clicked") }),
		responsive.NewButton("Button3", func() { log.Println("Button3 clicked") }),
		responsive.NewButton("Button4", func() { log.Println("Button4 clicked") }),
	}

	ui := responsive.NewUI("Responsive UI", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &Game{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

func (g *Game) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != g.prevWidth || screenHeight != g.prevHeight {
		log.Printf("Window resized to %dx%d\n", screenWidth, screenHeight)
		g.prevWidth = screenWidth
		g.prevHeight = screenHeight
	}

	g.ui.Update(screenWidth, screenHeight)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.ui.HandleClick(x, y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF})

	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
