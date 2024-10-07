package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Screen interface {
	Update(*Navigator) error
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
	shortcutKey         ebiten.Key
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

func (b *Button) IsKeyPressed() bool {
	return inpututil.IsKeyJustPressed(b.shortcutKey)
}

type ButtonConfig struct {
	x, y, width, height int
	text                string
	target              *GeneralScreen
	shortcutKey         ebiten.Key
}

type GeneralScreen struct {
	color         color.RGBA
	label         string
	screenButtons []*Button
}

func (gs *GeneralScreen) GenerateButtons(buttonConfigs []ButtonConfig) {
	gs.screenButtons = nil
	for _, config := range buttonConfigs {
		gs.screenButtons = append(gs.screenButtons, &Button{
			x:           config.x,
			y:           config.y,
			width:       config.width,
			height:      config.height,
			text:        config.text,
			target:      config.target,
			shortcutKey: config.shortcutKey,
		})
	}
}

func (s *GeneralScreen) Update(navigator *Navigator) error {
	// Handle button clicks and keyboard shortcuts
	for _, button := range s.screenButtons {
		if button.IsClicked() || button.IsKeyPressed() {
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
	ebitenutil.DebugPrint(screen, s.label)
	for _, button := range s.screenButtons {
		button.Draw(screen)
	}
}

func main() {
	navigator := NewNavigator()

	// Create the screens
	screenA := &GeneralScreen{
		color: color.RGBA{255, 0, 0, 255},
		label: "Screen A\nClick or press 1 or 2 to navigate",
	}
	screenB := &GeneralScreen{
		color: color.RGBA{0, 255, 0, 255},
		label: "Screen B\nClick or press 1 or 2 to navigate",
	}
	screenC := &GeneralScreen{
		color: color.RGBA{0, 0, 255, 255},
		label: "Screen C\nClick or press 1 to navigate",
	}
	screenD := &GeneralScreen{
		color: color.RGBA{255, 255, 0, 255},
		label: "Screen D",
	}
	screenE := &GeneralScreen{
		color: color.RGBA{0, 255, 255, 255},
		label: "Screen E",
	}
	screenF := &GeneralScreen{
		color: color.RGBA{255, 0, 255, 255},
		label: "Screen F",
	}

	// Generate buttons for each screen with shortcut keys
	screenA.GenerateButtons([]ButtonConfig{
		{50, 50, 200, 50, "Go to Screen B", screenB, ebiten.Key1},
		{50, 110, 200, 50, "Go to Screen C", screenC, ebiten.Key2},
	})
	screenB.GenerateButtons([]ButtonConfig{
		{50, 50, 200, 50, "Go to Screen D", screenD, ebiten.Key1},
		{50, 110, 200, 50, "Go to Screen E", screenE, ebiten.Key2},
	})
	screenC.GenerateButtons([]ButtonConfig{
		{50, 50, 200, 50, "Go to Screen F", screenF, ebiten.Key1},
	})

	navigator.Push(screenA)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Navigation Example")

	game := &Game{navigator: navigator}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	navigator *Navigator
}

func (g *Game) Update() error {
	if currentScreen := g.navigator.CurrentScreen(); currentScreen != nil {
		return currentScreen.Update(g.navigator)
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
