package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SettingsPage represents the settings UI.
type SettingsPage struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
	navigator  *navigator.Navigator
}

// NewSettingsPage initializes the settings page with specific breakpoints and buttons.
func NewSettingsPage(nv *navigator.Navigator) *SettingsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Audio", func() {
			log.Println("Audio clicked")
			nv.SwitchTo("audio")
		}),
		responsive.NewButton("Graphics", func() {
			log.Println("Graphics clicked")
			nv.SwitchTo("graphics")
		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("main")
		}),
	}

	ui := responsive.NewUI("Settings", breakpoints, buttons)

	// Initialize screen dimensions
	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &SettingsPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
		navigator:  nv,
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
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

// HandleInput processes input specific to the page.
func (p *SettingsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}

// ResetButtonStates resets all button states.
func (p *SettingsPage) ResetButtonStates() {
	p.ui.ResetButtonStates()
}
