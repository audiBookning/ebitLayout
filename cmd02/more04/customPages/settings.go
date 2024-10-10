package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SettingsPage represents the settings UI.
type SettingsPage struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
	navigator  *navigator.Navigator
}

func NewSettingsPage(nv *navigator.Navigator, screenWidth, screenHeight int) *SettingsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Audio", func() {
			log.Println("Audio clicked")
			nv.SwitchTo("audio")
		}),
		responsive.NewButton("Graphics", func() {
			log.Println("Graphics clicked")
			nv.SwitchTo("graphics")
		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("main")
		}),
	}

	ui := responsive.NewUI("Settings", breakpoints, buttons)

	// Initialize screen dimensions

	ui.Update(screenWidth, screenHeight)

	return &SettingsPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
		navigator:  nv,
	}
}

func (p *SettingsPage) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *SettingsPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("SettingsPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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

func (p *SettingsPage) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.ui.Draw(screen)
}

func (p *SettingsPage) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x1F, 0x1F, 0x1F, 0xFF})
}

func (p *SettingsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}

func (p *SettingsPage) ResetButtonStates() {
	p.ui.ResetButtonStates()
}
