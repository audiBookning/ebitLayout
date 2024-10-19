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

type Level02Page struct {
	*pagemodel.SinglePageBase
}

func NewLevel02Page(subNav *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int, id string, label string) *Level02Page {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	fields := []types.Element{
		responsive.NewButton("Start Challenge", func() {
			log.Println("Start Challenge in Level 02")
		}, textWrapper),
		responsive.NewButton("Back to Start", func() {
			log.Println("Back to Start from Level 02")
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
		BackgroundClr: color.RGBA{0x3C, 0x2B, 0x1A, 0xFF},
	}
	return &Level02Page{
		SinglePageBase: page,
	}
}
