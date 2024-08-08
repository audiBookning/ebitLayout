package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ***** INPUT MANAGER *****
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

	// Update hover state
	for _, c := range im.Clickables {
		c.SetHovered(c.Contains(im.MouseX, im.MouseY))
	}
}

// ***** CLICKABLE INTERFACE *****
type Clickable interface {
	Contains(x, y int) bool
	OnClick()
	OnMouseDown()
	SetHovered(isHovered bool)
}

type Area struct {
	X, Y, Width, Height int
}

// ***** BUTTON *****
type Button struct {
	area        Area
	Label       string
	Color       color.Color
	HoverColor  color.Color
	ClickColor  color.Color
	isHovered   bool
	isPressed   bool
	OnClickFunc func() // Callback function
}

func NewButton(x, y, width, height int, label string, color, hoverColor, clickColor color.Color, onClick func()) *Button {
	return &Button{
		area: Area{
			X: x, Y: y, Width: width, Height: height,
		},
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
	vector.DrawFilledRect(screen, float32(b.area.X), float32(b.area.Y), float32(b.area.Width), float32(b.area.Height), drawColor, true)
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
	b.isPressed = false
}

func (b *Button) OnMouseDown() {
	b.isPressed = true
}

func (b *Button) SetHovered(isHovered bool) {
	b.isHovered = isHovered
}

// ***** GAME *****
type Game struct {
	InputManager *InputManager
	centerButton *Button
}

func NewGame() *Game {
	game := &Game{
		InputManager: &InputManager{},
	}

	// Create a button with normal, hover, and click colors
	centerButton := NewButton(200, 200, 80, 25, "CENTER",
		color.RGBA{200, 0, 0, 255},                      // Normal color
		color.RGBA{150, 0, 0, 255},                      // Hover color
		color.RGBA{100, 0, 0, 255},                      // Click color
		func() { log.Println("Center button clicked") }) // OnClick handler
	game.InputManager.Register(centerButton)
	game.centerButton = centerButton

	return game
}

func (g *Game) Update() error {
	g.InputManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.centerButton.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// ***** MAIN *****
func main() {
	game := NewGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Button with Hover and Click Effects")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
