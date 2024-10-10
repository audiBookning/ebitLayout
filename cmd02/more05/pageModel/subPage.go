package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SubPageBase represents the base functionality for a subpage.
// Custom subpages can embed this struct and override its methods as needed.
type SubPageBase struct {
	ID            string
	Label         string
	Ui            *responsive.UI
	PrevWidth     int
	PrevHeight    int
	BackgroundClr color.Color
}

// NewSubPageBase initializes a new SubPageBase.
func NewSubPageBase(id, label string, screenWidth, screenHeight int) *SubPageBase {
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

	ui := responsive.NewUI(label, breakpoints, buttons)
	ui.Update(screenWidth, screenHeight)

	return &SubPageBase{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
}

// Layout handles the layout of the subpage.
func (p *SubPageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SubPageBase (%s): Window resized to %dx%d\n", p.ID, outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.Ui.Update(p.PrevWidth, p.PrevHeight)
	}
	return outsideWidth, outsideHeight
}

// Update handles the update logic for the subpage.
func (p *SubPageBase) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.Ui.HandleClick(x, y)
	}
	return nil
}

// Draw renders the subpage.
func (p *SubPageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.Ui.Draw(screen)
}

// DrawBackGround draws the background color.
func (p *SubPageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

// HandleInput processes input events.
func (p *SubPageBase) HandleInput(x, y int) {
	p.Ui.HandleClick(x, y)
}

// ResetButtonStates resets the state of all buttons.
func (p *SubPageBase) ResetButtonStates() {
	p.Ui.ResetButtonStates()
}

// SubPage is a custom subpage that embeds SubPageBase.
// Override methods as needed.
type SubPage struct {
	*SubPageBase
}

// NewSubPage initializes a new SubPage.
// Override methods by assigning custom functions after initialization.
func NewSubPage(id, label string, screenWidth, screenHeight int) *SubPage {
	base := NewSubPageBase(id, label, screenWidth, screenHeight)
	return &SubPage{
		SubPageBase: base,
	}
}

// Example of overriding the Draw method.
func (p *SubPage) Draw(screen *ebiten.Image) {
	// Custom draw logic
	log.Println("Custom SubPage Draw")
	p.SubPageBase.Draw(screen)
}
