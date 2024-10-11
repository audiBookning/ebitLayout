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

type GraphicsPage struct {
	*pagemodel.SinglePageBase
}

func NewGraphicsPage(nv *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int) *GraphicsPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		responsive.NewButton("Resolution", func() {
			log.Println("Resolution clicked")
			// Add Resolution logic here
		}, textWrapper),
		responsive.NewButton("Fullscreen", func() {
			log.Println("Fullscreen clicked")
			// Add Fullscreen logic here
		}, textWrapper),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}, textWrapper),
	}

	ui := responsive.NewUI("Graphics Settings", breakpoints, fields, textWrapper, responsive.AlignCenter)

	ui.Update(screenWidth, screenHeight)

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
