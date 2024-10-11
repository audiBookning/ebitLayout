package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type InputManager struct {
	Clickables []Clickable
}

func (im *InputManager) Register(c Clickable) {
	im.Clickables = append(im.Clickables, c)
}

func (im *InputManager) Clear() {
	im.Clickables = nil
}

func (im *InputManager) Update() {
	mouseX, mouseY := ebiten.CursorPosition()

	for _, c := range im.Clickables {
		if area, ok := c.(*ClickableArea); ok && area.Active {
			if area.Contains(mouseX, mouseY) {
				if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
					area.OnMouseDown()
				}
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
					area.OnClick()
				}

				return
			}
		}
	}

	for _, c := range im.Clickables {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if c.Contains(mouseX, mouseY) {
				c.OnMouseDown()
				break
			}
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			if c.Contains(mouseX, mouseY) {
				c.OnClick()
				break
			}
		}

		c.SetHovered(c.Contains(mouseX, mouseY))
	}
}

type Clickable interface {
	Contains(x, y int) bool
	OnClick()
	OnMouseDown()
	SetHovered(isHovered bool)
}

type Button struct {
	X, Y, Width, Height int
	Label               string
	Color               color.Color
	HoverColor          color.Color
	ClickColor          color.Color
	isHovered           bool
	isPressed           bool
	OnClickFunc         func()
}

func NewButton(x, y, width, height int, label string, color, hoverColor, clickColor color.Color, onClick func()) *Button {
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
	return x >= b.X && x <= b.X+b.Width &&
		y >= b.Y && y <= b.Y+b.Height
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

type Screen interface {
	Update() error
	Draw(screen *ebiten.Image)
	RegisterClickables(inputManager *InputManager)
}

type Navigator struct {
	stack []Screen
}

func NewNavigator() *Navigator {
	return &Navigator{
		stack: make([]Screen, 0),
	}
}

func (n *Navigator) Push(screen Screen, inputManager *InputManager) {
	n.stack = append(n.stack, screen)
	screen.RegisterClickables(inputManager)
}

func (n *Navigator) Pop(inputManager *InputManager) {
	if len(n.stack) > 0 {
		n.stack = n.stack[:len(n.stack)-1]
	}
	if len(n.stack) > 0 {
		n.stack[len(n.stack)-1].RegisterClickables(inputManager)
	}
}

func (n *Navigator) Replace(screen Screen, inputManager *InputManager) {
	if len(n.stack) > 0 {
		n.stack[len(n.stack)-1] = screen
	} else {
		n.Push(screen, inputManager)
	}
	screen.RegisterClickables(inputManager)
}

func (n *Navigator) CurrentScreen() Screen {
	if len(n.stack) == 0 {
		return nil
	}
	return n.stack[len(n.stack)-1]
}

type GeneralScreen struct {
	Color   color.RGBA
	Label   string
	Buttons []*Button
}

func (gs *GeneralScreen) AddButton(x, y, width, height int, label string, color, hoverColor, clickColor color.Color, onClick func()) {
	button := NewButton(x, y, width, height, label, color, hoverColor, clickColor, onClick)
	gs.Buttons = append(gs.Buttons, button)
}

func (gs *GeneralScreen) RegisterClickables(inputManager *InputManager) {
	for _, button := range gs.Buttons {
		inputManager.Register(button)
	}
}

func (gs *GeneralScreen) Draw(screen *ebiten.Image) {
	screen.Fill(gs.Color)
	for _, button := range gs.Buttons {
		button.Draw(screen)
	}
}

type TopBar struct {
	Buttons             []*Button
	X, Y, Width, Height int
	Color               color.RGBA
}

func NewTopBar(width, height int, numButtons int, inputManager *InputManager) *TopBar {
	buttons := make([]*Button, numButtons)
	for i := range buttons {
		buttons[i] = NewButton(
			i*100+10, 10, 80, 25,
			"sidebar",
			color.RGBA{200, 0, 0, 255},
			color.RGBA{150, 0, 0, 255},
			color.RGBA{100, 0, 0, 255},
			nil,
		)
		inputManager.Register(buttons[i])
	}
	return &TopBar{
		X:       0,
		Y:       0,
		Width:   width,
		Height:  height,
		Buttons: buttons,
		Color:   color.RGBA{0, 0, 0, 255},
	}
}

func (tb *TopBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(tb.X), float32(tb.Y), float32(tb.Width), float32(tb.Height), color.RGBA{30, 30, 30, 255}, true)
	for _, button := range tb.Buttons {
		button.Draw(screen)
	}
}

func (tb *TopBar) Update() {

}

type SidebarController struct {
	Sidebar       *Sidebar
	ClickableArea *ClickableArea
}

func NewSidebarController(width, height, topOffset int, screenWidth, screenHeight int, inputManager *InputManager) *SidebarController {
	sidebar := NewSidebar(width, height, topOffset)
	clickableArea := &ClickableArea{

		X:      sidebar.Width,
		Width:  screenWidth,
		Y:      sidebar.Y,
		Height: sidebar.Height,
		OnClickFunc: func() {

		},
		Active: sidebar.Visible,
	}

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

	inputManager.Register(sidebarController.ClickableArea)

	return sidebarController
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

type Sidebar struct {
	X, Y, Width, Height int
	Speed               int
	TargetX             int
	Visible             bool
}

func NewSidebar(width, height, topOffset int) *Sidebar {
	return &Sidebar{
		X:       0,
		Y:       topOffset,
		Width:   width,
		Height:  height,
		Visible: false,
		Speed:   4,
		TargetX: -width,
	}
}

func (s *Sidebar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(s.X),
		float32(s.Y),
		float32(s.Width),
		float32(s.Height),
		color.RGBA{0, 0, 128, 255},
		true,
	)
}

type ClickableArea struct {
	X, Y, Width, Height int
	OnClickFunc         func()
	Active              bool
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
}

func (ca *ClickableArea) OnMouseDown() {}

func (ca *ClickableArea) SetHovered(isHovered bool) {}

func (ca *ClickableArea) Draw(screen *ebiten.Image) {
	if !ca.Active {
		return
	}

	vector.DrawFilledRect(
		screen,
		float32(ca.X),
		float32(ca.Y),
		float32(ca.Width),
		float32(ca.Height),
		color.RGBA{30, 30, 30, 150},
		true,
	)
}

type Game struct {
	navigator    *Navigator
	inputManager *InputManager
	screenWidth  int
	screenHeight int
}

func NewGame() *Game {
	screenWidth, screenHeight := 640, 480
	inputManager := &InputManager{}
	navigator := NewNavigator()

	mainScreen := NewMainScreen(screenWidth, screenHeight, inputManager)
	navigator.Push(mainScreen, inputManager)

	return &Game{
		navigator:    navigator,
		inputManager: inputManager,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

func (g *Game) Update() error {
	if currentScreen := g.navigator.CurrentScreen(); currentScreen != nil {
		g.inputManager.Update()
		return currentScreen.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if currentScreen := g.navigator.CurrentScreen(); currentScreen != nil {
		currentScreen.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Sliding Sidebar Menu with Top Bar and Navigation")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type MainScreen struct {
	TopBar            *TopBar
	SidebarController *SidebarController
	centerButton      *Button
}

func NewMainScreen(screenWidth, screenHeight int, inputManager *InputManager) *MainScreen {
	topbar := NewTopBar(screenWidth, 50, 1, inputManager)
	sidebarController := NewSidebarController(200, screenHeight-50, 50, screenWidth, screenHeight, inputManager)

	topbar.Buttons[0].OnClickFunc = func() {
		sidebarController.ToggleSidebar()
	}

	centerButton := NewButton(
		screenWidth/2-50,
		screenHeight/2-25,
		100,
		50,
		"Go to Next",
		color.RGBA{0, 128, 0, 255},
		color.RGBA{0, 255, 0, 255},
		color.RGBA{0, 100, 0, 255},
		func() {
			log.Println("Center button clicked")
		},
	)

	screen := &MainScreen{
		TopBar:            topbar,
		SidebarController: sidebarController,
		centerButton:      centerButton,
	}

	screen.RegisterClickables(inputManager)

	return screen
}

func (s *MainScreen) Update() error {
	s.TopBar.Update()
	s.SidebarController.Update()
	return nil
}

func (s *MainScreen) Draw(screen *ebiten.Image) {
	s.centerButton.Draw(screen)
	s.TopBar.Draw(screen)
	s.SidebarController.Draw(screen)
}

func (s *MainScreen) RegisterClickables(inputManager *InputManager) {
	inputManager.Register(s.centerButton)
	for _, button := range s.TopBar.Buttons {
		inputManager.Register(button)
	}
	inputManager.Register(s.SidebarController.ClickableArea)
}

type SecondaryScreen struct {
	GeneralScreen
}

func NewSecondaryScreen(screenWidth, screenHeight int) *SecondaryScreen {
	screen := &SecondaryScreen{
		GeneralScreen: GeneralScreen{
			Color: color.RGBA{0, 0, 128, 255},
			Label: "Secondary Screen",
		},
	}

	screen.AddButton(50, 100, 150, 40, "Back", color.RGBA{128, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{100, 0, 0, 255}, nil)

	return screen
}
