package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more06/navigator"
	"example.com/menu/cmd02/more06/pagemodel"
	"example.com/menu/cmd02/more06/responsive"
	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
)

type Level01Page struct {
	*pagemodel.SinglePageBase
}

func NewLevel01Page(subNav *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int, id string, label string) *Level01Page {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	fields := []types.Element{
		responsive.NewButton("Play", func() {
			log.Println("Play Level 01")

		}, textWrapper),
		responsive.NewButton("Back to Start", func() {
			log.Println("Back to Start")
			subNav.SwitchTo("start")
		}, textWrapper),
	}
	ui := responsive.NewUI(label, breakpoints, fields, textWrapper, responsive.AlignCenter)
	ui.LayoutUpdate(screenWidth, screenHeight)
	page := &pagemodel.SinglePageBase{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
	return &Level01Page{
		SinglePageBase: page,
	}
}
