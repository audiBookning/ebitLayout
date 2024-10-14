package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ToggleButton02 struct {
	X, Y, Width, Height int
	Label               string
	DefaultColor        color.Color
	CurrentColor        color.Color
	ToggleColor         color.Color
	IsToggled           bool
	OnClickFunc         func()
}

func (b *ToggleButton02) OnMouseDown() {

	b.IsToggled = !b.IsToggled
	if b.IsToggled {
		b.CurrentColor = b.ToggleColor
	} else {
		b.CurrentColor = b.DefaultColor
	}
}

func (b *ToggleButton02) SetHovered(isHovered bool) {

}

func NewToggleButton02(
	x, y,
	width, height int,
	label string,
	defaultColor, toggleColor color.Color,
	onClick func()) *ToggleButton02 {
	return &ToggleButton02{
		X:            x,
		Y:            y,
		Width:        width,
		Height:       height,
		Label:        label,
		DefaultColor: defaultColor,
		CurrentColor: defaultColor,
		ToggleColor:  toggleColor,
		IsToggled:    false,
		OnClickFunc:  onClick,
	}
}

func (b *ToggleButton02) Contains(x, y int) bool {
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

func (b *ToggleButton02) OnClick() {
	b.IsToggled = !b.IsToggled
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
}

func (b *ToggleButton02) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(
		screen,
		float32(b.X), float32(b.Y),
		float32(b.Width), float32(b.Height),
		b.CurrentColor, true)
}
