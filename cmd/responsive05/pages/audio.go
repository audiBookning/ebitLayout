package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd/responsive05/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

type AudioPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

func NewAudioPage(switchPage func(pageName string)) *AudioPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Volume Up", func() { log.Println("Volume Up clicked") }),
		responsive.NewButton("Volume Down", func() { log.Println("Volume Down clicked") }),
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

func (p *AudioPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x4E, 0x4E, 0x4E, 0xFF})
	p.ui.Draw(screen)
}

func (p *AudioPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
