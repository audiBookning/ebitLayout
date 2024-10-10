package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/pagemodel"
	"example.com/menu/cmd02/more05/responsive"
)

type Level02Page struct {
	*pagemodel.SinglePageBase
}

func NewLevel02Page(subNav *navigator.Navigator, screenWidth, screenHeight int, id string, label string) *Level02Page {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	buttons := []*responsive.Button{
		responsive.NewButton("Start Challenge", func() {
			log.Println("Start Challenge in Level 02")
		}),
		responsive.NewButton("Back to Start", func() {
			log.Println("Back to Start from Level 02")
		}),
	}
	ui := responsive.NewUI(label, breakpoints, buttons)
	ui.Update(screenWidth, screenHeight)
	page := &pagemodel.SinglePageBase{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
	return &Level02Page{
		SinglePageBase: page,
	}
}
