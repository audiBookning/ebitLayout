package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Content interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Navigator struct {
	stack      []Content
	transition float64
	animating  bool
	direction  int
}

func NewNavigator() *Navigator {
	return &Navigator{
		stack:      []Content{},
		transition: 1.0,
	}
}

func (n *Navigator) Push(content Content) {
	if len(n.stack) > 0 {
		n.animating = true
		n.transition = 0.0
		n.direction = 1
	}
	n.stack = append(n.stack, content)
}

func (n *Navigator) Pop() {
	if len(n.stack) > 1 && !n.animating {
		n.animating = true
		n.transition = 0.0
		n.direction = -1
	}
}

func (n *Navigator) Update() error {
	if len(n.stack) == 0 {
		return nil
	}

	if n.animating {
		n.transition += 0.05
		if n.transition >= 1.0 {
			n.transition = 1.0
			n.animating = false
			if n.direction == -1 {
				n.stack = n.stack[:len(n.stack)-1]
			}
		}
	}

	return n.stack[len(n.stack)-1].Update()
}

func (n *Navigator) Draw(screen *ebiten.Image) {
	if len(n.stack) == 0 {
		return
	}

	if n.animating {

		offsetX := 0.0
		if n.direction == 1 {
			offsetX = (1.0 - n.transition) * float64(screen.Bounds().Dx())
		} else if n.direction == -1 {
			offsetX = -n.transition * float64(screen.Bounds().Dx())
		}

		topContent := n.stack[len(n.stack)-1]
		tempScreen := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		topContent.Draw(tempScreen)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, 0)
		screen.DrawImage(tempScreen, op)

		if n.direction == -1 && len(n.stack) > 1 {
			previousContent := n.stack[len(n.stack)-2]
			tempScreen := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
			previousContent.Draw(tempScreen)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(screen.Bounds().Dx())+offsetX, 0)
			screen.DrawImage(tempScreen, op)
		}

	} else {

		n.stack[len(n.stack)-1].Draw(screen)
	}
}

type SimpleContent struct {
	color   color.Color
	message string
}

func NewSimpleContent(color color.Color, message string) *SimpleContent {
	return &SimpleContent{
		color:   color,
		message: message,
	}
}

func (c *SimpleContent) Update() error {
	return nil
}

func (c *SimpleContent) Draw(screen *ebiten.Image) {
	screen.Fill(c.color)
	ebitenutil.DebugPrint(screen, c.message)
}

type Game struct {
	navigator *Navigator
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.navigator.Push(NewSimpleContent(color.RGBA{0, 0, 255, 255}, "Blue Screen"))
	}
	if ebiten.IsKeyPressed(ebiten.KeyO) {
		g.navigator.Pop()
	}
	return g.navigator.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.navigator.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Navigator Example")

	navigator := NewNavigator()
	navigator.Push(NewSimpleContent(color.RGBA{255, 0, 0, 255}, "Red Screen"))

	game := &Game{navigator: navigator}

	go func() {
		<-time.After(2 * time.Second)
		navigator.Push(NewSimpleContent(color.RGBA{0, 255, 0, 255}, "Green Screen"))
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
