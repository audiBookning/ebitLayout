package responsive

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	Text            string
	Position        Position
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

	titleBounds := text.BoundString(basicfont.Face7x13, u.Title.Text)
	u.Title.X = (screenWidth - titleBounds.Dx()) / 2
	u.Title.Y = 50
	//log.Printf("UI.Update: Title position=(%d, %d)\n", u.Title.X, u.Title.Y)

	positions := u.manager.CalculatePositions(screenWidth, screenHeight, u.elements)

	for i, btn := range u.Buttons {
		pos, exists := positions[u.elements[i]]
		if exists {
			btn.Position = pos
			//log.Printf("UI.Update: Button '%s' position=(%d, %d)\n", btn.Text, btn.Position.X, btn.Position.Y)
		}

		btn.Update()
	}
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
