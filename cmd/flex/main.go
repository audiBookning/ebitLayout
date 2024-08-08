package main

import (
	"image/color"
	"log"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Element struct {
	X, Y, Width, Height int
	Color               color.Color
	Flex                int // This can represent the flex-grow/shrink basis
}

type FlexBox struct {
	Elements       []Element
	Direction      string // "row" or "column"
	JustifyContent string // "flex-start", "center", "flex-end", "space-between", "space-around"
	AlignItems     string // "flex-start", "center", "flex-end"
}

func (fb *FlexBox) Layout(width, height int) {
	totalFlex := 0
	totalSize := 0

	for _, element := range fb.Elements {
		totalFlex += element.Flex
		if fb.Direction == "row" {
			totalSize += element.Width
		} else {
			totalSize += element.Height
		}
	}

	var offset int
	if fb.Direction == "row" {
		// Calculate available space
		availableSpace := width - totalSize
		for i := range fb.Elements {
			if fb.Elements[i].Flex > 0 {
				// Distribute available space according to flex-grow
				flexSize := (availableSpace * fb.Elements[i].Flex) / totalFlex
				fb.Elements[i].Width += flexSize
			}
		}

		offset = 0
		for i := range fb.Elements {
			fb.Elements[i].X = offset
			fb.Elements[i].Y = (height - fb.Elements[i].Height) / 2
			offset += fb.Elements[i].Width
		}
	} else {
		// Handle column layout similarly
		availableSpace := height - totalSize
		for i := range fb.Elements {
			if fb.Elements[i].Flex > 0 {
				flexSize := (availableSpace * fb.Elements[i].Flex) / totalFlex
				fb.Elements[i].Height += flexSize
			}
		}

		offset = 0
		for i := range fb.Elements {
			fb.Elements[i].X = (width - fb.Elements[i].Width) / 2
			fb.Elements[i].Y = offset
			offset += fb.Elements[i].Height
		}
	}
}

func (fb *FlexBox) Draw(screen *ebiten.Image) {
	for _, element := range fb.Elements {
		vector.DrawFilledRect(screen, float32(element.X), float32(element.Y), float32(element.Width), float32(element.Height), element.Color, true)
	}
}

type Game struct {
	flexbox *FlexBox
}

func NewGame() *Game {
	flexbox := &FlexBox{
		Elements: []Element{
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
