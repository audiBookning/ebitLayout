package main

import (
	"image/color"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/internals/layout"
)

type Game struct {
	flexbox *layout.FlexBox
}

func NewGame() *Game {
	flexbox := &layout.FlexBox{
		Elements: []layout.Element{
			{Width: 100, Height: 100, Color: color.RGBA{255, 0, 0, 255}, Flex: 0},
			{Width: 150, Height: 100, Color: color.RGBA{0, 255, 0, 255}, Flex: 2},
			{Width: 100, Height: 100, Color: color.RGBA{0, 0, 255, 255}, Flex: 0},
		},
		Direction:      "row",
		JustifyContent: "space-around",
		AlignItems:     "center",
	}
	return &Game{flexbox: flexbox}
}

func (g *Game) Update() error {

	windowWidth, windowHeight := ebiten.WindowSize()

	g.flexbox.Layout(windowWidth, windowHeight)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.flexbox.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Flex Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
