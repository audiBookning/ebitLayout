package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/responsive"
	"example.com/menu/cmd02/more05/textwrapper"
	"example.com/menu/cmd02/more05/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SinglePageBase struct {
	ID            string
	Label         string
	Ui            *responsive.UI
	PrevWidth     int
	PrevHeight    int
	Navigator     *navigator.Navigator
	BackgroundClr color.Color
}

func NewSinglePageBase(nv *navigator.Navigator, textWrapper *textwrapper.TextWrapper, id string, label string, screenWidth, screenHeight int) *SinglePageBase {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		responsive.NewButton("Button 01", func() {
			log.Println("Button 01 clicked")
		}, textWrapper),
		responsive.NewButton("Button 02", func() {
			log.Println("Button 02 clicked")
		}, textWrapper),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}, textWrapper),
	}

	ui := responsive.NewUI(label, breakpoints, fields, textWrapper, responsive.AlignCenter)
	ui.LayoutUpdate(screenWidth, screenHeight)

	return &SinglePageBase{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x4E, 0x4E, 0x4E, 0xFF},
	}
}

func (p *SinglePageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SinglePageBase: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.Ui.LayoutUpdate(p.PrevWidth, p.PrevHeight)
	}
	return outsideWidth, outsideHeight
}

func (p *SinglePageBase) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.Ui.HandleClick(x, y)
	}
	return nil
}

func (p *SinglePageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.Ui.Draw(screen)
}

func (p *SinglePageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

func (p *SinglePageBase) HandleInput(x, y int) {
	p.Ui.HandleClick(x, y)
}

func (p *SinglePageBase) ResetFieldStates() {
	p.Ui.ResetFieldStates()
}
