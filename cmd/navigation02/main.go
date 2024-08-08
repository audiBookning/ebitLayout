package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Screen interface {
	Update() error
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

type Button struct {
	x, y, width, height int
	text                string
	target              *GeneralScreen
}

func (b *Button) Draw(screen *ebiten.Image) {
	btnColor := color.RGBA{255, 255, 255, 255}
	btnRect := ebiten.NewImage(b.width, b.height)
	btnRect.Fill(btnColor)

	// Create a GeoM matrix and set translation
	var geoM ebiten.GeoM
	geoM.Translate(float64(b.x), float64(b.y))

	// Draw the button rectangle with the transformation
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM = geoM
	screen.DrawImage(btnRect, opts)

	// Draw the text on the button
	ebitenutil.DebugPrintAt(screen, b.text, b.x+10, b.y+10)
}

func (b *Button) IsClicked() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x >= b.x && x <= b.x+b.width && y >= b.y && y <= b.y+b.height
	}
	return false
}

type GeneralScreen struct {
	color   color.RGBA
	text    string
	targets []*GeneralScreen
	buttons []*Button
}

// GenerateButtons generates buttons based on the targets
func (s *GeneralScreen) GenerateButtons() {
	s.buttons = nil
	for i, target := range s.targets {
		buttonText := fmt.Sprintf("Go to Screen %d", i+1)
		s.buttons = append(s.buttons, &Button{
			x:      50,
			y:      50 + i*60,
			width:  200,
			height: 50,
			text:   buttonText,
			target: target,
		})
	}
}

func (s *GeneralScreen) Update() error {
	// Handle button clicks
	for _, button := range s.buttons {
		if button.IsClicked() {
			navigator.Push(button.target)
			break
		}
	}

	// Handle ESC key to go back
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if len(navigator.stack) > 1 {
			navigator.Pop()
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

var navigator *Navigator

func main() {
	navigator = NewNavigator()

	// Create the screens
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

	// Establish the links
	screenA.targets = []*GeneralScreen{screenB, screenC}
	screenB.targets = []*GeneralScreen{screenD, screenE}
	screenC.targets = []*GeneralScreen{screenF}

	// Generate buttons for each screen
	for _, screen := range []*GeneralScreen{screenA, screenB, screenC, screenD, screenE, screenF} {
		screen.GenerateButtons()
	}

	navigator.Push(screenA)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Navigation Example")

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	if currentScreen := navigator.CurrentScreen(); currentScreen != nil {
		return currentScreen.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if currentScreen := navigator.CurrentScreen(); currentScreen != nil {
		currentScreen.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}
