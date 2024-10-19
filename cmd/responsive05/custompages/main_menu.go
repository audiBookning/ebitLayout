package custompages

import (
	"image/color"
	"log"

	"example.com/menu/cmd/responsive05/responsive"
	"example.com/menu/cmd/responsive05/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

type MainMenuPage struct {
	ui         *widgets.UI
	manager    *responsive.LayoutManager
	prevWidth  int
	prevHeight int
}

func NewMainMenuPage(switchPage func(pageName string)) *MainMenuPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*widgets.Button{
		widgets.NewButton("Start Game", func() {
			log.Println("Start Game clicked")
			switchPage("start")
		}),
		widgets.NewButton("Settings", func() {
			log.Println("Settings clicked")
			switchPage("settings")
		}),
		widgets.NewButton("Exit", func() {
			log.Println("Exit clicked")
			switchPage("exit")
		}),
	}

	ui := widgets.NewUI("Main Menu", breakpoints, buttons)

	screenWidth, screenHeight := 800, 600
	ui.Update(screenWidth, screenHeight)

	return &MainMenuPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
	}
}

func (p *MainMenuPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("MainMenuPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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

func (p *MainMenuPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x2E, 0x2E, 0x2E, 0xFF})
	p.ui.Draw(screen)
}

func (p *MainMenuPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}
