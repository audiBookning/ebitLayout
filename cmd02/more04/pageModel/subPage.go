package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more04/navigator"
	"example.com/menu/cmd02/more04/responsive"
	"github.com/hajimehoshi/ebiten/v2"
)

type SubPage struct {
	ID            string
	Label         string
	Ui            *responsive.UI
	PrevWidth     int
	PrevHeight    int
	BackgroundClr color.Color
}

func NewSubPage(subNav *navigator.Navigator, screenWidth, screenHeight int) *SubPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Button 01", func() {
			log.Println("Button 01 clicked")
		}),
		responsive.NewButton("Button 02", func() {
			log.Println("Button 02 clicked")
		}),
	}

	ui := responsive.NewUI("Sub Page", breakpoints, buttons)
	ui.Update(screenWidth, screenHeight)

	return &SubPage{
		ID:            "subPage",
		Label:         "Sub Page",
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
}

func (p *SubPage) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("Level01Page: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.Ui.Update(p.PrevWidth, p.PrevHeight)
	}
	return outsideWidth, outsideHeight
}

func (p *SubPage) Update() error {
	return nil
}

func (p *SubPage) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

func (p *SubPage) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.Ui.Draw(screen)
}

func (p *SubPage) HandleInput(x, y int) {
	p.Ui.HandleClick(x, y)
}

func (p *SubPage) ResetButtonStates() {
	p.Ui.ResetButtonStates()
}
