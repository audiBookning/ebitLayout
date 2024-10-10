package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/pagemodel"
	"example.com/menu/cmd02/more05/responsive"
)

type MainMenuPage struct {
	*pagemodel.SinglePageBase
}

func NewMainMenuPage(nv *navigator.Navigator, screenWidth, screenHeight int) *MainMenuPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Start Game", func() {
			log.Println("Start Game clicked")
			nv.SwitchTo("start")
		}),
		responsive.NewButton("Settings", func() {
			log.Println("Settings clicked")
			nv.SwitchTo("settings")
		}),
		responsive.NewButton("Exit", func() {
			log.Println("Exit clicked")
			nv.SwitchTo("exit")
		}),
	}

	ui := responsive.NewUI("Main Menu", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

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
