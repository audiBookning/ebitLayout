package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd/responsive05/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

// StartGamePage represents the start game UI.
type StartGamePage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

// NewStartGamePage initializes the start game page with specific breakpoints and buttons.
func NewStartGamePage(switchPage func(pageName string)) *StartGamePage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Play", func() { log.Println("Play clicked") /* Add Play logic here */ }),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("main")
		}),
	}

	ui := responsive.NewUI("Start Game", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &StartGamePage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *StartGamePage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("StartGamePage: Window resized to %dx%d\n", screenWidth, screenHeight)
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
func (p *StartGamePage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF}) // Slightly lighter gray background
	p.ui.Draw(screen)
}

// HandleInput processes input specific to the page (if any).
func (p *StartGamePage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
