package responsive

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"example.com/menu/cmd02/more02/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Title struct {
	Text      string
	X, Y      int
	FontScale float64
}

func NewTitle(text string) *Title {
	return &Title{
		Text:      text,
		FontScale: 1.0,
	}
}

func (t *Title) Draw(screen *ebiten.Image) {

	text.Draw(screen, t.Text, basicfont.Face7x13, t.X, t.Y, color.White)
}

type Button struct {
	Text            string
	OnClickFunc     func()
	Position        types.Position
	Clicked         bool
	clickCooldown   int
	currentCooldown int
	mutex           sync.Mutex
	lastClickTime   int64
}

func NewButton(text string, onClick func()) *Button {
	return &Button{
		Text:        text,
		OnClickFunc: onClick,
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

	x, y := ebiten.CursorPosition()
	isHover := x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height

	var btnColor color.Color
	if b.Clicked && time.Now().UnixNano()-b.lastClickTime < int64(time.Millisecond*100) {
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

func (b *Button) ResetState() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Clicked = false
	b.currentCooldown = 0
}

type UI struct {
	Title    *Title
	Buttons  []*Button
	manager  *LayoutManager
	elements []string
	mutex    sync.RWMutex
}

func NewUI(titleText string, breakpoints []Breakpoint, buttons []*Button) *UI {
	elements := make([]string, len(buttons))
	for i := range buttons {
		elements[i] = fmt.Sprintf("button%d", i+1)
	}

	return &UI{
		Title:    NewTitle(titleText),
		Buttons:  buttons,
		manager:  NewLayoutManager(breakpoints),
		elements: elements,
	}
}

func (u *UI) Update(screenWidth, screenHeight int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.manager.DetermineLayout(screenWidth)

	positions := u.manager.CalculatePositions(screenWidth, screenHeight, u.elements)

	for i, btn := range u.Buttons {
		elem := u.elements[i]
		pos, exists := positions[elem]
		if exists {
			btn.Position = pos
		} else {

			btn.Position = types.Position{X: 0, Y: i * 50, Width: 100, Height: 40}
		}

		btn.Update()
	}

	titleBounds := text.BoundString(basicfont.Face7x13, u.Title.Text)
	u.Title.X = (screenWidth - titleBounds.Dx()) / 2
	u.Title.Y = 50
}

func (u *UI) HandleClick(x, y int) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	for _, btn := range u.Buttons {
		if btn.IsClicked(x, y) {
			btn.HandleClick()
		}
	}
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	u.Title.Draw(screen)

	for _, btn := range u.Buttons {
		btn.Draw(screen)
	}
}

func (u *UI) ResetButtonStates() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, btn := range u.Buttons {
		btn.ResetState()
	}
}
