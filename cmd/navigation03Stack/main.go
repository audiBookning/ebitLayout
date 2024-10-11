package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth     = 640
	screenHeight    = 480
	leftColumnWidth = 160
)

type Content interface {
	Update() error
	Draw(screen *ebiten.Image)
	Width() int
	Height() int
}

type Navigator struct {
	stack          []Content
	transition     float64
	animating      bool
	direction      int
	contentOptions []ContentOptions
}

func NewNavigator() *Navigator {
	return &Navigator{
		stack:      []Content{},
		transition: 1.0,
	}
}

func (n *Navigator) SetContentOptions(options []ContentOptions) {
	n.contentOptions = options
}

func (n *Navigator) Push() {
	if len(n.contentOptions) == 0 {
		log.Println("No content options available to push.")
		return
	}

	if len(n.stack) > 0 {
		n.animating = true
		n.transition = 0.0
		n.direction = 1
	}

	index := len(n.stack) % len(n.contentOptions)
	opt := n.contentOptions[index]
	content := NewSimpleContent(opt.color, opt.message, opt.width, opt.height)
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
		tempScreen := ebiten.NewImage(topContent.Width(), topContent.Height())
		topContent.Draw(tempScreen)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, 0)
		screen.DrawImage(tempScreen, op)

		if n.direction == -1 && len(n.stack) > 1 {
			previousContent := n.stack[len(n.stack)-2]
			tempScreen := ebiten.NewImage(previousContent.Width(), previousContent.Height())
			previousContent.Draw(tempScreen)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(screen.Bounds().Dx())+offsetX, 0)
			screen.DrawImage(tempScreen, op)
		}

	} else {

		topContent := n.stack[len(n.stack)-1]
		tempScreen := ebiten.NewImage(topContent.Width(), topContent.Height())
		topContent.Draw(tempScreen)
		screen.DrawImage(tempScreen, nil)
	}
}

type SimpleContent struct {
	color   color.Color
	message string
	width   int
	height  int
}

func NewSimpleContent(color color.Color, message string, width, height int) *SimpleContent {
	return &SimpleContent{
		color:   color,
		message: message,
		width:   width,
		height:  height,
	}
}

func (c *SimpleContent) Update() error {
	return nil
}

func (c *SimpleContent) Draw(screen *ebiten.Image) {
	screen.Fill(c.color)
	height := c.height
	if height > screenHeight {
		height = screenHeight
	}
	ebitenutil.DebugPrintAt(screen, c.message, 10, height-20)
}

func (c *SimpleContent) Width() int {
	return c.width
}

func (c *SimpleContent) Height() int {
	return c.height
}

type Game struct {
	navigator     *Navigator
	lastKeyState  map[ebiten.Key]bool
	leftColumnMsg string
}

func NewGame(navigator *Navigator) *Game {
	return &Game{
		navigator:     navigator,
		lastKeyState:  make(map[ebiten.Key]bool),
		leftColumnMsg: "Static Left Column",
	}
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && !g.lastKeyState[ebiten.KeyArrowRight] {
		g.navigator.Push()
		g.lastKeyState[ebiten.KeyArrowRight] = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && !g.lastKeyState[ebiten.KeyArrowLeft] {
		g.navigator.Pop()
		g.lastKeyState[ebiten.KeyArrowLeft] = true
	}

	if !ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.lastKeyState[ebiten.KeyArrowRight] = false
	}
	if !ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.lastKeyState[ebiten.KeyArrowLeft] = false
	}

	return g.navigator.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {

	leftColumn := ebiten.NewImage(leftColumnWidth, screenHeight)
	leftColumn.Fill(color.RGBA{50, 50, 50, 255})
	ebitenutil.DebugPrintAt(leftColumn, g.leftColumnMsg, 10, 10)
	screen.DrawImage(leftColumn, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(leftColumnWidth, 0)

	navigatorAreaWidth := screenWidth - leftColumnWidth
	navigatorArea := ebiten.NewImage(navigatorAreaWidth, screenHeight)
	g.navigator.Draw(navigatorArea)

	screen.DrawImage(navigatorArea, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type ContentOptions struct {
	color   color.Color
	message string
	width   int
	height  int
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Navigator Example with Static Left Column")

	navigator := NewNavigator()

	contentOptions := []ContentOptions{
		{color.RGBA{255, 0, 0, 255}, "Red Screen", 480, 480},
		{color.RGBA{0, 0, 255, 255}, "Blue Screen", 480, 480},
		{color.RGBA{0, 255, 0, 255}, "Green Screen", 320, 240},
		{color.RGBA{255, 255, 0, 255}, "Yellow Screen", 480, 600},
		{color.RGBA{255, 0, 255, 255}, "Magenta Screen", 400, 300},
	}

	navigator.SetContentOptions(contentOptions)

	navigator.Push()

	game := NewGame(navigator)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
