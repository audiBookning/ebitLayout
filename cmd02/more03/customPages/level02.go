package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"example.com/menu/cmd02/more03/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Level02Page represents the second level's UI.
type Level02Page struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
}

func NewLevel02Page(subNav *navigator.Navigator, screenWidth, screenHeight int) types.Page {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Start Challenge", func() {
			log.Println("Start Challenge in Level 02")
			// Implement Start Challenge logic here
		}),
		responsive.NewButton("Back to Start", func() {
			log.Println("Back to Start from Level 02")
			//subNav.SwitchTo("start") // Navigate back to start within sub-navigator
		}),
	}

	ui := responsive.NewUI("Level 02 - The Challenge", breakpoints, buttons)

	// Initialize screen dimensions

	ui.Update(screenWidth, screenHeight)

	return &Level02Page{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

func (p *Level02Page) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *Level02Page) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("Level02Page: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	p.ui.Update(screenWidth, screenHeight)

	return nil
}

func (p *Level02Page) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3C, 0x2B, 0x1A, 0xFF}) // Example color

	// Draw the title text in Yellow
	text.Draw(screen, "Level 02 - The Challenge", basicfont.Face7x13, 50, 50, color.RGBA{255, 255, 0, 255})

	p.ui.Draw(screen)
}

func (p *Level02Page) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}

func (p *Level02Page) ResetButtonStates() {
	p.ui.ResetButtonStates()
}
