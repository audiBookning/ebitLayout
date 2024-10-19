package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more06/navigator"
	"example.com/menu/cmd02/more06/pagemodel"
	"example.com/menu/cmd02/more06/responsive"
	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
	"example.com/menu/cmd02/more06/widgets"
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
		widgets.NewButton("Volume Up", func() {
			log.Println("Volume Up clicked")
		}, textWrapper),
		widgets.NewButton("Volume Down", func() {
			log.Println("Volume Down clicked")
		}, textWrapper),
		widgets.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}, textWrapper),
	}

	ui := widgets.NewUI("Audio Settings", breakpoints, fields, textWrapper, responsive.AlignCenter)

	ui.LayoutUpdate(screenWidth, screenHeight)

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
