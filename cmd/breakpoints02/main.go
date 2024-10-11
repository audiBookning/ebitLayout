package main

import (
	"image/color"
	"log"

	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
)

type UIElement interface {
	Draw(screen *ebiten.Image, layout layout.BreakpointLayout)
}

type RedRectangle struct {
	Width  int
	Height int
}

func (r *RedRectangle) Draw(screen *ebiten.Image, layout layout.BreakpointLayout) {
	uiElement := ebiten.NewImage(r.Width, r.Height)
	uiElement.Fill(color.RGBA{255, 0, 0, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-r.Width/2), float64(layout.Height/2-r.Height/2))
	screen.DrawImage(uiElement, op)
}

type BlueCircle struct {
	Radius int
}

func (b *BlueCircle) Draw(screen *ebiten.Image, layout layout.BreakpointLayout) {
	diameter := b.Radius * 2
	circle := ebiten.NewImage(diameter, diameter)
	circle.Fill(color.Transparent)

	for x := 0; x < diameter; x++ {
		for y := 0; y < diameter; y++ {
			dx := x - b.Radius
			dy := y - b.Radius
			if dx*dx+dy*dy <= b.Radius*b.Radius {
				circle.Set(x, y, color.RGBA{0, 0, 255, 255})
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-b.Radius), float64(layout.Height/2-b.Radius))
	screen.DrawImage(circle, op)
}

type Game struct {
	layoutSystem *layout.BreakpointLayoutSystem
	currentBP    layout.Breakpoint
	elements     []UIElement
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	layout := g.layoutSystem.GetLayout(g.currentBP)

	for _, element := range g.elements {
		element.Draw(screen, layout)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	g.currentBP = g.layoutSystem.DetermineBreakpoint(outsideWidth, outsideHeight)
	currentLayout := g.layoutSystem.GetLayout(g.currentBP)
	return currentLayout.Width, currentLayout.Height
}

func main() {
	layoutSystem := layout.NewLayoutSystem(nil)
	game := &Game{
		layoutSystem: layoutSystem,
		elements: []UIElement{
			&RedRectangle{Width: 100, Height: 50},
			&BlueCircle{Radius: 40},
		},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
