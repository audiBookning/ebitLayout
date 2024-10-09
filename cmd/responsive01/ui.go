package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type UI struct {
	Title         string
	ButtonText    string
	ButtonClicked bool
	// Positions and sizes
	TitleX, TitleY            int
	ButtonX, ButtonY          int
	ButtonWidth, ButtonHeight int
}

func NewUI() *UI {
	return &UI{
		Title:      "Responsive Layout",
		ButtonText: "Click Me",
	}
}

// Update recalculates UI element positions and sizes based on screen dimensions.
func (u *UI) Update(screenWidth, screenHeight int) {
	// Title positioned at top center
	titleBounds := text.BoundString(basicfont.Face7x13, u.Title)
	u.TitleX = (screenWidth - titleBounds.Dx()) / 2
	u.TitleY = 50

	// Button positioned at bottom center
	u.ButtonWidth = int(200 * float64(screenWidth) / 800)  // Scale with width
	u.ButtonHeight = int(50 * float64(screenHeight) / 600) // Scale with height
	u.ButtonX = (screenWidth - u.ButtonWidth) / 2
	u.ButtonY = screenHeight - u.ButtonHeight - 50
}

func (u *UI) Draw(screen *ebiten.Image) {
	// Draw title
	text.Draw(screen, u.Title, basicfont.Face7x13, u.TitleX, u.TitleY, color.White)

	// Draw button (simple rectangle with text)
	button := ebiten.NewImage(u.ButtonWidth, u.ButtonHeight)
	if u.ButtonClicked {
		button.Fill(color.RGBA{0xFF, 0x45, 0x00, 0xFF}) // Orange color when clicked
	} else {
		button.Fill(color.RGBA{0x00, 0x7A, 0xCC, 0xFF}) // Button color
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(u.ButtonX), float64(u.ButtonY))
	screen.DrawImage(button, op)

	// Draw button text
	textBounds := text.BoundString(basicfont.Face7x13, u.ButtonText)
	textX := u.ButtonX + (u.ButtonWidth-textBounds.Dx())/2
	textY := u.ButtonY + (u.ButtonHeight+textBounds.Dy())/2 + 4 // Adjust for better vertical alignment
	text.Draw(screen, u.ButtonText, basicfont.Face7x13, textX, textY, color.White)
}

// IsButtonClicked checks if the mouse coordinates are within the button's area.
func (u *UI) IsButtonClicked(x, y int) bool {
	return x >= u.ButtonX && x <= u.ButtonX+u.ButtonWidth &&
		y >= u.ButtonY && y <= u.ButtonY+u.ButtonHeight
}

// OnButtonClick handles the button click event.
func (u *UI) OnButtonClick() {
	u.ButtonClicked = true
	fmt.Println("Button clicked!")

	// Reset the button state after a short delay
	go func() {
		// Simple delay using a loop; for more accurate timing, use time.Sleep
		for i := 0; i < 100000000; i++ { // Dummy loop for delay
		}
		u.ButtonClicked = false
	}()
}
