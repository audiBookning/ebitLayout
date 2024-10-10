package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more04/navigator"
	"example.com/menu/cmd02/more04/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SinglePage represents the audio settings UI.
type SinglePage struct {
	Ui            *responsive.UI
	PrevWidth     int
	PrevHeight    int
	Navigator     *navigator.Navigator
	BackgroundClr color.Color
}

func NewAudioPage(nv *navigator.Navigator, screenWidth, screenHeight int) *SinglePage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Button 01", func() {
			log.Println("Button 01 clicked")
		}),
		responsive.NewButton("Button 02", func() {
			log.Println("Button 02 clicked")

		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}),
	}

	ui := responsive.NewUI("Single Page", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

	return &SinglePage{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x4E, 0x4E, 0x4E, 0xFF},
	}
}

func (p *SinglePage) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SinglePage: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.Ui.Update(p.PrevWidth, p.PrevHeight)
	}

	return outsideWidth, outsideHeight
}

func (p *SinglePage) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.Ui.HandleClick(x, y)
	}

	return nil
}

func (p *SinglePage) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.Ui.Draw(screen)
}

func (p *SinglePage) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

func (p *SinglePage) HandleInput(x, y int) {
	p.Ui.HandleClick(x, y)
}

func (p *SinglePage) ResetButtonStates() {
	p.Ui.ResetButtonStates()
}
