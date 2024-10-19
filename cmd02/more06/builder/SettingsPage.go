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

type SettingsPage struct {
	*pagemodel.SinglePageBase
}

func NewSettingsPage(nv *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int) *SettingsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		widgets.NewButton("Audio", func() {
			log.Println("Audio clicked")
			nv.SwitchTo("audio")
		}, textWrapper),
		widgets.NewButton("Graphics", func() {
			log.Println("Graphics clicked")
			nv.SwitchTo("graphics")
		}, textWrapper),
		widgets.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("main")
		}, textWrapper),
	}

	ui := widgets.NewUI("Settings", breakpoints, fields, textWrapper, responsive.AlignCenter)

	ui.LayoutUpdate(screenWidth, screenHeight)

	page := &pagemodel.SinglePageBase{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x1F, 0x1F, 0x1F, 0xFF},
	}
	return &SettingsPage{
		SinglePageBase: page,
	}
}
