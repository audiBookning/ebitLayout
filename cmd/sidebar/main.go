package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ***** INPUT MANAGER *****
type InputManager struct {
	Clickables []Clickable
}

func (im *InputManager) Register(c Clickable) {
	im.Clickables = append(im.Clickables, c)
}

func (im *InputManager) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, c := range im.Clickables {
			if c.Contains(mx, my) {
				c.OnClick()
				break
			}
		}
	}
}

// ***** CLICKABLE INTERFACE *****
type Clickable interface {
	Contains(x, y int) bool
	OnClick()
}

// ***** BUTTON *****
type Button struct {
	area        Area
	Color       color.Color
	Label       string
	OnClickFunc func() // Callback function
}

func NewButton(x, y, width, height int, label string, color color.Color, onClick func()) *Button {
	return &Button{
		area: Area{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		Label:       label,
		Color:       color,
		OnClickFunc: onClick,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(b.area.X), float32(b.area.Y), float32(b.area.Width), float32(b.area.Height), b.Color, true)
	ebitenutil.DebugPrintAt(screen, b.Label, b.area.X+10, b.area.Y+10)
}

func (b *Button) Contains(x, y int) bool {
	return x >= b.area.X && x <= b.area.X+b.area.Width &&
		y >= b.area.Y && y <= b.area.Y+b.area.Height
}

func (b *Button) OnClick() {
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
}

// ***** TOP BAR *****
type TopBar struct {
	area    Area
	Buttons []*Button
}

func NewTopBar(width, height int, numButtons int) *TopBar {
	buttons := make([]*Button, numButtons)
	for i := range buttons {
		buttons[i] = &Button{
			area: Area{
				X:      i*100 + 10,
				Y:      10,
				Width:  80,
				Height: 25,
			},
			Label: "Menu " + strconv.Itoa(i),
			Color: color.RGBA{200, 0, 0, 255},
			OnClickFunc: func() {
				// Placeholder for button-specific logic
			},
		}
	}
	return &TopBar{
		area: Area{
			X:      0,
			Y:      0,
			Width:  width,
			Height: height,
		},
		Buttons: buttons,
	}
}

func (tb *TopBar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(tb.area.X), float32(tb.area.Y), float32(tb.area.Width), float32(tb.area.Height), color.RGBA{30, 30, 30, 255}, true)
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

func NewSidebarController(width, height, topOffset int, screenWidth, screenHeight int) *SidebarController {
	sidebar := NewSidebar(width, height, topOffset)
	clickableArea := &ClickableArea{
		area: Area{
			X:      sidebar.area.Width,
			Width:  screenWidth,
			Y:      sidebar.area.Y,
			Height: sidebar.area.Height,
		},
		OnClickFunc: func() {
			// Placeholder for button-specific logic
		},
		Active: sidebar.Visible,
	}
	sidebarControler := &SidebarController{
		Sidebar:       sidebar,
		ClickableArea: clickableArea,
	}

	sidebarControler.ClickableArea.OnClickFunc = func() {
		log.Println("clickableArea clicked")
		if sidebar.Visible {
			sidebarControler.ToggleSidebar()
		}
	}

	return sidebarControler
}

func (sc *SidebarController) ToggleSidebar() {
	if sc.Sidebar.TargetX == 0 {
		sc.Sidebar.TargetX = -sc.Sidebar.area.Width
		sc.ClickableArea.Active = false
	} else {
		sc.Sidebar.TargetX = 0
		sc.ClickableArea.Active = true
	}
}

func (sc *SidebarController) Update() {
	if sc.Sidebar.area.X < sc.Sidebar.TargetX {
		sc.Sidebar.area.X += sc.Sidebar.Speed
		if sc.Sidebar.area.X >= sc.Sidebar.TargetX {
			sc.Sidebar.area.X = sc.Sidebar.TargetX
			sc.Sidebar.Visible = sc.Sidebar.area.X == 0
		}
	} else if sc.Sidebar.area.X > sc.Sidebar.TargetX {
		sc.Sidebar.area.X -= sc.Sidebar.Speed
		if sc.Sidebar.area.X <= sc.Sidebar.TargetX {
			sc.Sidebar.area.X = sc.Sidebar.TargetX
			sc.Sidebar.Visible = sc.Sidebar.area.X == 0
		}
	}
}

func (sc *SidebarController) Draw(screen *ebiten.Image) {
	sc.Sidebar.Draw(screen)
	sc.ClickableArea.Draw(screen)
}

type Area struct {
	X, Y, Width, Height int
}

// ***** CLICKABLE AREA (OUTSIDE SIDEBAR) *****
type ClickableArea struct {
	area        Area
	OnClickFunc func()
	Active      bool
}

func (ca *ClickableArea) Contains(x, y int) bool {
	if !ca.Active {
		return false
	}
	// Check if the click is within the main area excluding the sidebar and top bar
	return x > ca.area.X &&
		y > ca.area.Y &&
		x <= ca.area.Width &&
		y <= ca.area.Height
}

func (ca *ClickableArea) OnClick() {
	if ca.Active && ca.OnClickFunc != nil {
		ca.OnClickFunc()
	}
}

func (ca *ClickableArea) Draw(screen *ebiten.Image) {
	if !ca.Active {
		return
	}
	// todo: add a subtle easing animation to the alpha color of the rect?
	// can be added to the SidebarController update method
	vector.DrawFilledRect(screen, float32(ca.area.X), float32(ca.area.Y), float32(ca.area.Width), float32(ca.area.Height), color.RGBA{30, 30, 30, 150}, true)
}

// ***** SIDEBAR *****
type Sidebar struct {
	area    Area
	Visible bool
	TargetX int
	Speed   int
}

func NewSidebar(width, height, topOffset int) *Sidebar {
	return &Sidebar{
		area: Area{
			X:      -width,
			Y:      topOffset,
			Width:  width,
			Height: height,
		},
		Visible: false,
		Speed:   10.0,
		TargetX: -width,
	}
}

func (s *Sidebar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(s.area.X), float32(s.area.Y), float32(s.area.Width), float32(s.area.Height), color.RGBA{50, 50, 50, 255}, true)
}

// ***** GAME *****
type Game struct {
	TopBar                    *TopBar
	centerButton              *Button
	InputManager              *InputManager
	SidebarController         *SidebarController
	screenWidth, screenHeight int
}

func NewGame() *Game {
	var topBarHeight int = 50
	screenWidth, screenHeight := 640, 480

	// Define the number of buttons in the topbar
	numButtons := 2
	topbar := NewTopBar(screenWidth, topBarHeight, numButtons)

	sidebarController := NewSidebarController(200, screenHeight-topBarHeight, topBarHeight, screenWidth, screenHeight)

	game := &Game{
		TopBar:            topbar,
		InputManager:      &InputManager{},
		SidebarController: sidebarController,
		screenWidth:       screenWidth,
		screenHeight:      screenHeight,
	}

	// Assign specific functions to button click handlers
	topbar.Buttons[0].OnClickFunc = func() {
		log.Println("Button 1 clicked")
		game.SidebarController.ToggleSidebar()
	}

	topbar.Buttons[1].OnClickFunc = func() {
		log.Println("Button 2 clicked")
		game.SidebarController.ToggleSidebar()
	}

	// Register all buttons with the input manager
	for _, button := range topbar.Buttons {
		game.InputManager.Register(button)
	}

	// Register the clickable area from the SidebarController with the input manager
	game.InputManager.Register(sidebarController.ClickableArea)

	// Add another centerButton in the center of the screen for testing purposes
	centerButton := &Button{
		area: Area{
			X:      200,
			Y:      200,
			Width:  80,
			Height: 25,
		},
		Label: "CENTER",
		Color: color.RGBA{200, 0, 0, 255},
		OnClickFunc: func() {
			log.Println("Center button clicked")
		},
	}
	game.InputManager.Register(centerButton)
	game.centerButton = centerButton

	return game
}

func (g *Game) Update() error {
	g.TopBar.Update()
	g.SidebarController.Update()
	g.InputManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.centerButton.Draw(screen)
	g.TopBar.Draw(screen)
	g.SidebarController.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// ***** MAIN *****
func main() {
	game := NewGame()
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Sliding Sidebar Menu with Top Bar")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
