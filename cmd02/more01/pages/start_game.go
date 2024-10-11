package pages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more01/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

type StartGamePage struct {
	mainUI     *responsive.UI
	sidebarUI  *responsive.UI
	prevWidth  int
	prevHeight int
}

const sidebarFixedWidth = 200

func NewStartGamePage(switchPage func(pageName string)) *StartGamePage {

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	mainButtons := []*responsive.Button{
		responsive.NewButton("Play", func() { log.Println("Play clicked") }),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Level01", func() { log.Println("Level01 clicked") /* Add Level01 logic here */ }),
		responsive.NewButton("Level02", func() { log.Println("Level02 clicked") /* Add Level02 logic here */ }),
		responsive.NewButton("Level03", func() { log.Println("Level03 clicked") /* Add Level03 logic here */ }),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			switchPage("main")
		}),
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	screenWidth, screenHeight := 800, 600
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	return &StartGamePage{
		mainUI:     mainUI,
		sidebarUI:  sidebarUI,
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

	p.mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	p.sidebarUI.Update(sidebarFixedWidth, screenHeight)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x < sidebarFixedWidth {

			p.sidebarUI.HandleClick(x, y)
		} else {

			p.mainUI.HandleClick(x-sidebarFixedWidth, y)
		}
	}

	return nil
}

func (p *StartGamePage) Draw(screen *ebiten.Image) {

	screenWidth, screenHeight := screen.Size()

	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF})

	sidebarImage := ebiten.NewImage(sidebarFixedWidth, screenHeight)
	sidebarImage.Fill(color.RGBA{0x2E, 0x2E, 0x2E, 0xFF})

	p.sidebarUI.Draw(sidebarImage)

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(sidebarImage, op)

	mainUIImage := ebiten.NewImage(screenWidth-sidebarFixedWidth, screenHeight)
	mainUIImage.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF})

	p.mainUI.Draw(mainUIImage)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sidebarFixedWidth), 0)
	screen.DrawImage(mainUIImage, op)

	separatorColor := color.RGBA{0x00, 0x00, 0x00, 0xFF}
	separatorImg := ebiten.NewImage(2, screenHeight)
	separatorImg.Fill(separatorColor)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sidebarFixedWidth), 0)
	screen.DrawImage(separatorImg, op)
}

func (p *StartGamePage) HandleInput(x, y int) {
	if x < sidebarFixedWidth {
		p.sidebarUI.HandleClick(x, y)
	} else {
		p.mainUI.HandleClick(x-sidebarFixedWidth, y)
	}
}
