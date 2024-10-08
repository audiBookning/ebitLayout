package responsive

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// MainMenuPage represents the main menu UI.
type MainMenuPage struct {
	ui         *UI
	manager    *LayoutManager
	prevWidth  int
	prevHeight int
}

// NewMainMenuPage initializes the main menu page with specific breakpoints and buttons.
func NewMainMenuPage() *MainMenuPage {
	breakpoints := []Breakpoint{
		{Width: 1200, LayoutMode: LayoutGrid},
		{Width: 800, LayoutMode: LayoutVertical},
		{Width: 0, LayoutMode: LayoutHorizontal},
	}

	buttons := []*Button{
		NewButton("Start Game", func() { log.Println("Start Game clicked") }),
		NewButton("Settings", func() { log.Println("Settings clicked") }),
		NewButton("Exit", func() { log.Println("Exit clicked") }),
	}

	ui := NewUI("Main Menu", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &MainMenuPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *MainMenuPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("MainMenuPage: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	p.ui.Update(screenWidth, screenHeight)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.ui.HandleClick(x, y)
	}

	return nil
}

// Draw renders the page.
func (p *MainMenuPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x2E, 0x2E, 0x2E, 0xFF}) // Slightly lighter gray background
	p.ui.Draw(screen)
}

// HandleInput processes input specific to the page (if any).
func (p *MainMenuPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
