package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ***** INPUT MANAGER *****
type InputManager struct {
	clickables       []Clickable
	screenClickables []Clickable
}

func (im *InputManager) Register(c Clickable) {
	im.clickables = append(im.clickables, c)
}

func (im *InputManager) RegisterScreenClickable(c Clickable) {
	im.screenClickables = append(im.screenClickables, c)
}

func (im *InputManager) ClearScreenClickable() {
	im.screenClickables = nil
}

func (im *InputManager) Update() {
	mouseX, mouseY := ebiten.CursorPosition()

	returnFlag := false

	for _, c := range im.clickables {
		if area, ok := c.(*ClickableArea); ok && area.Active {
			if area.Contains(mouseX, mouseY) {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					area.OnMouseDown()
				}
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
					area.OnClick()
				}

				break
			}
			returnFlag = true
		}
		breakFlag := checkClickable(c, mouseX, mouseY)
		if breakFlag {
			break
		}
	}

	if returnFlag {
		return
	}

	for _, c := range im.screenClickables {
		breakFlag := checkClickable(c, mouseX, mouseY)
		if breakFlag {
			break
		}
	}
}

func checkClickable(c Clickable, mouseX, mouseY int) bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if c.Contains(mouseX, mouseY) {
			c.OnMouseDown()
			return true
		}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if c.Contains(mouseX, mouseY) {
			c.OnClick()
			return true
		}
	}
	// Update hover state
	c.SetHovered(c.Contains(mouseX, mouseY))

	return false
}

// ***** CLICKABLE INTERFACE *****
type Clickable interface {
	Contains(x, y int) bool
	OnClick()
	OnMouseDown()
	SetHovered(isHovered bool)
}

// ***** BUTTON *****
type Button struct {
	X, Y, Width, Height int
	Label               string
	Color               color.Color
	HoverColor          color.Color
	ClickColor          color.Color
	isHovered           bool
	isPressed           bool
	OnClickFunc         func() // Callback function
	offsetY             int
	target              *GeneralScreen
}

func NewButton(x, y, width, height int, offsetY int, label string, color, hoverColor, clickColor color.Color, target *GeneralScreen, onClick func()) *Button {
	return &Button{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Label:       label,
		Color:       color,
		HoverColor:  hoverColor,
		ClickColor:  clickColor,
		OnClickFunc: onClick,
		offsetY:     offsetY,
		target:      target,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	var drawColor color.Color
	if b.isPressed {
		drawColor = b.ClickColor
	} else if b.isHovered {
		drawColor = b.HoverColor
	} else {
		drawColor = b.Color
	}
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), drawColor, true)
	ebitenutil.DebugPrintAt(screen, b.Label, b.X+10, b.Y+10)
}

func (b *Button) Contains(x, y int) bool {
	if x >= b.X && x <= b.X+b.Width &&
		y-b.offsetY >= b.Y && y-b.offsetY <= b.Y+b.Height {
		return true
	}
	return false
}

func (b *Button) OnClick() {
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
	b.isPressed = false
}

func (b *Button) OnMouseDown() {
	b.isPressed = true
}

func (b *Button) SetHovered(isHovered bool) {
	b.isHovered = isHovered
}

// ***** SCREEN INTERFACE *****
type Screen interface {
	Update(navigator *Navigator) error
	Draw(screen *ebiten.Image)
	RegisterClickables(inputManager *InputManager)
}

// ***** NAVIGATOR *****
type Navigator struct {
	stack        []Screen
	inputManager *InputManager
}

func NewNavigator(inputManager *InputManager) *Navigator {
	return &Navigator{
		stack:        make([]Screen, 0),
		inputManager: inputManager,
	}
}

func (n *Navigator) Push(screen Screen) {
	n.stack = append(n.stack, screen)
	screen.RegisterClickables(n.inputManager)
}

func (n *Navigator) Pop() {
	if len(n.stack) > 0 {
		n.stack = n.stack[:len(n.stack)-1]
	}
	if len(n.stack) > 0 {
		n.stack[len(n.stack)-1].RegisterClickables(n.inputManager)
	}
}

func (n *Navigator) Replace(screen Screen) {
	if len(n.stack) > 0 {
		n.stack[len(n.stack)-1] = screen
	} else {
		n.Push(screen)
	}
	screen.RegisterClickables(n.inputManager)
}

func (n *Navigator) CurrentScreen() Screen {
	if len(n.stack) == 0 {
		return nil
	}
	return n.stack[len(n.stack)-1]
}

func (n *Navigator) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if len(n.stack) > 1 {
			n.Pop()
			n.inputManager.ClearScreenClickable()
			if prevScreen := n.CurrentScreen(); prevScreen != nil {
				for _, btn := range prevScreen.(*GeneralScreen).screenButtons {
					n.inputManager.RegisterScreenClickable(btn)
				}
			}
		}
	}
	n.inputManager.Update()
}

// ***** GENERAL SCREEN *****
type GeneralScreen struct {
	color         color.RGBA
	label         string
	screenButtons []*Button
	targets       []*GeneralScreen
	offsetY       int
}

func NewGeneralScreen(screenWidth, screenHeight int, label string, offset int) *GeneralScreen {
	return &GeneralScreen{
		color:   color.RGBA{128, 128, 255, 255},
		label:   label,
		offsetY: offset,
	}
}

func (gs *GeneralScreen) GenerateButtons() {
	gs.screenButtons = nil
	for i, target := range gs.targets {
		buttonText := fmt.Sprintf("Go to Screen %d", i+1)

		button := NewButton(
			50,
			50+i*60,
			200,
			50,
			gs.offsetY,
			buttonText,
			color.RGBA{0, 128, 128, 255},
			color.RGBA{0, 255, 255, 255},
			color.RGBA{0, 100, 100, 255},
			target,
			nil,
		)
		gs.screenButtons = append(gs.screenButtons, button)

	}
}

func (gs *GeneralScreen) RegisterClickables(inputManager *InputManager) {
	for _, button := range gs.screenButtons {
		inputManager.RegisterScreenClickable(button)
	}
}

func (gs *GeneralScreen) Draw(screen *ebiten.Image) {
	screen.Fill(gs.color)
	ebitenutil.DebugPrint(screen, gs.label)
	for _, button := range gs.screenButtons {
		button.Draw(screen)
	}
}

func (gs *GeneralScreen) Update(navigator *Navigator) error {
	// ...
	return nil
}

// ***** TOP BAR *****
type TopBar struct {
	X, Y, Width, Height int
	Buttons             []*Button
	Color               color.RGBA
}

func NewTopBar(width, height, buttonCount int) *TopBar {
	topbar := &TopBar{
		X:      0,
		Y:      0,
		Height: height,
		Color:  color.RGBA{0, 0, 0, 255},
	}
	sideBarBtn := NewButton(
		100,
		0,
		100,
		height,
		0,
		"sideBar",
		color.RGBA{255, 255, 255, 255},
		color.RGBA{200, 200, 200, 255},
		color.RGBA{150, 150, 150, 255},
		nil,
		nil,
	)
	topbar.Buttons = append(topbar.Buttons, sideBarBtn)

	return topbar
}

func (tb *TopBar) RegisterClickables(inputManager *InputManager) {
	for _, button := range tb.Buttons {
		inputManager.Register(button)
	}
}

func (tb *TopBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(tb.X), float32(tb.Y), float32(tb.Width), float32(tb.Height), color.RGBA{30, 30, 30, 255}, true)
	for _, button := range tb.Buttons {
		button.Draw(screen)
	}
}

func (tb *TopBar) Update() {
	// ...
}

// ***** SIDEBAR CONTROLLER *****
type SidebarController struct {
	Sidebar       *Sidebar
	ClickableArea *ClickableArea
}

func NewSidebarController(width, height, topOffset, screenWidth, screenHeight int) *SidebarController {
	sidebar := NewSidebar(width, height, topOffset)

	clickableArea := NewClickableArea(
		sidebar.Width,
		sidebar.Y,
		screenWidth,
		sidebar.Height,
		color.RGBA{30, 30, 30, 150},
		sidebar.Visible,
		func() {
			// Placeholder for button-specific logic
		},
	)

	sidebarController := &SidebarController{
		Sidebar:       sidebar,
		ClickableArea: clickableArea,
	}

	sidebarController.ClickableArea.OnClickFunc = func() {
		log.Println("clickableArea clicked")
		if sidebar.Visible {
			sidebarController.ToggleSidebar()
		}
	}

	return sidebarController
}

func (sc *SidebarController) RegisterClickables(inputManager *InputManager) {
	inputManager.Register(sc.ClickableArea)
}

func (sc *SidebarController) ToggleSidebar() {
	if sc.Sidebar.TargetX == 0 {
		sc.Sidebar.TargetX = -sc.Sidebar.Width
		sc.ClickableArea.Active = false
	} else {
		sc.Sidebar.TargetX = 0
		sc.ClickableArea.Active = true
	}
}

func (sc *SidebarController) Update() {
	if sc.Sidebar.X < sc.Sidebar.TargetX {
		sc.Sidebar.X += sc.Sidebar.Speed
		if sc.Sidebar.X >= sc.Sidebar.TargetX {
			sc.Sidebar.X = sc.Sidebar.TargetX
			sc.Sidebar.Visible = sc.Sidebar.X == 0
		}
	} else if sc.Sidebar.X > sc.Sidebar.TargetX {
		sc.Sidebar.X -= sc.Sidebar.Speed
		if sc.Sidebar.X <= sc.Sidebar.TargetX {
			sc.Sidebar.X = sc.Sidebar.TargetX
			sc.Sidebar.Visible = sc.Sidebar.X == 0
		}
	}
}

func (sc *SidebarController) Draw(screen *ebiten.Image) {
	sc.Sidebar.Draw(screen)
	sc.ClickableArea.Draw(screen)
}

// ***** SIDEBAR *****
type Sidebar struct {
	X, Y, Width, Height int
	Speed               int
	Color               color.RGBA
	TargetX             int

	Visible bool
}

func NewSidebar(width, height, topOffset int) *Sidebar {
	return &Sidebar{
		X:       -width,
		Y:       topOffset,
		Width:   width,
		Height:  height,
		Visible: false,
		Speed:   10.0,
		TargetX: -width,
		Color:   color.RGBA{50, 50, 50, 255},
	}
}

func (s *Sidebar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(s.X),
		float32(s.Y),
		float32(s.Width),
		float32(s.Height),
		s.Color,
		true,
	)

}

// ***** CLICKABLE AREA *****
type ClickableArea struct {
	X, Y, Width, Height int
	Color               color.Color
	OnClickFunc         func()
	isHovered           bool
	isPressed           bool
	Active              bool
}

func NewClickableArea(x, y, width, height int, color color.Color, active bool, onClick func()) *ClickableArea {
	return &ClickableArea{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Color:       color,
		OnClickFunc: onClick,
		Active:      active,
	}
}

func (ca *ClickableArea) Contains(x, y int) bool {
	if !ca.Active {
		return false
	}
	return ca.Active && x >= ca.X && x <= ca.X+ca.Width &&
		y >= ca.Y && y <= ca.Y+ca.Height
}

func (ca *ClickableArea) OnClick() {
	if ca.Active && ca.OnClickFunc != nil {
		ca.OnClickFunc()
	}
	ca.isPressed = false
}

func (ca *ClickableArea) OnMouseDown() {
	ca.isPressed = true
}

func (ca *ClickableArea) SetHovered(isHovered bool) {
	ca.isHovered = isHovered
}
func (ca *ClickableArea) Draw(screen *ebiten.Image) {
	if !ca.Active {
		return
	}
	var drawColor color.Color
	if ca.isPressed {
		drawColor = color.RGBA{150, 0, 0, 255}
	} else if ca.isHovered {
		drawColor = color.RGBA{200, 0, 0, 255}
	} else {
		drawColor = ca.Color
	}
	vector.DrawFilledRect(
		screen,
		float32(ca.X),
		float32(ca.Y),
		float32(ca.Width),
		float32(ca.Height),
		drawColor,
		true,
	)
}

// ***** GAME *****
type Game struct {
	navigator         *Navigator
	screenWidth       int
	screenHeight      int
	topBar            *TopBar
	sidebarController *SidebarController
}

func NewGame(screenWidth, screenHeight int) *Game {
	var topBarHeight int = 50
	// Generate buttons for each screen
	inputManager := &InputManager{}
	navigator := NewNavigator(inputManager)

	sidebarController := NewSidebarController(200, screenHeight-topBarHeight, topBarHeight, screenWidth, screenHeight)
	sidebarController.RegisterClickables(inputManager)
	topBar := NewTopBar(screenWidth, topBarHeight, 3)

	topBar.Buttons[0].OnClickFunc = func() {
		log.Println("sidebar clicked")
		sidebarController.ToggleSidebar()
	}
	topBar.RegisterClickables(inputManager)

	// Create the example screens
	screenA := &GeneralScreen{
		color:   color.RGBA{255, 0, 0, 255},
		label:   "Screen A\nClick to navigate",
		targets: nil,
		offsetY: topBar.Height,
	}
	screenB := &GeneralScreen{
		color:   color.RGBA{0, 255, 0, 255},
		label:   "Screen B\nClick to navigate",
		targets: nil,
		offsetY: topBar.Height,
	}
	screenC := &GeneralScreen{
		color:   color.RGBA{0, 0, 255, 255},
		label:   "Screen C\nClick to navigate",
		targets: nil,
		offsetY: topBar.Height,
	}
	screenD := &GeneralScreen{
		color:   color.RGBA{255, 255, 0, 255},
		label:   "Screen D\nClick to navigate",
		targets: nil,
		offsetY: topBar.Height,
	}
	screenE := &GeneralScreen{
		color:   color.RGBA{0, 255, 255, 255},
		label:   "Screen E\nClick to navigate",
		targets: nil,
		offsetY: topBar.Height,
	}
	screenF := &GeneralScreen{
		color:   color.RGBA{255, 0, 255, 255},
		label:   "Screen F\nClick to navigate",
		targets: nil,
		offsetY: topBar.Height,
	}

	// Establish the links
	screenA.targets = []*GeneralScreen{screenB, screenC}
	screenB.targets = []*GeneralScreen{screenD, screenE}
	screenC.targets = []*GeneralScreen{screenF}

	for _, screen := range []*GeneralScreen{screenA, screenB, screenC, screenD, screenE, screenF} {
		screen.GenerateButtons()
		for _, button := range screen.screenButtons {
			button.OnClickFunc = func() {
				navigator.inputManager.ClearScreenClickable()
				navigator.Push(button.target)
			}
		}
	}

	navigator.Push(screenA)

	return &Game{
		navigator:         navigator,
		topBar:            topBar,
		sidebarController: sidebarController,
		screenWidth:       screenWidth,
		screenHeight:      screenHeight,
	}
}

func (g *Game) Update() error {
	g.sidebarController.Update()
	g.navigator.Update()
	currentScreen := g.navigator.CurrentScreen()
	if currentScreen != nil {
		if err := currentScreen.Update(g.navigator); err != nil {
			log.Fatalf("error updating screen: %v", err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the current screen in a new image
	ebitenImage := ebiten.NewImage(g.screenWidth, g.screenHeight-g.topBar.Height)

	var currentScreen Screen
	op := &ebiten.DrawImageOptions{}

	if g.navigator.CurrentScreen() != nil {
		currentScreen = g.navigator.CurrentScreen()
		currentScreen.Draw(ebitenImage)
		op.GeoM.Translate(0, float64(g.topBar.Height))
	}
	screen.DrawImage(ebitenImage, op)

	g.topBar.Draw(screen)
	g.sidebarController.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// ***** MAIN *****
func main() {
	game := NewGame(800, 600)

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Sidebar Navigation Example")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
