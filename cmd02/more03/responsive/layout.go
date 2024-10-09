package responsive

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"example.com/menu/cmd02/more03/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

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

// Button represents a clickable button.
type Button struct {
	Text            string
	OnClickFunc     func()
	Position        types.Position
	Clicked         bool
	clickCooldown   int
	currentCooldown int
	mutex           sync.Mutex
	lastClickTime   int64 // Add this field to track the last click time
}

// NewButton creates a new Button with the given text and click handler.
func NewButton(text string, onClick func()) *Button {
	return &Button{
		Text:        text,
		OnClickFunc: onClick,
	}
}

// IsClicked checks if the button was clicked based on x, y coordinates.
func (b *Button) IsClicked(x, y int) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height
}

// HandleClick triggers the button's click event.
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

// Draw renders the button on the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check for hover
	x, y := ebiten.CursorPosition()
	isHover := x >= b.Position.X && x <= b.Position.X+b.Position.Width &&
		y >= b.Position.Y && y <= b.Position.Y+b.Position.Height

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

// ResetState resets the button's state.
func (b *Button) ResetState() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Clicked = false
	b.currentCooldown = 0
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
func (u *UI) Update(screenWidth, screenHeight int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	// Determine the current layout based on screen width
	u.manager.DetermineLayout(screenWidth)

	// Calculate positions based on current layout
	positions := u.manager.CalculatePositions(screenWidth, screenHeight, u.elements)

	// Assign positions to buttons and update their states
	for i, btn := range u.Buttons {
		elem := u.elements[i]
		pos, exists := positions[elem]
		if exists {
			btn.Position = pos
		} else {
			// Set a default position if not calculated
			btn.Position = types.Position{X: 0, Y: i * 50, Width: 100, Height: 40}
		}
		// Call the button's Update method to handle cooldowns
		btn.Update()
	}

	// Update Title Position (centered at the top)
	titleBounds := text.BoundString(basicfont.Face7x13, u.Title.Text)
	u.Title.X = (screenWidth - titleBounds.Dx()) / 2
	u.Title.Y = 50
}

// HandleClick triggers button actions based on click coordinates.
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

// ResetButtonStates resets all button states.
func (u *UI) ResetButtonStates() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, btn := range u.Buttons {
		btn.ResetState()
	}
}
