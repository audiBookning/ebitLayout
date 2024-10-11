package responsive

import (
	"image/color"
	"sync"
	"time"

	"example.com/menu/cmd02/more05/textwrapper"
	"example.com/menu/cmd02/more05/types"
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
	padding         int // Add padding for the button
}

func NewButton(text string, onClick func(), tw *textwrapper.TextWrapper) *Button {
	return &Button{
		Text:        text,
		OnClickFunc: onClick,
		TextWrapper: tw,
		padding:     10, // Add some padding around the text
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
		b.currentCooldown = 10 // Frames to wait before next click
	}
}

// Update handles cooldown for button clicks.
func (b *Button) Update() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.currentCooldown > 0 {
		b.currentCooldown--
		if b.currentCooldown == 0 {
			b.Clicked = false
		}
	}
	// Automatically reset Clicked state after a short duration
	if b.Clicked && time.Now().UnixNano()-b.lastClickTime > int64(time.Millisecond*100) {
		b.Clicked = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Calculate button size
	width, height := b.calculateSize()

	// Check for hover
	x, y := ebiten.CursorPosition()
	isHover := x >= b.Position.X && x <= b.Position.X+width &&
		y >= b.Position.Y && y <= b.Position.Y+height

	// Choose button color based on state
	var btnColor color.Color
	if b.Clicked && time.Now().UnixNano()-b.lastClickTime < int64(time.Millisecond*100) {
		btnColor = color.RGBA{0xFF, 0xA5, 0x00, 0xFF} // Orange when clicked
	} else if isHover {
		btnColor = color.RGBA{0x00, 0x8B, 0x8B, 0xFF} // DarkCyan on hover
	} else {
		btnColor = color.RGBA{0x00, 0x7A, 0xCC, 0xFF} // Default button color
	}

	// Create button rectangle
	buttonImg := ebiten.NewImage(width, height)
	buttonImg.Fill(btnColor)

	// Draw button rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Position.X), float64(b.Position.Y))
	screen.DrawImage(buttonImg, op)

	// Draw button text
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

// GetSize returns the current size of the button.
func (b *Button) GetSize() (int, int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	textWidth, textHeight := b.TextWrapper.MeasureText(b.Text)
	width := int(textWidth) + b.padding*2
	height := int(textHeight) + b.padding*2
	return width, height
}
