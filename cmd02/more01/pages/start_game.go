package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

// StartGamePage represents the start game UI.
type StartGamePage struct {
	mainUI     *responsive.UI
	sidebarUI  *responsive.UI // Sidebar UI
	prevWidth  int
	prevHeight int
}

// Define a constant for sidebar width (for layout purposes)
const sidebarFixedWidth = 200 // Adjust this value as needed

// NewStartGamePage initializes the start game page with specific breakpoints and buttons.
func NewStartGamePage(switchPage func(pageName string)) *StartGamePage {
	// Main UI setup remains the same
	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	mainButtons := []*responsive.Button{
		responsive.NewButton("Play", func() { log.Println("Play clicked") /* Add Play logic here */ }),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("main")
		}),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	// Sidebar UI setup
	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical}, // Always vertical for sidebar
	}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Options", func() { log.Println("Options clicked") /* Add Options logic here */ }),
		responsive.NewButton("Help", func() { log.Println("Help clicked") /* Add Help logic here */ }),
		responsive.NewButton("Credits", func() { log.Println("Credits clicked") /* Add Credits logic here */ }),
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	// Initialize screen dimensions
	screenWidth, screenHeight := 800, 600
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	return &StartGamePage{
		mainUI:     mainUI,
		sidebarUI:  sidebarUI,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *StartGamePage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	// Check for window resize
	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("StartGamePage: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	// Update main UI and sidebar UI
	p.mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	p.sidebarUI.Update(sidebarFixedWidth, screenHeight)

	// Handle clicks for both UIs
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x < sidebarFixedWidth {
			// Adjust x coordinate for sidebar
			p.sidebarUI.HandleClick(x, y)
		} else {
			// Adjust x coordinate for main UI
			p.mainUI.HandleClick(x-sidebarFixedWidth, y)
		}
	}

	return nil
}

// Draw renders the page.
func (p *StartGamePage) Draw(screen *ebiten.Image) {
	// Get the current screen dimensions
	screenWidth, screenHeight := screen.Size()

	// Fill the background
	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF}) // Slightly lighter gray background

	// Create a sub-image for the sidebar
	sidebarImage := ebiten.NewImage(sidebarFixedWidth, screenHeight)
	sidebarImage.Fill(color.RGBA{0x2E, 0x2E, 0x2E, 0xFF}) // Darker background for sidebar

	// Draw the sidebar UI on the sub-image
	p.sidebarUI.Draw(sidebarImage)

	// Draw the sidebar sub-image onto the main screen
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(sidebarImage, op)

	// Create a sub-image for the main UI area
	mainUIImage := ebiten.NewImage(screenWidth-sidebarFixedWidth, screenHeight)
	mainUIImage.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF}) // Same as the overall background

	// Draw the main UI on the sub-image
	p.mainUI.Draw(mainUIImage)

	// Draw the main UI sub-image onto the main screen
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sidebarFixedWidth), 0)
	screen.DrawImage(mainUIImage, op)

	// Optional: Draw a separator line between sidebar and main UI
	separatorColor := color.RGBA{0x00, 0x00, 0x00, 0xFF} // Black
	separatorImg := ebiten.NewImage(2, screenHeight)
	separatorImg.Fill(separatorColor)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sidebarFixedWidth), 0)
	screen.DrawImage(separatorImg, op)
}

// HandleInput processes input specific to the page (if any).
func (p *StartGamePage) HandleInput(x, y int) {
	if x < sidebarFixedWidth {
		p.sidebarUI.HandleClick(x, y)
	} else {
		p.mainUI.HandleClick(x-sidebarFixedWidth, y)
	}
}
