package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

// SettingsPage represents the settings UI.
type SettingsPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

// NewSettingsPage initializes the settings page with specific breakpoints and buttons.
func NewSettingsPage(switchPage func(pageName string)) *SettingsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Audio", func() {
			log.Println("Audio clicked")
			switchPage("audio")
		}),
		responsive.NewButton("Graphics", func() {
			log.Println("Graphics clicked")
			switchPage("graphics")
		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("main")
		}),
	}

	ui := responsive.NewUI("Settings", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &SettingsPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *SettingsPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("SettingsPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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
func (p *SettingsPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF}) // Dark gray background
	p.ui.Draw(screen)
}

// HandleInput processes input specific to the page (if any).
func (p *SettingsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
