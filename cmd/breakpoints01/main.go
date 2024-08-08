package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Breakpoint int

const (
	ExtraSmall Breakpoint = iota
	Small
	Medium
	Large
	ExtraLarge
	ExtraExtraLarge
)

type Layout struct {
	Breakpoint Breakpoint
	Width      int
	Height     int
}

type LayoutSystem struct {
	Layouts map[Breakpoint]Layout
}

func NewLayoutSystem() *LayoutSystem {
	return &LayoutSystem{
		Layouts: map[Breakpoint]Layout{
			ExtraSmall:      {Breakpoint: ExtraSmall, Width: 320, Height: 480},
			Small:           {Breakpoint: Small, Width: 576, Height: 768},
			Medium:          {Breakpoint: Medium, Width: 768, Height: 1024},
			Large:           {Breakpoint: Large, Width: 992, Height: 1280},
			ExtraLarge:      {Breakpoint: ExtraLarge, Width: 1200, Height: 1600},
			ExtraExtraLarge: {Breakpoint: ExtraExtraLarge, Width: 1400, Height: 1920},
		},
	}
}

func (ls *LayoutSystem) GetLayout(breakpoint Breakpoint) Layout {
	return ls.Layouts[breakpoint]
}

func drawLayout(screen *ebiten.Image, layout Layout) {
	// Fill background with a color
	screen.Fill(color.RGBA{0, 0, 255, 255}) // Blue background

	// Draw a red rectangle to represent a UI element
	uiElement := ebiten.NewImage(100, 50)
	uiElement.Fill(color.RGBA{255, 0, 0, 255}) // Red rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-50), float64(layout.Height/2-25))
	screen.DrawImage(uiElement, op)
}

type Game struct {
	layoutSystem *LayoutSystem
	currentBP    Breakpoint
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	layout := g.layoutSystem.GetLayout(g.currentBP)
	drawLayout(screen, layout)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Determine breakpoint based on window size
	switch {
	case outsideWidth < 576:
		g.currentBP = ExtraSmall
	case outsideWidth < 768:
		g.currentBP = Small
	case outsideWidth < 992:
		g.currentBP = Medium
	case outsideWidth < 1200:
		g.currentBP = Large
	case outsideWidth < 1400:
		g.currentBP = ExtraLarge
	default:
		g.currentBP = ExtraExtraLarge
	}
	layout := g.layoutSystem.GetLayout(g.currentBP)
	return layout.Width, layout.Height
}

func main() {
	layoutSystem := NewLayoutSystem()
	game := &Game{layoutSystem: layoutSystem}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
