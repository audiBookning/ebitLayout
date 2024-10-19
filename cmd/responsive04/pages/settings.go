package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd/responsive04/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

type SettingsPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

func NewSettingsPage(switchPage func(pageName string)) *SettingsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Audio", func() { log.Println("Audio clicked") }),
		responsive.NewButton("Graphics", func() { log.Println("Graphics clicked") }),
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

func (p *SettingsPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF})
	p.ui.Draw(screen)
}

func (p *SettingsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
