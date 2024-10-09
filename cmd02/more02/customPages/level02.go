package customPages

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Level02Page represents the second level's UI.
type Level02Page struct {
	title string
}

func NewLevel02Page() *Level02Page {
	return &Level02Page{
		title: "Level 02 - The Challenge",
	}
}

// Update updates the Level02Page state.
func (p *Level02Page) Update() error {
	// Implement any level-specific logic here
	return nil
}

// Draw renders the Level02Page on the given screen.
func (p *Level02Page) Draw(screen *ebiten.Image) {
	// Fill the background with a different color
	screen.Fill(color.RGBA{0x3C, 0x2B, 0x1A, 0xFF}) // Example color

	// Draw the title text in Yellow
	text.Draw(screen, p.title, basicfont.Face7x13, 50, 50, color.RGBA{255, 255, 0, 255})

	// Add more rendering logic as needed
}

// HandleInput processes input specific to the Level02Page.
func (p *Level02Page) HandleInput(x, y int) {
	// Implement input handling for Level02Page if necessary
	log.Printf("Level02Page received input at (%d, %d)\n", x, y)
}
