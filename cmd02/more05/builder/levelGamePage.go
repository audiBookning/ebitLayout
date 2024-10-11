package builder

import (
	"image/color"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/pagemodel"
	"example.com/menu/cmd02/more05/responsive"
	"example.com/menu/cmd02/more05/textwrapper"
	"example.com/menu/cmd02/more05/types"
)

type PlayGamePage struct {
	*pagemodel.SidebarPageBase
}

func NewLevelGamePage(mainNav *navigator.Navigator, textWrapper *textwrapper.TextWrapper, screenWidth, screenHeight int, id string, label string) *PlayGamePage {
	// Initialize the sub-navigator for LevelGamePage
	subNav := navigator.NewNavigator(nil) // No onExit needed for sub-navigator

	// Initialize Level01 and Level02 pages with sub-navigator
	level01 := NewLevel01Page(subNav, textWrapper, screenWidth, screenHeight, "level01", "Level 01")
	level02 := NewLevel02Page(subNav, textWrapper, screenWidth, screenHeight, "level02", "Level 02")

	// Add subpages to sub-navigator
	subNav.AddPage("level01", level01)
	subNav.AddPage("level02", level02)

	// Optionally, set an initial subpage if needed
	subNav.SwitchTo("level01")

	// Main UI setup (could include additional buttons relevant to LevelGamePage)
	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainFields := []types.Element{}
	mainUI := responsive.NewUI(label, mainBreakpoints, mainFields, textWrapper, responsive.AlignCenter)

	// Sidebar UI setup
	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical}, // Always vertical for sidebar
	}
	sidebarFields := []types.Element{
		responsive.NewButton("Level 1", func() { subNav.SwitchTo("level01") }, textWrapper),
		responsive.NewButton("Level 2", func() { subNav.SwitchTo("level02") }, textWrapper),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }, textWrapper),
	}
	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarFields, textWrapper, responsive.AlignCenter)

	const sidebarFixedWidth = 200
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	// Initialize LevelGamePage
	page := &pagemodel.SidebarPageBase{
		ID:            id,
		Label:         label,
		MainUI:        mainUI,
		SidebarUI:     sidebarUI,
		SubNavigator:  subNav,
		PrevWidth:     screenWidth,
		PrevHeight:    screenHeight,
		SidebarWidth:  sidebarFixedWidth,
		Navigator:     mainNav, // Reference to main navigator to allow navigating back to main menu
		BackgroundClr: color.RGBA{0x3E, 0x3E, 0x3E, 0xFF},
	}

	// Reset button states when creating the page
	page.ResetAllButtonStates()

	return &PlayGamePage{
		SidebarPageBase: page,
	}
}
