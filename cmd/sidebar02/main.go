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

type InputManager struct {
	Clickables     []Clickable
	MouseX, MouseY int
}

func (im *InputManager) Register(c Clickable) {
	im.Clickables = append(im.Clickables, c)
}

func (im *InputManager) Update() {
	im.MouseX, im.MouseY = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, c := range im.Clickables {
			if c.Contains(im.MouseX, im.MouseY) {
				c.OnMouseDown()
				break
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, c := range im.Clickables {
			if c.Contains(im.MouseX, im.MouseY) {
				c.OnClick()
				break
			}
		}
	}

	for _, c := range im.Clickables {
		c.SetHovered(c.Contains(im.MouseX, im.MouseY))
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

type TopBar struct {
	X, Y, Width, Height int
	Buttons             []*Button
}

func NewTopBar(width, height int, numButtons int, inputManager *InputManager) *TopBar {
	buttons := make([]*Button, numButtons)
	for i := range buttons {
		buttons[i] = NewButton(
			i*100+10, 10, 80, 25,
			"Menu "+strconv.Itoa(i),
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

	inputManager.Register(sidebarControler.ClickableArea)

	return sidebarControler
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

type ClickableArea struct {
	X, Y, Width, Height int
	OnClickFunc         func()
	Active              bool
}

func (ca *ClickableArea) Contains(x, y int) bool {
	if !ca.Active {
		return false
	}

	return x > ca.X &&
		y > ca.Y &&
		x <= ca.Width &&
		y <= ca.Height
}

func (ca *ClickableArea) OnClick() {
	if ca.Active && ca.OnClickFunc != nil {
		ca.OnClickFunc()
	}
}

func (ca *ClickableArea) OnMouseDown() {

}

func (ca *ClickableArea) SetHovered(isHovered bool) {

}

func (ca *ClickableArea) Draw(screen *ebiten.Image) {
	if !ca.Active {
		return
	}
	vector.DrawFilledRect(screen, float32(ca.X), float32(ca.Y), float32(ca.Width), float32(ca.Height), color.RGBA{30, 30, 30, 150}, true)
}

type Sidebar struct {
	X, Y, Width, Height int
	Visible             bool
	TargetX             int
	Speed               int
}

func NewSidebar(width, height, topOffset int) *Sidebar {
	return &Sidebar{
		X:       -width,
		Y:       topOffset,
		Width:   width,
		Height:  height,
		Visible: false,
		Speed:   10,
		TargetX: -width,
	}
}

func (s *Sidebar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(s.X), float32(s.Y), float32(s.Width), float32(s.Height), color.RGBA{50, 50, 50, 255}, true)
}

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

	inputManager := &InputManager{}

	numButtons := 2
	topbar := NewTopBar(screenWidth, topBarHeight, numButtons, inputManager)

	sidebarController := NewSidebarController(200, screenHeight-topBarHeight, topBarHeight, screenWidth, screenHeight, inputManager)

	game := &Game{
		TopBar:            topbar,
		InputManager:      inputManager,
		SidebarController: sidebarController,
		screenWidth:       screenWidth,
		screenHeight:      screenHeight,
	}

	topbar.Buttons[0].OnClickFunc = func() {
		log.Println("Button 1 clicked")
		game.SidebarController.ToggleSidebar()
	}

	topbar.Buttons[1].OnClickFunc = func() {
		log.Println("Button 2 clicked")
		game.SidebarController.ToggleSidebar()
	}

	centerButton := NewButton(
		200, 200, 80, 25,
		"CENTER",
		color.RGBA{200, 0, 0, 255},
		color.RGBA{150, 0, 0, 255},
		color.RGBA{100, 0, 0, 255},
		func() {
			log.Println("Center button clicked")
		},
	)
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

func main() {
	game := NewGame()
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Sliding Sidebar Menu with Top Bar")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
