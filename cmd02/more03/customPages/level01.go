package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"example.com/menu/cmd02/more03/types"
	"github.com/hajimehoshi/ebiten/v2"
)

// Level01Page represents the first level of the game.
type Level01Page struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
	// No need for a navigator here since it's managed by LevelGamePage's subNavigator
}

func NewLevel01Page(subNav *navigator.Navigator, screenWidth, screenHeight int) types.Page {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Play", func() {
			log.Println("Play Level 01")
			// Implement Play Level 01 logic here
		}),
		responsive.NewButton("Back to Start", func() {
			log.Println("Back to Start")
			subNav.SwitchTo("start") // Navigate back to start within sub-navigator
		}),
	}

	ui := responsive.NewUI("Level 01", breakpoints, buttons)

	// Initialize screen dimensions

	ui.Update(screenWidth, screenHeight)

	return &Level01Page{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

func (p *Level01Page) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *Level01Page) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("Level01Page: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	p.ui.Update(screenWidth, screenHeight)

	return nil
}

func (p *Level01Page) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x6E, 0x6E, 0x6E, 0xFF}) // Example background color
	p.ui.Draw(screen)
}

func (p *Level01Page) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}

func (p *Level01Page) ResetButtonStates() {
	p.ui.ResetButtonStates()
}
