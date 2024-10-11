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

func (g *Game) Update() error {
	g.UI.Update(g.ScreenWidth, g.ScreenHeight)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		g.UI.HandleClick(mouseX, mouseY)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x30, 0x30, 0x30, 0xFF})

	g.UI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight
	g.UI.Update(g.ScreenWidth, g.ScreenHeight)
	return outsideWidth, outsideHeight
}
