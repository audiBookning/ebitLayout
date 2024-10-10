package customPages

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more03/navigator"
	"example.com/menu/cmd02/more03/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// LevelGamePage manages the game levels with a sidebar for navigation.
type LevelGamePage struct {
	mainUI       *responsive.UI
	sidebarUI    *responsive.UI // Sidebar UI
	subNavigator *navigator.Navigator
	prevWidth    int
	prevHeight   int
	sidebarWidth int
	navigator    *navigator.Navigator // Main navigator, used to navigate to LevelGamePage itself
}

// Define a constant for sidebar width (for layout purposes)
const sidebarFixedWidth = 200 // Adjust this value as needed

func NewLevelGamePage(mainNav *navigator.Navigator, screenWidth, screenHeight int) *LevelGamePage {
	// Initialize the sub-navigator for LevelGamePage
	subNav := navigator.NewNavigator(nil) // No onExit needed for sub-navigator

	// Initialize Level01 and Level02 pages with sub-navigator
	level01 := NewLevel01Page(subNav, screenWidth, screenHeight)
	level02 := NewLevel02Page(subNav, screenWidth, screenHeight)

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
	mainButtons := []*responsive.Button{}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Level 1", func() { subNav.SwitchTo("level01") }),
		responsive.NewButton("Level 2", func() { subNav.SwitchTo("level02") }),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	// Sidebar UI setup
	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical}, // Always vertical for sidebar
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	// Initialize screen dimensions

	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	// Initialize LevelGamePage
	page := &LevelGamePage{
		mainUI:       mainUI,
		sidebarUI:    sidebarUI,
		subNavigator: subNav,
		prevWidth:    screenWidth,
		prevHeight:   screenHeight,
		sidebarWidth: sidebarFixedWidth,
		navigator:    mainNav, // Reference to main navigator to allow navigating back to main menu
	}

	// Reset button states when creating the page
	page.ResetAllButtonStates()

	return page
}

func (p *LevelGamePage) Layout(outsideWidth, outsideHeight int) (int, int) {
	p.prevWidth = outsideWidth
	p.prevHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func (p *LevelGamePage) Update() error {
	screenWidth, screenHeight := ebiten.WindowSize()

	// Check for window resize
	if screenWidth != p.prevWidth || screenHeight != p.prevHeight {
		log.Printf("LevelGamePage: Window resized to %dx%d\n", screenWidth, screenHeight)
		p.prevWidth = screenWidth
		p.prevHeight = screenHeight
	}

	// Update main UI and sidebar UI
	p.mainUI.Update(screenWidth-p.sidebarWidth, screenHeight)
	p.sidebarUI.Update(p.sidebarWidth, screenHeight)

	// Handle clicks for both UIs
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.HandleInput(x, y)
	}

	// Update the current subpage
	p.subNavigator.CurrentActivePage().Update()

	return nil
}

func (p *LevelGamePage) HandleInput(x, y int) {
	if x < p.sidebarWidth {
		p.sidebarUI.HandleClick(x, y)
	} else {
		// Pass the click to the current subpage
		if p.subNavigator.CurrentActivePage() != nil {
			p.subNavigator.CurrentActivePage().HandleInput(x-p.sidebarWidth, y)
		}
	}
}

func (p *LevelGamePage) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF})
}

func (p *LevelGamePage) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	// Draw the sidebar and main UI
	p.sidebarUI.Draw(screen)
	p.mainUI.Draw(screen)

	// Draw the current subpage in the play-render-space
	if p.subNavigator.CurrentActivePage() != nil {
		screenWidth, screenHeight := screen.Size()
		// Create a subimage for the play-render-space
		playRenderSpace := ebiten.NewImage(screenWidth-p.sidebarWidth, screenHeight)
		p.subNavigator.CurrentActivePage().Draw(playRenderSpace)

		// Draw the subimage onto the main screen with translation
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.sidebarWidth), 0)
		screen.DrawImage(playRenderSpace, op)
	}

	// Optional: Draw a separator line between sidebar and main UI
	separatorColor := color.RGBA{0x00, 0x00, 0x00, 0xFF} // Black
	separatorImg := ebiten.NewImage(2, p.prevHeight)
	separatorImg.Fill(separatorColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.sidebarWidth), 0)
	screen.DrawImage(separatorImg, op)
}

func (p *LevelGamePage) ResetAllButtonStates() {
	p.mainUI.ResetButtonStates()
	p.sidebarUI.ResetButtonStates()
	if p.subNavigator.CurrentActivePage() != nil {
		p.subNavigator.CurrentActivePage().ResetButtonStates()
	}
}

func (p *LevelGamePage) ResetButtonStates() {
	p.mainUI.ResetButtonStates()
	p.sidebarUI.ResetButtonStates()
	if p.subNavigator.CurrentActivePage() != nil {
		p.subNavigator.CurrentActivePage().ResetButtonStates()
	}
}
