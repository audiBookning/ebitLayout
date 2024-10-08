package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

// GraphicsPage represents the graphics settings UI.
type GraphicsPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

// NewGraphicsPage initializes the graphics settings page with specific breakpoints and buttons.
func NewGraphicsPage(switchPage func(pageName string)) *GraphicsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Resolution", func() { log.Println("Resolution clicked") /* Add Resolution logic */ }),
		responsive.NewButton("Fullscreen", func() { log.Println("Fullscreen clicked") /* Add Fullscreen logic */ }),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("settings")
		}),
	}

	ui := responsive.NewUI("Graphics Settings", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &GraphicsPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *GraphicsPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("GraphicsPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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
func (p *GraphicsPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x5E, 0x5E, 0x5E, 0xFF}) // Light gray background
	p.ui.Draw(screen)
}

// HandleInput processes input specific to the page (if any).
func (p *GraphicsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
