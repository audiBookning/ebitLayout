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

type InputManager struct {
	Clickables     []Clickable
	MouseX, MouseY int
}

func (im *InputManager) Register(c Clickable) {
	im.Clickables = append(im.Clickables, c)
}

func (im *InputManager) Clear() {
	im.Clickables = nil
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
	Text                string
	Color               color.Color
	HoverColor          color.Color
	ClickColor          color.Color
	isHovered           bool
	isPressed           bool
	OnClickFunc         func()
}

func NewButton(x, y, width, height int, text string, color, hoverColor, clickColor color.Color, onClick func()) *Button {
	return &Button{
		X: x, Y: y, Width: width, Height: height,
		Text:        text,
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
	ebitenutil.DebugPrintAt(screen, b.Text, b.X+10, b.Y+10)
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
	Update(*Navigator, *InputManager) error
	Draw(screen *ebiten.Image)
}

type Navigator struct {
	stack []Screen
}

func NewNavigator() *Navigator {
	return &Navigator{
		stack: make([]Screen, 0),
	}
}

func (n *Navigator) Push(screen Screen) {
	n.stack = append(n.stack, screen)
}

func (n *Navigator) Pop() {
	if len(n.stack) > 0 {
		n.stack = n.stack[:len(n.stack)-1]
	}
}

func (n *Navigator) Replace(screen Screen) {
	if len(n.stack) > 0 {
		n.stack[len(n.stack)-1] = screen
	} else {
		n.Push(screen)
	}
}

func (n *Navigator) CurrentScreen() Screen {
	if len(n.stack) == 0 {
		return nil
	}
	return n.stack[len(n.stack)-1]
}

type GeneralScreen struct {
	color   color.RGBA
	text    string
	targets []*GeneralScreen
	buttons []*Button
}

func (s *GeneralScreen) GenerateButtons(navigator *Navigator, inputManager *InputManager) {
	s.buttons = nil
	for i, target := range s.targets {
		buttonText := fmt.Sprintf("Go to Screen %d", i+1)
		button := NewButton(50, 50+i*60, 200, 50, buttonText,
			color.RGBA{255, 255, 255, 255},
			color.RGBA{200, 200, 200, 255},
			color.RGBA{150, 150, 150, 255},
			func(target *GeneralScreen) func() {
				return func() {
					navigator.Push(target)
					inputManager.Clear()
					for _, btn := range target.buttons {
						inputManager.Register(btn)
					}
				}
			}(target))
		s.buttons = append(s.buttons, button)
		inputManager.Register(button)
	}
}

func (s *GeneralScreen) Update(navigator *Navigator, inputManager *InputManager) error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if len(navigator.stack) > 1 {
			navigator.Pop()
			inputManager.Clear()
			if prevScreen := navigator.CurrentScreen(); prevScreen != nil {
				for _, btn := range prevScreen.(*GeneralScreen).buttons {
					inputManager.Register(btn)
				}
			}
		}
	}
	return nil
}

func (s *GeneralScreen) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)
	ebitenutil.DebugPrint(screen, s.text)
	for _, button := range s.buttons {
		button.Draw(screen)
	}
}

type Game struct {
	navigator    *Navigator
	inputManager *InputManager
}

func (g *Game) Update() error {
	if currentScreen := g.navigator.CurrentScreen(); currentScreen != nil {
		g.inputManager.Update()
		return currentScreen.Update(g.navigator, g.inputManager)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if currentScreen := g.navigator.CurrentScreen(); currentScreen != nil {
		currentScreen.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	navigator := NewNavigator()
	inputManager := &InputManager{}

	screenA := &GeneralScreen{
		color:   color.RGBA{255, 0, 0, 255},
		text:    "Screen A\nClick to navigate",
		targets: nil,
	}
	screenB := &GeneralScreen{
		color:   color.RGBA{0, 255, 0, 255},
		text:    "Screen B\nClick to navigate",
		targets: nil,
	}
	screenC := &GeneralScreen{
		color:   color.RGBA{0, 0, 255, 255},
		text:    "Screen C\nClick to navigate",
		targets: nil,
	}
	screenD := &GeneralScreen{
		color:   color.RGBA{255, 255, 0, 255},
		text:    "Screen D\nClick to navigate",
		targets: nil,
	}
	screenE := &GeneralScreen{
		color:   color.RGBA{0, 255, 255, 255},
		text:    "Screen E\nClick to navigate",
		targets: nil,
	}
	screenF := &GeneralScreen{
		color:   color.RGBA{255, 0, 255, 255},
		text:    "Screen F\nClick to navigate",
		targets: nil,
	}

	screenA.targets = []*GeneralScreen{screenB, screenC}
	screenB.targets = []*GeneralScreen{screenD, screenE}
	screenC.targets = []*GeneralScreen{screenF}

	for _, screen := range []*GeneralScreen{screenA, screenB, screenC, screenD, screenE, screenF} {
		screen.GenerateButtons(navigator, inputManager)
	}

	navigator.Push(screenA)
	inputManager.Clear()
	for _, button := range screenA.buttons {
		inputManager.Register(button)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Navigation Example with Input Manager")

	game := &Game{
		navigator:    navigator,
		inputManager: inputManager,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
