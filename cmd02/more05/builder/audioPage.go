package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/pagemodel"
	"example.com/menu/cmd02/more05/responsive"
	"example.com/menu/cmd02/more05/textwrapper"
	"example.com/menu/cmd02/more05/types"
)

type AudioPage struct {
	*pagemodel.SinglePageBase
}

func NewAudioPage(nv *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int) *AudioPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		responsive.NewButton("Volume Up", func() {
			log.Println("Volume Up clicked")
		}, textWrapper),
		responsive.NewButton("Volume Down", func() {
			log.Println("Volume Down clicked")
		}, textWrapper),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}, textWrapper),
	}

	ui := responsive.NewUI("Audio Settings", breakpoints, fields, textWrapper, responsive.AlignCenter)

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
