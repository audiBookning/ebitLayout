package widgets

import (
	"image/color"
	"sync"
	"time"

	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	Text            string
	OnClickFunc     func()
	Position        types.Position
	Clicked         bool
	currentCooldown int
	mutex           sync.Mutex
	lastClickTime   int64
	TextWrapper     *textwrapper.TextWrapper
	padding         int
}

func NewButton(text string, onClick func(), tw *textwrapper.TextWrapper) *Button {
	return &Button{
		Text:        text,
		OnClickFunc: onClick,
		TextWrapper: tw,
		padding:     10,
	}
}

func (b *Button) IsClicked(x, y int) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height
}

func (b *Button) HandleClick() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	currentTime := time.Now().UnixNano()
	if b.OnClickFunc != nil && currentTime-b.lastClickTime > int64(time.Millisecond*200) {
		b.Clicked = true
		b.OnClickFunc()
		b.lastClickTime = currentTime
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

	if b.Clicked && time.Now().UnixNano()-b.lastClickTime > int64(time.Millisecond*100) {
		b.Clicked = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	width, height := b.calculateSize()

	x, y := ebiten.CursorPosition()
	isHover := x >= b.Position.X && x <= b.Position.X+width &&
		y >= b.Position.Y && y <= b.Position.Y+height

	var btnColor color.Color
	if b.Clicked && time.Now().UnixNano()-b.lastClickTime < int64(time.Millisecond*100) {
		btnColor = color.RGBA{0xFF, 0xA5, 0x00, 0xFF}
	} else if isHover {
		btnColor = color.RGBA{0x00, 0x8B, 0x8B, 0xFF}
	} else {
		btnColor = color.RGBA{0x00, 0x7A, 0xCC, 0xFF}
	}

	buttonImg := ebiten.NewImage(width, height)
	buttonImg.Fill(btnColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Position.X), float64(b.Position.Y))
	screen.DrawImage(buttonImg, op)

	textWidth, textHeight := b.TextWrapper.MeasureText(b.Text)
	textX := float64(b.Position.X) + (float64(width)-textWidth)/2
	textY := float64(b.Position.Y) + (float64(height)-textHeight)/2

	b.TextWrapper.DrawText(screen, b.Text, textX, textY)
}

func (b *Button) GetPosition() types.Position {
	return b.Position
}

func (b *Button) SetPosition(pos types.Position) {
	b.Position = pos
}

func (b *Button) ResetState() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Clicked = false
	b.currentCooldown = 0
}

func (b *Button) calculateSize() (int, int) {
	textWidth, textHeight := b.TextWrapper.MeasureText(b.Text)
	width := int(textWidth) + b.padding*2
	height := int(textHeight) + b.padding*2
	return width, height
}

func (b *Button) GetSize() (int, int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	textWidth, textHeight := b.TextWrapper.MeasureText(b.Text)
	width := int(textWidth) + b.padding*2
	height := int(textHeight) + b.padding*2
	return width, height
}
