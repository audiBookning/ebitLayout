package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

// MainMenuPage represents the main menu UI.
type MainMenuPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

// NewMainMenuPage initializes the main menu page with specific breakpoints and buttons.
func NewMainMenuPage(switchPage func(pageName string)) *MainMenuPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Start Game", func() {
			log.Println("Start Game clicked")
			switchPage("start")
		}),
		responsive.NewButton("Settings", func() {
			log.Println("Settings clicked")
			switchPage("settings")
		}),
		responsive.NewButton("Exit", func() {
			log.Println("Exit clicked")
			switchPage("exit")
		}),
	}

	ui := responsive.NewUI("Main Menu", breakpoints, buttons)

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
