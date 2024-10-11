package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more02/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GraphicsPage struct {
	ui         *responsive.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

func NewGraphicsPage(switchPage func(pageName string)) *GraphicsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Resolution", func() { log.Println("Resolution clicked") }),
		responsive.NewButton("Fullscreen", func() { log.Println("Fullscreen clicked") }),
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

func (p *GraphicsPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("GraphicsPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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

func (p *GraphicsPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x5E, 0x5E, 0x5E, 0xFF})
	p.ui.Draw(screen)
}

func (p *GraphicsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
