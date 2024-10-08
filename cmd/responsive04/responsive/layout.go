package responsive

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Button represents a clickable button.
type Button struct {
	Text            string
	Position        Position
	Clicked         bool
	OnClickFunc     func()
	mutex           sync.RWMutex
	clickCooldown   int
	currentCooldown int
}

// NewButton creates a new Button with the given text and click handler.
func NewButton(text string, onClick func()) *Button {
	return &Button{
		Text:        text,
		OnClickFunc: onClick,
	}
}

// IsClicked checks if the given coordinates are within the button's area.
func (b *Button) IsClicked(x, y int) bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height
}

// HandleClick triggers the button's click event.
func (b *Button) HandleClick() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.OnClickFunc != nil && b.clickCooldown == 0 {
		b.Clicked = true
		b.OnClickFunc()
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
}

// Draw renders the button on the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Check for hover
	x, y := ebiten.CursorPosition()
	isHover := x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height

	// Choose button color based on state
	var btnColor color.Color
	if b.Clicked {
		btnColor = color.RGBA{0xFF, 0xA5, 0x00, 0xFF} // Orange when clicked
	} else if isHover {
		btnColor = color.RGBA{0x00, 0x8B, 0x8B, 0xFF} // DarkCyan on hover
	} else {
		btnColor = color.RGBA{0x00, 0x7A, 0xCC, 0xFF} // Default button color
	}

	// Create button rectangle
	buttonImg := ebiten.NewImage(b.Position.Width, b.Position.Height)
	buttonImg.Fill(btnColor)

	// Draw button rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.Position.X), float64(b.Position.Y))
	screen.DrawImage(buttonImg, op)

	// Draw button text
	textBounds := text.BoundString(basicfont.Face7x13, b.Text)
	textX := b.Position.X + (b.Position.Width-textBounds.Dx())/2
	textY := b.Position.Y + (b.Position.Height+textBounds.Dy())/2 + 4 // Adjust for better vertical alignment
	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, color.White)
}

// Title represents a text title.
type Title struct {
	Text      string
	X, Y      int
	FontScale float64
}

// NewTitle creates a new Title with the given text.
func NewTitle(text string) *Title {
	return &Title{
		Text:      text,
		FontScale: 1.0,
	}
}

// Draw renders the title on the screen.
func (t *Title) Draw(screen *ebiten.Image) {
	// Example: Adjust FontScale based on screen width
	// This requires a font library that supports scaling, such as truetype
	// For simplicity, we'll keep it fixed
	text.Draw(screen, t.Text, basicfont.Face7x13, t.X, t.Y, color.White)
}

// UI represents a page's UI, containing a title and buttons.
type UI struct {
	Title    *Title
	Buttons  []*Button
	manager  *LayoutManager
	elements []string // Identifiers for layout positioning
	mutex    sync.RWMutex
}

// NewUI creates a new UI with the given title, breakpoints, and buttons.
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

// Update recalculates UI element positions based on screen dimensions.
// Update recalculates UI element positions based on screen dimensions and updates buttons.
func (u *UI) Update(screenWidth, screenHeight int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Update Title Position (always centered at top)
	titleBounds := text.BoundString(basicfont.Face7x13, u.Title.Text)
	u.Title.X = (screenWidth - titleBounds.Dx()) / 2
	u.Title.Y = 50
	//log.Printf("UI.Update: Title position=(%d, %d)\n", u.Title.X, u.Title.Y)

	// Calculate positions for buttons
	positions := u.manager.CalculatePositions(screenWidth, screenHeight, u.elements)

	// Assign positions to buttons and update their states
	for i, btn := range u.Buttons {
		pos, exists := positions[u.elements[i]]
		if exists {
			btn.Position = pos
			//log.Printf("UI.Update: Button '%s' position=(%d, %d)\n", btn.Text, btn.Position.X, btn.Position.Y)
		}
		// Call the button's Update method to handle cooldowns
		btn.Update()
	}
}

// HandleClick determines if any button was clicked and triggers its action.
func (u *UI) HandleClick(x, y int) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	for _, btn := range u.Buttons {
		if btn.IsClicked(x, y) {
			btn.HandleClick()
		}
	}
}

// Draw renders the UI elements on the screen.
func (u *UI) Draw(screen *ebiten.Image) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	// Draw Title
	u.Title.Draw(screen)

	// Draw Buttons
	for _, btn := range u.Buttons {
		btn.Draw(screen)
	}
}
