package layout

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Breakpoint int

const (
	Bp_ExtraSmall Breakpoint = iota
	Bp_Small
	Bp_Medium
	Bp_Large
	Bp_ExtraLarge
	Bp_ExtraExtraLarge
)

type BreakpointLayout struct {
	Breakpoint Breakpoint
	Width      int
	Height     int
}

func (layout BreakpointLayout) DrawLayout(screen *ebiten.Image) {
	// Fill background with a color
	screen.Fill(color.RGBA{0, 0, 255, 255}) // Blue background

	// Draw a red rectangle to represent a UI element
	uiElement := ebiten.NewImage(100, 50)
	uiElement.Fill(color.RGBA{255, 0, 0, 255}) // Red rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(layout.Width/2-50), float64(layout.Height/2-25))
	screen.DrawImage(uiElement, op)
}

type BreakpointLayoutSystem struct {
	Layouts map[Breakpoint]BreakpointLayout
}

func NewLayoutSystem() *BreakpointLayoutSystem {
	return &BreakpointLayoutSystem{
		Layouts: map[Breakpoint]BreakpointLayout{
			Bp_ExtraSmall:      {Breakpoint: Bp_ExtraSmall, Width: 320, Height: 480},
			Bp_Small:           {Breakpoint: Bp_Small, Width: 576, Height: 768},
			Bp_Medium:          {Breakpoint: Bp_Medium, Width: 768, Height: 1024},
			Bp_Large:           {Breakpoint: Bp_Large, Width: 992, Height: 1280},
			Bp_ExtraLarge:      {Breakpoint: Bp_ExtraLarge, Width: 1200, Height: 1600},
			Bp_ExtraExtraLarge: {Breakpoint: Bp_ExtraExtraLarge, Width: 1400, Height: 1920},
		},
	}
}

func (ls *BreakpointLayoutSystem) GetLayout(breakpoint Breakpoint) BreakpointLayout {
	return ls.Layouts[breakpoint]
}
