package main

import (
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

type Arrow struct {
	direction ebiten.Key
	target    *GeneralScreen
}

type GeneralScreen struct {
	color  color.RGBA
	text   string
	arrows []*Arrow
}

func (s *GeneralScreen) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if len(navigator.stack) > 1 {
			navigator.Pop()
		}
		return nil
	}

	for _, arrow := range s.arrows {
		if inpututil.IsKeyJustPressed(arrow.direction) {
			navigator.Push(arrow.target)
		}
	}
	return nil
}

func (s *GeneralScreen) Draw(screen *ebiten.Image) {
	screen.Fill(s.color)
	ebitenutil.DebugPrint(screen, s.text)
}

var navigator *Navigator

func main() {
	navigator = NewNavigator()

	screenA := &GeneralScreen{
		color: color.RGBA{255, 0, 0, 255},
		text:  "Screen A\nPress B to go to Screen B\nPress C to go to Screen C",
	}
	screenB := &GeneralScreen{
		color: color.RGBA{0, 255, 0, 255},
		text:  "Screen B\nPress D to go to Screen D\nPress E to go to Screen E\nPress Esc to go back",
	}
	screenC := &GeneralScreen{
		color: color.RGBA{0, 0, 255, 255},
		text:  "Screen C\nPress F to go to Screen F\nPress Esc to go back",
	}
	screenD := &GeneralScreen{
		color: color.RGBA{255, 255, 0, 255},
		text:  "Screen D\nPress Esc to go back",
	}
	screenE := &GeneralScreen{
		color: color.RGBA{0, 255, 255, 255},
		text:  "Screen E\nPress Esc to go back",
	}
	screenF := &GeneralScreen{
		color: color.RGBA{255, 0, 255, 255},
		text:  "Screen F\nPress Esc to go back",
	}

	screenA.arrows = []*Arrow{
		{direction: ebiten.KeyB, target: screenB},
		{direction: ebiten.KeyC, target: screenC},
	}
	screenB.arrows = []*Arrow{
		{direction: ebiten.KeyD, target: screenD},
		{direction: ebiten.KeyE, target: screenE},
	}
	screenC.arrows = []*Arrow{
		{direction: ebiten.KeyF, target: screenF},
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
