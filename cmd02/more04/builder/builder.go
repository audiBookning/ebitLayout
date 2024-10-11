package builder

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more04/navigator"
	pagemodel "example.com/menu/cmd02/more04/pageModel"
	"example.com/menu/cmd02/more04/responsive"
)

func NewLevelGamePage(mainNav *navigator.Navigator, screenWidth, screenHeight int, id string, label string) *pagemodel.SidebarPageBase {

	subNav := navigator.NewNavigator(nil)

	level01 := NewSubPage01(subNav, screenWidth, screenHeight, "level01", "Level 01")
	level02 := NewSubPage02(subNav, screenWidth, screenHeight, "level02", "Level 02")

	subNav.AddPage("level01", level01)
	subNav.AddPage("level02", level02)

	subNav.SwitchTo("level01")

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainButtons := []*responsive.Button{}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Level 1", func() { subNav.SwitchTo("level01") }),
		responsive.NewButton("Level 2", func() { subNav.SwitchTo("level02") }),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }),
	}

	mainUI := responsive.NewUI(label, mainBreakpoints, mainButtons)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	const sidebarFixedWidth = 200
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	page := &pagemodel.SidebarPageBase{
		ID:            id,
		Label:         label,
		MainUI:        mainUI,
		SidebarUI:     sidebarUI,
		SubNavigator:  subNav,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		SidebarWidth:  sidebarFixedWidth,
		Navigator:     mainNav,
		BackgroundClr: color.RGBA{0x3E, 0x3E, 0x3E, 0xFF},
	}

	page.ResetAllButtonStates()

	return page
}

func NewSubPage01(subNav *navigator.Navigator, screenWidth, screenHeight int, id string, label string) *pagemodel.SubPage {
	breakpoints := []responsive.Breakpoint{
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	buttons := []*responsive.Button{
		responsive.NewButton("Play", func() {
			log.Println("Play Level 01")

		}),
		responsive.NewButton("Back to Start", func() {
			log.Println("Back to Start")
			subNav.SwitchTo("start")
		}),
	}
	ui := responsive.NewUI(label, breakpoints, buttons)
	ui.Update(screenWidth, screenHeight)
	return &pagemodel.SubPage{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
}
func NewSubPage02(subNav *navigator.Navigator, screenWidth, screenHeight int, id string, label string) *pagemodel.SubPage {
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
	return &pagemodel.SubPage{
		ID:            id,
		Label:         label,
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		BackgroundClr: color.RGBA{0x6E, 0x6E, 0x6E, 0xFF},
	}
}

func NewAudioPage(nv *navigator.Navigator, screenWidth, screenHeight int) *pagemodel.SinglePage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Volume Up", func() {
			log.Println("Volume Up clicked")
		}),
		responsive.NewButton("Volume Down", func() {
			log.Println("Volume Down clicked")

		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}),
	}

	ui := responsive.NewUI("Audio Settings", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

	return &pagemodel.SinglePage{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x4E, 0x4E, 0x4E, 0xFF},
	}
}

func NewGraphicsPage(nv *navigator.Navigator, screenWidth, screenHeight int) *pagemodel.SinglePage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Resolution", func() {
			log.Println("Resolution clicked")

		}),
		responsive.NewButton("Fullscreen", func() {
			log.Println("Fullscreen clicked")

		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("settings")
		}),
	}

	ui := responsive.NewUI("Graphics Settings", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

	return &pagemodel.SinglePage{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x5E, 0x5E, 0x5E, 0xFF},
	}
}

func NewMainMenuPage(nv *navigator.Navigator, screenWidth, screenHeight int) *pagemodel.SinglePage {
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

	return &pagemodel.SinglePage{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x2E, 0x2E, 0x2E, 0xFF},
	}
}

func NewSettingsPage(nv *navigator.Navigator, screenWidth, screenHeight int) *pagemodel.SinglePage {
	breakpoints := []responsive.Breakpoint{
		{Width: 1000, LayoutMode: responsive.LayoutVertical},
		{Width: 600, LayoutMode: responsive.LayoutHorizontal},
	}

	buttons := []*responsive.Button{
		responsive.NewButton("Audio", func() {
			log.Println("Audio clicked")
			nv.SwitchTo("audio")
		}),
		responsive.NewButton("Graphics", func() {
			log.Println("Graphics clicked")
			nv.SwitchTo("graphics")
		}),
		responsive.NewButton("Back", func() {
			log.Println("Back clicked")
			nv.SwitchTo("main")
		}),
	}

	ui := responsive.NewUI("Settings", breakpoints, buttons)

	ui.Update(screenWidth, screenHeight)

	return &pagemodel.SinglePage{
		Ui:            ui,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		Navigator:     nv,
		BackgroundClr: color.RGBA{0x1F, 0x1F, 0x1F, 0xFF},
	}
}
