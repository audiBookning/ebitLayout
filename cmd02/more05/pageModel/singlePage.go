package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SinglePageBase represents the base functionality for a single page.
// Custom pages can embed this struct and override its methods as needed.
type SinglePageBase struct {
	ID            string
	Label         string
	Ui            *responsive.UI
	PrevWidth     int
	PrevHeight    int
	Navigator     *navigator.Navigator
	BackgroundClr color.Color
}

// NewSinglePageBase initializes a new SinglePageBase.
func NewSinglePageBase(nv *navigator.Navigator, id string, label string, screenWidth, screenHeight int) *SinglePageBase {
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

	ui := responsive.NewUI(label, breakpoints, buttons)
	ui.Update(screenWidth, screenHeight)

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

// Layout handles the layout of the page.
func (p *SinglePageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SinglePageBase: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.Ui.Update(p.PrevWidth, p.PrevHeight)
	}

	return outsideWidth, outsideHeight
}

// Update handles the update logic for the page.
func (p *SinglePageBase) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.Ui.HandleClick(x, y)
	}
	return nil
}

// Draw renders the page.
func (p *SinglePageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.Ui.Draw(screen)
}

// DrawBackGround draws the background color.
func (p *SinglePageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

// HandleInput processes input events.
func (p *SinglePageBase) HandleInput(x, y int) {
	p.Ui.HandleClick(x, y)
}

// ResetButtonStates resets the state of all buttons.
func (p *SinglePageBase) ResetButtonStates() {
	p.Ui.ResetButtonStates()
}
