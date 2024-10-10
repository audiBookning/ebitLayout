package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/pagemodel"
	"example.com/menu/cmd02/more05/responsive"
)

type AudioPage struct {
	*pagemodel.SinglePageBase
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

	page := &pagemodel.SinglePageBase{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x4E, 0x4E, 0x4E, 0xFF},
	}
	return &AudioPage{
		SinglePageBase: page,
	}
}
