package widgets

import (
	"image/color"
	"sync"

	"example.com/menu/cmd/responsive05/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	Text            string
	Position        responsive.Position
	Clicked         bool
	OnClickFunc     func()
	mutex           sync.RWMutex
	clickCooldown   int
	currentCooldown int
}

func NewButton(text string, onClick func()) *Button {
	return &Button{
		Text:        text,
		OnClickFunc: onClick,
	}
}

func (b *Button) IsClicked(x, y int) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height
}

func (b *Button) HandleClick() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.OnClickFunc != nil && b.clickCooldown == 0 {
		b.Clicked = true
		b.OnClickFunc()
		b.currentCooldown = 10
	}
}

func (b *Button) Update() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.currentCooldown > 0 {
		b.currentCooldown--
		if b.currentCooldown == 0 {
			b.Clicked = false
		}
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	x, y := ebiten.CursorPosition()
	isHover := x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height

	var btnColor color.Color
	if b.Clicked {
		btnColor = color.RGBA{0xFF, 0xA5, 0x00, 0xFF}
	} else if isHover {
		btnColor = color.RGBA{0x00, 0x8B, 0x8B, 0xFF}
	} else {
		btnColor = color.RGBA{0x00, 0x7A, 0xCC, 0xFF}
	}

	buttonImg := ebiten.NewImage(b.Position.Width, b.Position.Height)
	buttonImg.Fill(btnColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Position.X), float64(b.Position.Y))
	screen.DrawImage(buttonImg, op)

	textBounds := text.BoundString(basicfont.Face7x13, b.Text)
	textX := b.Position.X + (b.Position.Width-textBounds.Dx())/2
	textY := b.Position.Y + (b.Position.Height+textBounds.Dy())/2 + 4
	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, color.White)
}
