package main

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Breakpoint defines a screen width and the corresponding layout mode.
type Breakpoint struct {
	Width      int    // Maximum width for this breakpoint
	LayoutMode string // e.g., "horizontal", "vertical", "grid"
}

type Button struct {
	Text          string
	X, Y          int
	Width, Height int
	Clicked       bool
	OnClickFunc   func()
}

// IsClicked checks if the given coordinates are within the button's area.
func (b *Button) IsClicked(x, y int) bool {
	return x >= b.X && x <= b.X+b.Width &&
		y >= b.Y && y <= b.Y+b.Height
}

// OnClick handles the button click event.
func (b *Button) OnClick() {
	if b.OnClickFunc != nil {
		b.Clicked = true
		b.OnClickFunc()
		// Reset Clicked state after action
		// For simplicity, we reset immediately. You can add delay if needed.
		b.Clicked = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	// Choose button color based on click state
	var btnColor color.Color
	if b.Clicked {
		btnColor = color.RGBA{0xFF, 0xA5, 0x00, 0xFF} // Orange when clicked
	} else {
		btnColor = color.RGBA{0x00, 0x7A, 0xCC, 0xFF} // Default button color
	}

	// Create button rectangle
	button := ebiten.NewImage(b.Width, b.Height)
	button.Fill(btnColor)

	// Draw button rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	screen.DrawImage(button, op)

	// Draw button text
	textBounds := text.BoundString(basicfont.Face7x13, b.Text)
	textX := b.X + (b.Width-textBounds.Dx())/2
	textY := b.Y + (b.Height+textBounds.Dy())/2
	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, color.White)
}

// UI manages all UI elements.
type UI struct {
	Title          string
	TitleX, TitleY int
	Button1        Button
	Button2        Button
	Breakpoints    []Breakpoint
	CurrentMode    string // Current layout mode based on window size
}

// NewUI initializes the UI with customizable breakpoints.
func NewUI(breakpoints []Breakpoint) *UI {
	// Sort breakpoints in descending order of Width for proper selection
	sort.Slice(breakpoints, func(i, j int) bool {
		return breakpoints[i].Width > breakpoints[j].Width
	})

	ui := &UI{
		Title:       "Responsive Layout",
		Breakpoints: breakpoints,
	}

	// Initialize Buttons
	ui.Button1 = Button{
		Text: "Button 1",
		OnClickFunc: func() {
			fmt.Println("Button 1 clicked!")
		},
	}
	ui.Button2 = Button{
		Text: "Button 2",
		OnClickFunc: func() {
			fmt.Println("Button 2 clicked!")
		},
	}

	return ui
}

// Update recalculates UI element positions and sizes based on screen dimensions.
func (u *UI) Update(screenWidth, screenHeight int) {
	// Update Title Position (always centered at top)
	titleBounds := text.BoundString(basicfont.Face7x13, u.Title)
	u.TitleX = (screenWidth - titleBounds.Dx()) / 2
	u.TitleY = 50

	// Determine Current Layout Mode based on breakpoints
	u.determineLayoutMode(screenWidth)

	// Adjust layout based on CurrentMode
	switch u.CurrentMode {
	case "horizontal":
		u.layoutHorizontal(screenWidth, screenHeight)
	case "grid":
		u.layoutGrid(screenWidth, screenHeight)
	case "vertical":
		u.layoutVertical(screenWidth, screenHeight)
	default:
		u.layoutVertical(screenWidth, screenHeight)
	}
}

// determineLayoutMode sets the current layout mode based on screen width and breakpoints.
func (u *UI) determineLayoutMode(screenWidth int) {
	for _, bp := range u.Breakpoints {
		if screenWidth <= bp.Width {
			u.CurrentMode = bp.LayoutMode
		} else {
			continue
		}
	}
	// If no breakpoint matched, use the first (largest) layout mode
	if u.CurrentMode == "" && len(u.Breakpoints) > 0 {
		u.CurrentMode = u.Breakpoints[0].LayoutMode
	}
}

// layoutHorizontal arranges UI elements horizontally.
func (u *UI) layoutHorizontal(screenWidth, screenHeight int) {
	// Define button size
	btnWidth := 200
	btnHeight := 50

	// Calculate positions for horizontal layout
	totalWidth := 2*btnWidth + 50 // 50px space between buttons
	startX := (screenWidth - totalWidth) / 2
	yPos := screenHeight - btnHeight - 50

	u.Button1.X = startX
	u.Button1.Y = yPos
	u.Button1.Width = btnWidth
	u.Button1.Height = btnHeight

	u.Button2.X = startX + btnWidth + 50
	u.Button2.Y = yPos
	u.Button2.Width = btnWidth
	u.Button2.Height = btnHeight
}

// layoutVertical arranges UI elements vertically.
func (u *UI) layoutVertical(screenWidth, screenHeight int) {
	// Define button size
	btnWidth := 150
	btnHeight := 40

	// Calculate positions for vertical layout
	totalHeight := 2*btnHeight + 20 // 20px space between buttons
	startY := screenHeight - totalHeight - 50

	u.Button1.X = (screenWidth - btnWidth) / 2
	u.Button1.Y = startY
	u.Button1.Width = btnWidth
	u.Button1.Height = btnHeight

	u.Button2.X = (screenWidth - btnWidth) / 2
	u.Button2.Y = startY + btnHeight + 20
	u.Button2.Width = btnWidth
	u.Button2.Height = btnHeight
}

// layoutGrid arranges UI elements in a grid layout.
func (u *UI) layoutGrid(screenWidth, screenHeight int) {
	// Example grid layout: 2x2 grid (for future scalability)
	// For now, arranging two buttons side by side with smaller spacing
	btnWidth := 180
	btnHeight := 45

	// Calculate positions for grid layout
	totalWidth := 2*btnWidth + 30 // 30px space between buttons
	startX := (screenWidth - totalWidth) / 2
	yPos := screenHeight - btnHeight - 60

	u.Button1.X = startX
	u.Button1.Y = yPos
	u.Button1.Width = btnWidth
	u.Button1.Height = btnHeight

	u.Button2.X = startX + btnWidth + 30
	u.Button2.Y = yPos
	u.Button2.Width = btnWidth
	u.Button2.Height = btnHeight
}

// HandleClick determines if any button was clicked and triggers its action.
func (u *UI) HandleClick(x, y int) {
	if u.Button1.IsClicked(x, y) {
		u.Button1.OnClick()
	}
	if u.Button2.IsClicked(x, y) {
		u.Button2.OnClick()
	}
}

func (u *UI) Draw(screen *ebiten.Image) {
	// Draw Title
	text.Draw(screen, u.Title, basicfont.Face7x13, u.TitleX, u.TitleY, color.White)

	// Draw Buttons
	u.Button1.Draw(screen)
	u.Button2.Draw(screen)
}
