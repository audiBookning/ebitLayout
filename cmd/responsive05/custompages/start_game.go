package custompages

import (
	"image/color"
	"log"

	"example.com/menu/cmd/responsive05/responsive"
	"example.com/menu/cmd/responsive05/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

type StartGamePage struct {
	ui         *widgets.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

func NewStartGamePage(switchPage func(pageName string)) *StartGamePage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*widgets.Button{
		widgets.NewButton("Play", func() { log.Println("Play clicked") /* Add Play logic here */ }),
		widgets.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("main")
		}),
	}

	ui := widgets.NewUI("Start Game", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &StartGamePage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

func (p *StartGamePage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("StartGamePage: Window resized to %dx%d\n", screenWidth, screenHeight)
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

func (p *StartGamePage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF})
	p.ui.Draw(screen)
}

func (p *StartGamePage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
