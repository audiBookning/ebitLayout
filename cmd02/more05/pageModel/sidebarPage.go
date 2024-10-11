package pagemodel

import (
	"image/color"
	"log"

	"example.com/menu/cmd02/more05/navigator"
	"example.com/menu/cmd02/more05/responsive"
	"example.com/menu/cmd02/more05/textwrapper"
	"example.com/menu/cmd02/more05/types"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// SidebarPageBase manages the game levels with a sidebar for navigation.
// It serves as a base template for custom sidebar pages.
type SidebarPageBase struct {
	ID            string
	Label         string
	MainUI        *responsive.UI
	SidebarUI     *responsive.UI // Sidebar UI
	SubNavigator  *navigator.Navigator
	PrevWidth     int
	PrevHeight    int
	SidebarWidth  int
	Navigator     *navigator.Navigator // Main navigator, used to navigate to SidebarPage itself
	BackgroundClr color.Color
}

// NewSidebarPageBase initializes a new SidebarPageBase.
func NewSidebarPageBase(mainNav *navigator.Navigator, textWrapper *textwrapper.TextWrapper, id, label string, screenWidth, screenHeight int) *SidebarPageBase {
	// Initialize the sub-navigator for SidebarPage
	subNav := navigator.NewNavigator(nil) // No onExit needed for sub-navigator

	// Initialize Sub1 and Sub2 pages with sub-navigator
	sub1 := NewSubPageBase(textWrapper, "sub01", "Sub 1", screenWidth, screenHeight)
	sub2 := NewSubPageBase(textWrapper, "sub02", "Sub 2", screenWidth, screenHeight)

	// Add subpages to sub-navigator
	subNav.AddPage(sub1.ID, sub1)
	subNav.AddPage(sub2.ID, sub2)

	// Optionally, set an initial subpage if needed
	subNav.SwitchTo("sub01")

	// Main UI setup (could include additional buttons relevant to SidebarPage)
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
		responsive.NewButton("Sub 1", func() { subNav.SwitchTo("sub01") }, textWrapper),
		responsive.NewButton("Sub 2", func() { subNav.SwitchTo("sub02") }, textWrapper),
		responsive.NewButton("Back", func() { mainNav.SwitchTo("main") }, textWrapper),
	}

	sidebarUI := responsive.NewUI("Sidebar Menu", sidebarBreakpoints, sidebarFields, textWrapper, responsive.AlignCenter)

	const sidebarFixedWidth = 200
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	// Initialize SidebarPageBase
	page := &SidebarPageBase{
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

	return page
}

// Layout handles the layout of the sidebar page.
func (p *SidebarPageBase) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != p.PrevWidth || outsideHeight != p.PrevHeight {
		log.Printf("SidebarPageBase: Window resized to %dx%d\n", outsideWidth, outsideHeight)
		p.PrevWidth = outsideWidth
		p.PrevHeight = outsideHeight
		p.MainUI.Update(p.PrevWidth-p.SidebarWidth, p.PrevHeight)
		p.SidebarUI.Update(p.SidebarWidth, p.PrevHeight)
	}
	return outsideWidth, outsideHeight
}

// Update handles the update logic for the sidebar page.
func (p *SidebarPageBase) Update() error {
	// Handle clicks for both UIs
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.HandleInput(x, y)
	}

	// Update the current subpage
	p.SubNavigator.CurrentActivePage().Update()

	return nil
}

// HandleInput processes input events.
func (p *SidebarPageBase) HandleInput(x, y int) {
	if x < p.SidebarWidth {
		p.SidebarUI.HandleClick(x, y)
	} else {
		// Pass the click to the current subpage
		if p.SubNavigator.CurrentActivePage() != nil {
			p.SubNavigator.CurrentActivePage().HandleInput(x-p.SidebarWidth, y)
		}
	}
}

// Draw renders the sidebar page.
func (p *SidebarPageBase) Draw(screen *ebiten.Image) {
	p.DrawBackGround(screen)
	// Draw the sidebar and main UI
	p.SidebarUI.Draw(screen)
	p.MainUI.Draw(screen)

	// Draw the current subpage in the play-render-space
	if p.SubNavigator.CurrentActivePage() != nil {
		screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
		// Create a subimage for the play-render-space
		playRenderSpace := ebiten.NewImage(screenWidth-p.SidebarWidth, screenHeight)
		p.SubNavigator.CurrentActivePage().Draw(playRenderSpace)

		// Draw the subimage onto the main screen with translation
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.SidebarWidth), 0)
		screen.DrawImage(playRenderSpace, op)
	}

	// Optional: Draw a separator line between sidebar and main UI
	separatorColor := color.RGBA{0x00, 0x00, 0x00, 0xFF} // Black
	separatorImg := ebiten.NewImage(2, p.PrevHeight)
	separatorImg.Fill(separatorColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.SidebarWidth), 0)
	screen.DrawImage(separatorImg, op)
}

// DrawBackGround draws the background color.
func (p *SidebarPageBase) DrawBackGround(screen *ebiten.Image) {
	screen.Fill(p.BackgroundClr)
}

// ResetAllButtonStates resets the state of all buttons in both UIs and the current subpage.
func (p *SidebarPageBase) ResetAllButtonStates() {
	p.ResetFieldStates()
	if p.SubNavigator.CurrentActivePage() != nil {
		p.SubNavigator.CurrentActivePage().ResetFieldStates()
	}
}

// ResetButtonStates resets the state of all buttons.
func (p *SidebarPageBase) ResetFieldStates() {
	p.MainUI.ResetFieldStates()
	p.SidebarUI.ResetFieldStates()
}
