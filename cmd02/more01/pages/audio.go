package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

// AudioPage represents the audio settings UI.
type AudioPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

// NewAudioPage initializes the audio settings page with specific breakpoints and buttons.
func NewAudioPage(switchPage func(pageName string)) *AudioPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Volume Up", func() { log.Println("Volume Up clicked") /* Add Volume Up logic */ }),
		responsive.NewButton("Volume Down", func() { log.Println("Volume Down clicked") /* Add Volume Down logic */ }),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("settings")
		}),
	}

	ui := responsive.NewUI("Audio Settings", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &AudioPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

// Update updates the page state.
func (p *AudioPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("AudioPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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
func (p *AudioPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x4E, 0x4E, 0x4E, 0xFF}) // Even lighter gray background
	p.ui.Draw(screen)
}

// HandleInput processes input specific to the page (if any).
func (p *AudioPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
