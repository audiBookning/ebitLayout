package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more06/responsive"
	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SubPageBase struct {
	ID            string
	Label         string
	Ui            *responsive.UI
	PrevWidth     int
	PrevHeight    int
	BackgroundClr color.Color
}

func NewSubPageBase(textWrapper *textwrapper.TextWrapper, id, label string, screenWidth, screenHeight int) *SubPageBase {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		responsive.NewButton("Button 01", func() {
			log.Println("Button 01 clicked")
		}, textWrapper),
		responsive.NewButton("Button 02", func() {
			log.Println("Button 02 clicked")
		}, textWrapper),
	}

	ui := responsive.NewUI(label, breakpoints, fields, textWrapper, responsive.AlignCenter)
	ui.LayoutUpdate(screenWidth, screenHeight)

	return &SubPageBase{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
}

func (p *SubPageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SubPageBase (%s): Window resized to %dx%d\n", p.ID, outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.Ui.LayoutUpdate(p.PrevWidth, p.PrevHeight)
	}
	return outsideWidth, outsideHeight
}

func (p *SubPageBase) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.Ui.HandleClick(x, y)
	}
	return nil
}

func (p *SubPageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.Ui.Draw(screen)
}

func (p *SubPageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

func (p *SubPageBase) HandleInput(x, y int) {
	p.Ui.HandleClick(x, y)
}

func (p *SubPageBase) ResetFieldStates() {
	p.Ui.ResetFieldStates()
}
