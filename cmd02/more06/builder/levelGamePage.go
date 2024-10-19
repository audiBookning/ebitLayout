package builder

import (
	"image/color"

	"example.com/menu/cmd02/more06/navigator"
	"example.com/menu/cmd02/more06/pagemodel"
	"example.com/menu/cmd02/more06/responsive"
	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
)

type PlayGamePage struct {
	*pagemodel.SidebarPageBase
}

func NewLevelGamePage(mainNav *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int, id string, label string) *PlayGamePage {

	subNav := navigator.NewNavigator(nil)

	level01 := NewLevel01Page(subNav, textWrapper, screenWidth, screenHeight, "level01", "Level 01")
	level02 := NewLevel02Page(subNav, textWrapper, screenWidth, screenHeight, "level02", "Level 02")

	subNav.AddPage("level01", level01)
	subNav.AddPage("level02", level02)

	subNav.SwitchTo("level01")

	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainFields := []types.Element{}
	mainUI := responsive.NewUI(label, mainBreakpoints, mainFields, textWrapper, responsive.AlignCenter)

	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical},
	}
	sidebarFields := []types.Element{
		responsive.NewButton("Level 1", func() { subNav.SwitchTo("level01") }, textWrapper),
		responsive.NewButton("Level 2", func() { subNav.SwitchTo("level02") }, textWrapper),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }, textWrapper),
	}
	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarFields, textWrapper, responsive.AlignCenter)

	const sidebarFixedWidth = 200
	mainUI.LayoutUpdate(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.LayoutUpdate(sidebarFixedWidth, screenHeight)

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

	return &PlayGamePage{
		SidebarPageBase: page,
	}
}
