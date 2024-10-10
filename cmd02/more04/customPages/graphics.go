package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more04/navigator"
	"example.com/menu/cmd02/more04/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GraphicsPage represents the graphics settings UI.
type GraphicsPage struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
	navigator  *navigator.Navigator
}

func NewGraphicsPage(nv *navigator.Navigator, screenWidth, screenHeight int) *GraphicsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Resolution", func() {
			log.Println("Resolution clicked")
			// Add Resolution logic here
		}),
		responsive.NewButton("Fullscreen", func() {
			log.Println("Fullscreen clicked")
			// Add Fullscreen logic here
		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}),
	}

	ui := responsive.NewUI("Graphics Settings", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

	return &GraphicsPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
		navigator:  nv,
	}
}

func (p *GraphicsPage) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *GraphicsPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("GraphicsPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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

func (p *GraphicsPage) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	p.ui.Draw(screen)
}

func (p *GraphicsPage) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x5E, 0x5E, 0x5E, 0xFF})
}

func (p *GraphicsPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}

func (p *GraphicsPage) ResetButtonStates() {
	p.ui.ResetButtonStates()
}
