package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more02/responsive"
	"example.com/menu/cmd02/more02/types"
	"github.com/hajimehoshi/ebiten/v2"
)

// Level01Page represents the first level of the game.
type Level01Page struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
}

// NewLevel01Page initializes the Level 01 page.
func NewLevel01Page() types.Page {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Play", func() { log.Println("Play Level 01") /* Add Play logic */ }),
		responsive.NewButton("Back to Start", func() { log.Println("Back to Start") /* Add navigation logic if needed */ }),
	}

	ui := responsive.NewUI("Level 01", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &Level01Page{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *Level01Page) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("Level01Page: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	p.ui.Update(screenWidth, screenHeight)

	// Remove the click handling from here, as it will be handled in HandleInput

	return nil
}

// Draw renders the page.
func (p *Level01Page) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x6E, 0x6E, 0x6E, 0xFF}) // Example background color
	p.ui.Draw(screen)
}

// HandleInput processes input specific to the page (if any).
func (p *Level01Page) HandleInput(x, y int) {
	// Apply the offset here
	sidebarWidth := 200 // This should match the sidebar width in LevelGamePage
	p.ui.HandleClick(x-sidebarWidth, y)
}
