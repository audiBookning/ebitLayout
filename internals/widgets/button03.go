package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ToggleButton03 struct {
	X, Y, Width, Height int
	Label               string
	DefaultColor        color.Color
	CurrentColor        color.Color
	ToggleColor         color.Color
	IsToggled           bool
	knobX               float64
	OnClickFunc         func()
}

func (b *ToggleButton03) OnMouseDown() {

	b.IsToggled = !b.IsToggled
	if b.IsToggled {
		b.CurrentColor = b.ToggleColor
	} else {
		b.CurrentColor = b.DefaultColor
	}

}

func (b *ToggleButton03) SetHovered(isHovered bool) {

}

func NewToggleButton03(
	x, y,
	width, height int,
	label string,
	color, toggleColor color.Color,
	onClick func()) *ToggleButton03 {
	return &ToggleButton03{
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		Label:        label,
		DefaultColor: color,
		ToggleColor:  toggleColor,
		knobX:        float64(x),
		OnClickFunc:  onClick,
	}
}

func (b *ToggleButton03) OnClick() {
	b.IsToggled = !b.IsToggled
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
}

func (b *ToggleButton03) Update() {
	targetX := float64(b.X)
	if b.IsToggled {
		targetX = float64(b.X + b.Width - b.Height)
	}
	b.knobX += (targetX - b.knobX) * 0.2
}

func (b *ToggleButton03) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), color.RGBA{200, 200, 200, 255}, true)

	knobSize := float32(b.Height)
	vector.DrawFilledRect(
		screen,
		float32(b.knobX), float32(b.Y),
		knobSize, knobSize,
		b.CurrentColor, true)
}

func (b *ToggleButton03) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}
