package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button01 struct {
	X, Y, Width, Height int
	Label               string
	Color               color.Color
	HoverColor          color.Color
	ClickColor          color.Color
	isHovered           bool
	isPressed           bool
	OnClickFunc         func()
}

func NewButton(x, y, width, height int, label string, color, hoverColor, clickColor color.Color, onClick func()) *Button01 {
	return &Button01{
		X: x, Y: y, Width: width, Height: height,
		Label:       label,
		Color:       color,
		HoverColor:  hoverColor,
		ClickColor:  clickColor,
		OnClickFunc: onClick,
	}
}

func (b *Button01) Draw(screen *ebiten.Image) {
	var drawColor color.Color
	if b.isPressed {
		drawColor = b.ClickColor
	} else if b.isHovered {
		drawColor = b.HoverColor
	} else {
		drawColor = b.Color
	}
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), drawColor, true)
	ebitenutil.DebugPrintAt(screen, b.Label, b.X+10, b.Y+10)
}

func (b *Button01) Contains(x, y int) bool {
	return x >= b.X && x <= b.X+b.Width &&
		y >= b.Y && y <= b.Y+b.Height
}

func (b *Button01) OnClick() {
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
	b.isPressed = false
}

func (b *Button01) OnMouseDown() {
	b.isPressed = true
}

func (b *Button01) SetHovered(isHovered bool) {
	b.isHovered = isHovered
}
