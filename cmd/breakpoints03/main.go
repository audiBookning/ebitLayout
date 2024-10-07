package main

import (
	"image/color"
	"log"

	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
)

type UIElement interface {
	Draw(screen *ebiten.Image, layout layout.BreakpointLayout)
	GetBreakpoint() layout.Breakpoint
}

type RedRectangle struct {
	Width      int
	Height     int
	Breakpoint layout.Breakpoint
}

func (r *RedRectangle) Draw(screen *ebiten.Image, layout layout.BreakpointLayout) {
	if layout.Breakpoint != r.Breakpoint {
		return // Skip drawing if the layout's breakpoint doesn't match the element's breakpoint
	}
	uiElement := ebiten.NewImage(r.Width, r.Height)
	uiElement.Fill(color.RGBA{255, 0, 0, 255}) // Red rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-r.Width/2), float64(layout.Height/2-r.Height/2))
	screen.DrawImage(uiElement, op)
}

func (r *RedRectangle) GetBreakpoint() layout.Breakpoint {
	return r.Breakpoint
}

type BlueCircle struct {
	Radius     int
	Breakpoint layout.Breakpoint
}

func (b *BlueCircle) Draw(screen *ebiten.Image, layout layout.BreakpointLayout) {
	if layout.Breakpoint != b.Breakpoint {
		return // Skip drawing if the layout's breakpoint doesn't match the element's breakpoint
	}
	diameter := b.Radius * 2
	circle := ebiten.NewImage(diameter, diameter)
	circle.Fill(color.Transparent) // Set transparent background

	for x := 0; x < diameter; x++ {
		for y := 0; y < diameter; y++ {
			dx := x - b.Radius
			dy := y - b.Radius
			if dx*dx+dy*dy <= b.Radius*b.Radius {
				circle.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue color
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-b.Radius), float64(layout.Height/2-b.Radius))
	screen.DrawImage(circle, op)
}

func (b *BlueCircle) GetBreakpoint() layout.Breakpoint {
	return b.Breakpoint
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
	// Draw all UI elements that match the current breakpoint
	for _, element := range g.elements {
		if element.GetBreakpoint() == g.currentBP {
			element.Draw(screen, layout)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Determine breakpoint based on window size
	switch {
	case outsideWidth < 576:
		g.currentBP = layout.Bp_ExtraSmall
	case outsideWidth < 768:
		g.currentBP = layout.Bp_Small
	case outsideWidth < 992:
		g.currentBP = layout.Bp_Medium
	case outsideWidth < 1200:
		g.currentBP = layout.Bp_Large
	case outsideWidth < 1400:
		g.currentBP = layout.Bp_ExtraLarge
	default:
		g.currentBP = layout.Bp_ExtraExtraLarge
	}
	layout := g.layoutSystem.GetLayout(g.currentBP)
	return layout.Width, layout.Height
}

func main() {
	layoutSystem := layout.NewLayoutSystem()
	game := &Game{
		layoutSystem: layoutSystem,
		elements: []UIElement{
			&RedRectangle{Width: 100, Height: 50, Breakpoint: layout.Bp_Medium},
			&BlueCircle{Radius: 40, Breakpoint: layout.Bp_Large},
		},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
