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

type GraphicsPage struct {
	*pagemodel.SinglePageBase
}

func NewGraphicsPage(nv *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int) *GraphicsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		widgets.NewButton("Resolution", func() {
			log.Println("Resolution clicked")

		}, textWrapper),
		widgets.NewButton("Fullscreen", func() {
			log.Println("Fullscreen clicked")

		}, textWrapper),
		widgets.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}, textWrapper),
	}

	ui := widgets.NewUI("Graphics Settings", breakpoints, fields, textWrapper, responsive.AlignCenter)

	ui.LayoutUpdate(screenWidth, screenHeight)

	page := &pagemodel.SinglePageBase{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x5E, 0x5E, 0x5E, 0xFF},
	}
	return &GraphicsPage{
		SinglePageBase: page,
	}
}
