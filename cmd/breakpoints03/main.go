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

type UIElement interface {
	Draw(screen *ebiten.Image, layout Layout)
	GetBreakpoint() Breakpoint
}

type RedRectangle struct {
	Width      int
	Height     int
	Breakpoint Breakpoint
}

func (r *RedRectangle) Draw(screen *ebiten.Image, layout Layout) {
	if layout.Breakpoint != r.Breakpoint {
		return // Skip drawing if the layout's breakpoint doesn't match the element's breakpoint
	}
	uiElement := ebiten.NewImage(r.Width, r.Height)
	uiElement.Fill(color.RGBA{255, 0, 0, 255}) // Red rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-r.Width/2), float64(layout.Height/2-r.Height/2))
	screen.DrawImage(uiElement, op)
}

func (r *RedRectangle) GetBreakpoint() Breakpoint {
	return r.Breakpoint
}

type BlueCircle struct {
	Radius     int
	Breakpoint Breakpoint
}

func (b *BlueCircle) Draw(screen *ebiten.Image, layout Layout) {
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

func (b *BlueCircle) GetBreakpoint() Breakpoint {
	return b.Breakpoint
}

type Game struct {
	layoutSystem *LayoutSystem
	currentBP    Breakpoint
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
	game := &Game{
		layoutSystem: layoutSystem,
		elements: []UIElement{
			&RedRectangle{Width: 100, Height: 50, Breakpoint: Medium},
			&BlueCircle{Radius: 40, Breakpoint: Large},
		},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
