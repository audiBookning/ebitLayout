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

type MainMenuPage struct {
	*pagemodel.SinglePageBase
}

func NewMainMenuPage(nv *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int) *MainMenuPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	fields := []types.Element{
		widgets.NewButton("Start Game", func() {
			log.Println("Start Game clicked")
			nv.SwitchTo("start")
		}, textWrapper),
		widgets.NewButton("Settings", func() {
			log.Println("Settings clicked")
			nv.SwitchTo("settings")
		}, textWrapper),
		widgets.NewButton("Exit", func() {
			log.Println("Exit clicked")
			nv.SwitchTo("exit")
		}, textWrapper),
	}

	ui := widgets.NewUI("Main Menu", breakpoints, fields, textWrapper, responsive.AlignCenter)

	ui.LayoutUpdate(screenWidth, screenHeight)

	page := &pagemodel.SinglePageBase{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x2E, 0x2E, 0x2E, 0xFF},
	}

	return &MainMenuPage{
		SinglePageBase: page,
	}
}
