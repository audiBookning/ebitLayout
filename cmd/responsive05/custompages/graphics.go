package custompages

import (
	"image/color"
	"log"

	"example.com/menu/cmd/responsive05/responsive"
	"example.com/menu/cmd/responsive05/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

type GraphicsPage struct {
	ui         *widgets.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

func NewGraphicsPage(switchPage func(pageName string)) *GraphicsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*widgets.Button{
		widgets.NewButton("Resolution", func() { log.Println("Resolution clicked") /* Add Resolution logic */ }),
		widgets.NewButton("Fullscreen", func() { log.Println("Fullscreen clicked") /* Add Fullscreen logic */ }),
		widgets.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("settings")
		}),
	}

	ui := widgets.NewUI("Graphics Settings", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &GraphicsPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

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

func (p *GraphicsPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x5E, 0x5E, 0x5E, 0xFF})
	p.ui.Draw(screen)
}

func (p *GraphicsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
