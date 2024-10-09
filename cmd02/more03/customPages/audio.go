package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// AudioPage represents the audio settings UI.
type AudioPage struct {
	ui         *responsive.UI
	prevWidth  int
	prevHeight int
	navigator  *navigator.Navigator
}

func NewAudioPage(nv *navigator.Navigator, screenWidth, screenHeight int) *AudioPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Volume Up", func() {
			log.Println("Volume Up clicked")
		}),
		responsive.NewButton("Volume Down", func() {
			log.Println("Volume Down clicked")

		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}),
	}

	ui := responsive.NewUI("Audio Settings", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

	return &AudioPage{
		ui:         ui,
		prevWidth:  screenWidth,
		prevHeight: screenHeight,
		navigator:  nv,
	}
}

func (p *AudioPage) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *AudioPage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("AudioPage: Window resized to %dx%d\n", screenWidth, screenHeight)
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

func (p *AudioPage) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x4E, 0x4E, 0x4E, 0xFF}) // Even lighter gray background
	p.ui.Draw(screen)
}

func (p *AudioPage) HandleInput(x, y int) {
	p.ui.HandleClick(x, y)
}

func (p *AudioPage) ResetButtonStates() {
	p.ui.ResetButtonStates()
}
