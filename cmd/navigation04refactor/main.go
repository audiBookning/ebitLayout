package main

import (
	"image"
	"image/color"
	"log"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth     = 800
	screenHeight    = 600
	leftColumnWidth = 160 // Width of the static left column
)

var (
	filePathTxt          string
	Assets_Relative_Path = "../../"
	// Globals var and singletons... We could pass these around as parameters
	// but that seem overcomplicating this basic use case without good adavantages
	pageRegistry    = make(map[string]*Page)
	globalNavigator *Navigator
)

// Define constants for page IDs
const (
	PageIDRed     = "redPage"
	PageIDBlue    = "bluePage"
	PageIDGreen   = "greenPage"
	PageIDYellow  = "yellowPage"
	PageIDMagenta = "magentaPage"
)

// PageConfig struct to hold page configuration
type PageConfig struct {
	id      string
	color   color.Color
	message string
	x, y    float32
	width   float32
	height  float32
}

// PageButton struct to hold button configuration
// TODO: remove this extra struct. Does not help in this situation
type PageButton struct {
	targetPageID string
	x, y         float32
	label        string
}

// GetFilePath constructs the full path for asset files.
func GetFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

// Game represents the entire game state.
type Game struct {
	navigator     *Navigator
	lastKeyState  map[ebiten.Key]bool
	leftColumnMsg string
}

// NewGame initializes a new Game instance.
func NewGame(navigator *Navigator) *Game {
	return &Game{
		navigator:     navigator,
		lastKeyState:  make(map[ebiten.Key]bool),
		leftColumnMsg: "Static Left Column",
	}
}

// Update handles game logic, including keyboard inputs and navigator updates.
func (g *Game) Update() error {
	// Handle left arrow key press to pop the current content
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && !g.lastKeyState[ebiten.KeyArrowLeft] {
		g.navigator.Pop()
		g.lastKeyState[ebiten.KeyArrowLeft] = true
	}

	// Reset key state when not pressed
	if !ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.lastKeyState[ebiten.KeyArrowLeft] = false
	}

	// Define navigatorOffsetX and navigatorOffsetY based on layout
	navigatorOffsetX := float32(leftColumnWidth)
	navigatorOffsetY := float32(0) // Assuming navigator starts at top

	// Delegate update to Navigator
	_, err := g.navigator.Update(navigatorOffsetX, navigatorOffsetY)
	if err != nil {
		return err
	}

	return nil
}

// Draw renders the game screen, including the left column and navigator content.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the static left column
	leftColumn := ebiten.NewImage(leftColumnWidth, screenHeight)
	leftColumn.Fill(color.RGBA{50, 50, 50, 255}) // Dark gray color
	ebitenutil.DebugPrintAt(leftColumn, g.leftColumnMsg, 10, 10)
	screen.DrawImage(leftColumn, nil)

	// Define the navigator area rectangle
	navigatorAreaRect := image.Rect(leftColumnWidth, 0, screenWidth, screenHeight)

	// Delegate the drawing of the navigator area to the Navigator
	g.navigator.Draw(screen, navigatorAreaRect)
}

// Layout defines the game's screen dimensions.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Navigator manages a stack of pages for navigation with animation capabilities.
type Navigator struct {
	stack      []*Page
	animating  bool
	transition float64 // Animation progress (0 to 1)
	direction  int     // 1 for push, -1 for pop
	Push       func(*Page)
	Pop        func()
}

// NewNavigator initializes a new Navigator instance with animation support.
func NewNavigator() *Navigator {
	return &Navigator{
		stack:      []*Page{},
		animating:  false,
		transition: 1.0, // Start with no transition
	}
}

// PushPage pushes a new page onto the stack and starts the push animation.
func (n *Navigator) PushPage(page *Page) {
	if len(n.stack) > 0 {
		n.animating = true
		n.transition = 0.0
		n.direction = 1 // Push direction
	}
	n.stack = append(n.stack, page)
}

// PopPage pops the top page from the stack and starts the pop animation.
func (n *Navigator) PopPage() {
	if len(n.stack) > 1 && !n.animating {
		n.animating = true
		n.transition = 0.0
		n.direction = -1 // Pop direction
	}
}

// Update updates the navigator's state and handles animations.
func (n *Navigator) Update(navigatorOffsetX, navigatorOffsetY float32) (bool, error) {
	if len(n.stack) == 0 {
		return false, nil
	}

	// Handle animation transitions
	if n.animating {
		n.transition += 0.05 // Adjust this value for animation speed
		if n.transition >= 1.0 {
			n.transition = 1.0
			n.animating = false
			if n.direction == -1 {
				// Complete the pop after animation
				n.stack = n.stack[:len(n.stack)-1]
			}
		}
	}

	// Update the current page with navigator offsets and animation state
	currentPage := n.stack[len(n.stack)-1]
	err := currentPage.Update(navigatorOffsetX, navigatorOffsetY, n.animating)
	return n.animating, err
}

// Draw renders the navigator area and manages page animations.
func (n *Navigator) Draw(screen *ebiten.Image, navigatorAreaRect image.Rectangle) {
	// Create an off-screen image for the navigator area
	navigatorArea := ebiten.NewImage(navigatorAreaRect.Dx(), navigatorAreaRect.Dy())
	navigatorArea.Fill(color.RGBA{30, 30, 30, 255}) // Optional: Background color for navigator area

	if n.animating && len(n.stack) > 1 {
		var prevOffsetX, currentOffsetX float64

		if n.direction == 1 { // Push
			// Current Page slides to the left
			prevOffsetX = -n.transition * float64(navigatorAreaRect.Dx())
			// New Page slides in from the right
			currentOffsetX = float64(navigatorAreaRect.Dx()) * (1.0 - n.transition)
		} else if n.direction == -1 { // Pop
			// Current Page slides to the right
			currentOffsetX = float64(navigatorAreaRect.Dx()) * n.transition
			// Previous Page slides in from the left
			prevOffsetX = -float64(navigatorAreaRect.Dx()) * (1.0 - n.transition)
		}

		// Debugging: Log animation state
		//fmt.Printf("PrevOffsetX: %f, CurrentOffsetX: %f\n", prevOffsetX, currentOffsetX)

		// Draw the Previous Page with its own Y position
		previousPage := n.stack[len(n.stack)-2]
		previousPage.Draw(navigatorArea, prevOffsetX, 0)

		// Draw the Current (New) Page with its own Y position
		currentPage := n.stack[len(n.stack)-1]
		currentPage.Draw(navigatorArea, currentOffsetX, 0)
	} else {
		// No animation, draw the top page normally
		if len(n.stack) > 0 {
			currentPage := n.stack[len(n.stack)-1]
			currentPage.Draw(navigatorArea, 0, 0)
		}
	}

	// Draw the navigator area onto the main screen within the defined rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(navigatorAreaRect.Min.X), float64(navigatorAreaRect.Min.Y))
	screen.DrawImage(navigatorArea, op)
}

// CurrentPage retrieves the current active page.
func (n *Navigator) CurrentPage() *Page {
	if len(n.stack) == 0 {
		return nil
	}
	return n.stack[len(n.stack)-1]
}

// UIElement interface for different UI components
type UIElement interface {
	Update(offsetX, offsetY float32, isAnimating bool)
	Draw(screen *ebiten.Image)
}

// Page represents a single page in the navigation stack.
// TOD: rename this to component?
// but components will have to be inside some kind of page...
type Page struct {
	X, Y            float32
	Width, Height   float32
	backgroundColor color.Color
	message         string
	elements        []UIElement
	textWrapper     *textwrapper.TextWrapper
	NextPageID      string
}

// NewPage creates a new Page instance with specified position and size.
func NewPage(
	bgColor color.Color,
	message string,
	tw *textwrapper.TextWrapper,
	x, y, width, height float32,
) *Page {
	return &Page{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		backgroundColor: bgColor,
		message:         message,
		elements:        make([]UIElement, 0),
		textWrapper:     tw,
		NextPageID:      "",
	}
}

// AddElement adds a new UI element to the page
func (p *Page) AddElement(element UIElement) {
	p.elements = append(p.elements, element)
}

func (p *Page) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error {
	for _, element := range p.elements {
		element.Update(navigatorOffsetX+p.X, navigatorOffsetY+p.Y, isAnimating)
	}
	return nil
}

// Draw renders the page content onto the provided navigatorArea with given offsets.
func (p *Page) Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64) {
	pageArea := ebiten.NewImage(int(p.Width), int(p.Height))
	pageArea.Fill(p.backgroundColor)
	ebitenutil.DebugPrintAt(pageArea, p.message, 10, 10)

	for _, element := range p.elements {
		element.Draw(pageArea)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X)+offsetX, float64(p.Y)+offsetY)
	navigatorArea.DrawImage(pageArea, op)
}

// Helper function to create and add a button to a page
func (page *Page) addButton(btn PageButton) {
	targetPage, exists := pageRegistry[btn.targetPageID]
	if !exists {
		log.Fatalf("Page with ID '%s' does not exist", btn.targetPageID)
	}

	onNext := func(target *Page) func() {
		return func() {
			if globalNavigator == nil {
				log.Println("Navigator not initialized")
				return
			}
			globalNavigator.PushPage(target)
		}
	}(targetPage)

	button := widgets.NewButtonStd(
		btn.x,
		btn.y,
		100,
		40,
		btn.label,
		page.textWrapper,
		color.RGBA{0, 128, 255, 255},
		color.White,
		16,
		onNext,
	)

	page.AddElement(button)
}

// registerPages initializes and registers all pages in the global registry.
func registerPages(tw *textwrapper.TextWrapper, navigatorAreaWidth, navigatorAreaHeight float32) {
	// Page configurations with constants for IDs
	pageConfigs := []PageConfig{
		{
			id:      PageIDRed,
			color:   color.RGBA{255, 0, 0, 255},
			message: "Red Page - Full Size",
			x:       0,
			y:       0,
			width:   navigatorAreaWidth,
			height:  navigatorAreaHeight,
		},
		{
			id:      PageIDBlue,
			color:   color.RGBA{0, 0, 255, 255},
			message: "Blue Page - Top Half",
			x:       0,
			y:       0,
			width:   navigatorAreaWidth,
			height:  navigatorAreaHeight / 2,
		},
		{
			id:      PageIDGreen,
			color:   color.RGBA{0, 255, 0, 255},
			message: "Green Page - Bottom Half",
			x:       0,
			y:       navigatorAreaHeight / 2,
			width:   navigatorAreaWidth,
			height:  navigatorAreaHeight / 2,
		},
		{
			id:      PageIDYellow,
			color:   color.RGBA{255, 255, 0, 255},
			message: "Yellow Page - Smaller Window",
			x:       navigatorAreaWidth / 4,
			y:       navigatorAreaHeight / 4,
			width:   navigatorAreaWidth / 2,
			height:  navigatorAreaHeight / 2,
		},
		{
			id:      PageIDMagenta,
			color:   color.RGBA{255, 0, 255, 255},
			message: "Magenta Page - Custom Size",
			x:       50,
			y:       50,
			width:   300,
			height:  200,
		},
	}

	// Create and register each page
	for _, cfg := range pageConfigs {
		page := NewPage(
			cfg.color,
			cfg.message,
			tw,
			cfg.x,
			cfg.y,
			cfg.width,
			cfg.height,
		)
		pageRegistry[cfg.id] = page
	}

	// After all pages are registered, add buttons
	pageRegistry[PageIDRed].addButton(PageButton{
		targetPageID: PageIDBlue,
		x:            (navigatorAreaWidth - 100) / 2,
		y:            navigatorAreaHeight - 60,
		label:        "To Blue",
	})
	pageRegistry[PageIDRed].addButton(PageButton{
		targetPageID: PageIDYellow,
		x:            20,
		y:            20,
		label:        "To Yellow",
	})

	pageRegistry[PageIDBlue].addButton(PageButton{
		targetPageID: PageIDGreen,
		x:            (navigatorAreaWidth - 100) / 2,
		y:            (navigatorAreaHeight / 2) - 60,
		label:        "To Green",
	})
	pageRegistry[PageIDBlue].addButton(PageButton{
		targetPageID: PageIDMagenta,
		x:            navigatorAreaWidth - 120,
		y:            20,
		label:        "To Magenta",
	})

	pageRegistry[PageIDGreen].addButton(PageButton{
		targetPageID: PageIDYellow,
		x:            (navigatorAreaWidth - 100) / 2,
		y:            (navigatorAreaHeight / 2) - 60,
		label:        "To Yellow",
	})
	pageRegistry[PageIDYellow].addButton(PageButton{
		targetPageID: PageIDMagenta,
		x:            (navigatorAreaWidth/4 - 100) / 2,
		y:            (navigatorAreaHeight / 2) - 60,
		label:        "To Magenta",
	})
	pageRegistry[PageIDMagenta].addButton(PageButton{
		targetPageID: PageIDRed,
		x:            100,
		y:            140,
		label:        "To Red",
	})
}

func main() {
	_, filePathTxt, _, _ = runtime.Caller(0)
	fontPath := GetFilePath("assets/fonts/roboto_regularTTF.ttf")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Navigator Example with Animations")

	navigator := NewNavigator()
	globalNavigator = navigator // Assign to the global variable

	// Initialize TextWrapper
	textWrapper, err := textwrapper.NewTextWrapper(fontPath, 16, false)
	if err != nil {
		log.Fatalf("Failed to create TextWrapper: %v", err)
	}

	// Calculate navigatorAreaWidth based on global constants
	navigatorAreaWidth := float32(screenWidth - leftColumnWidth)
	navigatorAreaHeight := float32(screenHeight)

	// Register all pages
	registerPages(textWrapper, navigatorAreaWidth, navigatorAreaHeight)

	// Initialize Push and Pop functions
	navigator.Push = func(page *Page) {
		navigator.PushPage(page)
	}
	navigator.Pop = func() {
		navigator.PopPage()
	}

	// Push the initial page
	initialPage, exists := pageRegistry[PageIDRed]
	if !exists {
		log.Fatal("Initial page 'redPage' not found in registry")
	}
	navigator.PushPage(initialPage)

	game := NewGame(navigator)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
