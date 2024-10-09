package customPages

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"example.com/menu/cmd02/more02/responsive"
	"example.com/menu/cmd02/more02/types"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// LevelGamePage represents the start game UI with a play-render-space for subpages.
type LevelGamePage struct {
	mainUI       *responsive.UI
	sidebarUI    *responsive.UI // Sidebar UI
	subPages     map[string]types.Page
	currentSub   types.Page
	prevWidth    int
	prevHeight   int
	sidebarWidth int
	switchPage   func(string) // Add this line
}

// Define a constant for sidebar width (for layout purposes)
const sidebarFixedWidth = 200 // Adjust this value as needed

// NewLevelGamePage initializes the start game page with specific breakpoints and buttons.
func NewLevelGamePage(switchPage func(string)) *LevelGamePage {
	// Main UI setup remains the same
	mainBreakpoints := []responsive.Breakpoint{
		{Width: 1200, LayoutMode: responsive.LayoutGrid},
		{Width: 800, LayoutMode: responsive.LayoutVertical},
		{Width: 0, LayoutMode: responsive.LayoutHorizontal},
	}
	mainButtons := []*responsive.Button{}

	sidebarButtons := []*responsive.Button{
		responsive.NewButton("Level 1", nil),
		responsive.NewButton("Level 2", nil),
		responsive.NewButton("Back", nil),
	}

	mainUI := responsive.NewUI("Start Game", mainBreakpoints, mainButtons)

	// Sidebar UI setup
	sidebarBreakpoints := []responsive.Breakpoint{
		{Width: 0, LayoutMode: responsive.LayoutVertical}, // Always vertical for sidebar
	}

	sidebarUI := responsive.NewUI("Menu", sidebarBreakpoints, sidebarButtons)

	// Initialize screen dimensions
	screenWidth, screenHeight := 800, 600
	mainUI.Update(screenWidth-sidebarFixedWidth, screenHeight)
	sidebarUI.Update(sidebarFixedWidth, screenHeight)

	// Initialize subpages
	subPages := make(map[string]types.Page)
	subPages["level01"] = NewLevel01Page()
	subPages["level02"] = NewLevel02Page()
	// Add more subpages as needed

	page := &LevelGamePage{
		mainUI:       mainUI,
		sidebarUI:    sidebarUI,
		subPages:     subPages,
		currentSub:   subPages["level01"], // Default subpage
		prevWidth:    screenWidth,
		prevHeight:   screenHeight,
		sidebarWidth: sidebarFixedWidth,
		switchPage:   switchPage, // Add this line
	}

	// Now that the page exists, set up the sidebar buttons
	page.setupSidebarButtons()

	// Add this line to reset button states when creating the page
	page.ResetAllButtonStates()

	return page
}

// setupSidebarButtons sets up the sidebar button callbacks
func (p *LevelGamePage) setupSidebarButtons() {
	p.sidebarUI.Buttons[0].OnClickFunc = func() { p.SwitchSubPage("level01") }
	p.sidebarUI.Buttons[1].OnClickFunc = func() { p.SwitchSubPage("level02") }
	p.sidebarUI.Buttons[2].OnClickFunc = func() {
		log.Println("Back clicked")
		p.switchPage("main") // Use the switchPage function to go back to main menu
	}
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
	if p.currentSub != nil {
		p.currentSub.Update()
	}

	return nil
}

// HandleInput processes input specific to the page (if any).
func (p *LevelGamePage) HandleInput(x, y int) {
	if x < p.sidebarWidth {
		p.sidebarUI.HandleClick(x, y)
	} else {
		// Pass the click to the current subpage
		if p.currentSub != nil {
			p.currentSub.HandleInput(x, y)
		}
	}
}

// SwitchSubPage switches the current subpage.
func (p *LevelGamePage) SwitchSubPage(pageName string) {
	if page, exists := p.subPages[pageName]; exists {
		log.Printf("Switching to subpage: %s\n", pageName)
		p.currentSub = page
		p.ResetAllButtonStates() // Reset button states when switching subpages
	} else {
		log.Printf("Subpage %s does not exist!\n", pageName)
	}
}

func (p *LevelGamePage) Draw(screen *ebiten.Image) {
	// Fill the background
	screen.Fill(color.RGBA{0x3E, 0x3E, 0x3E, 0xFF}) // Slightly lighter gray background

	// Draw the sidebar and main UI
	p.sidebarUI.Draw(screen)
	p.mainUI.Draw(screen)

	// Draw the current subpage in the play-render-space
	if p.currentSub != nil {
		screenWidth, screenHeight := screen.Size()
		// Create a subimage for the play-render-space
		playRenderSpace := ebiten.NewImage(screenWidth-p.sidebarWidth, screenHeight)
		p.currentSub.Draw(playRenderSpace)

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

// Add this method to the LevelGamePage struct
func (p *LevelGamePage) ResetAllButtonStates() {
	p.mainUI.ResetButtonStates()
	p.sidebarUI.ResetButtonStates()
	for _, subPage := range p.subPages {
		if subPageWithUI, ok := subPage.(interface{ ResetButtonStates() }); ok {
			subPageWithUI.ResetButtonStates()
		}
	}
}
