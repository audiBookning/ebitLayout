package main

import (
	"image/color"
	"log"

	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
)

// UIElement defines an interface for drawable UI elements.
type UIElement interface {
	SetPosition(x, y float64)
	Draw(screen *ebiten.Image)
	GetHeight() float64
}

// RedRectangle represents a red rectangle UI element.
type RedRectangle struct {
	Width  int
	Height int
	x, y   float64
}

func (r *RedRectangle) SetPosition(x, y float64) {
	r.x = x
	r.y = y
}

func (r *RedRectangle) Draw(screen *ebiten.Image) {
	uiElement := ebiten.NewImage(r.Width, r.Height)
	uiElement.Fill(color.RGBA{255, 0, 0, 255}) // Red color

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.x, r.y)
	screen.DrawImage(uiElement, op)
}

func (r *RedRectangle) GetHeight() float64 {
	return float64(r.Height)
}

// BlueCircle represents a blue circle UI element.
type BlueCircle struct {
	Radius int
	x, y   float64
}

func (b *BlueCircle) SetPosition(x, y float64) {
	b.x = x
	b.y = y
}

func (b *BlueCircle) Draw(screen *ebiten.Image) {
	diameter := b.Radius * 2
	circle := ebiten.NewImage(diameter, diameter)
	circle.Fill(color.Transparent)

	for x := 0; x < diameter; x++ {
		for y := 0; y < diameter; y++ {
			dx := x - b.Radius
			dy := y - b.Radius
			if dx*dx+dy*dy <= b.Radius*b.Radius {
				circle.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue color
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(circle, op)
}

func (b *BlueCircle) GetHeight() float64 {
	return float64(b.Radius * 2)
}

type Game struct {
	layoutSystem *layout.BreakpointLayoutSystem
	currentBP    layout.Breakpoint
	previousBP   layout.Breakpoint
	elements     []UIElement
}

func NewGame() *Game {
	layoutSystem := layout.NewLayoutSystem(nil)

	// Define UI elements
	redRect := &RedRectangle{Width: 100, Height: 50}
	blueCircle := &BlueCircle{Radius: 40}

	// Add elements to the game
	game := &Game{
		layoutSystem: layoutSystem,
		elements:     []UIElement{redRect, blueCircle},
		previousBP:   -1, // Initialize with an invalid breakpoint
	}

	return game
}

func (g *Game) Update() error {
	if g.currentBP != g.previousBP {
		g.layoutElements()
		g.previousBP = g.currentBP
	}
	return nil
}

// layoutElements calculates and sets the positions of UI elements based on the current layout.
func (g *Game) layoutElements() {
	layout := g.layoutSystem.GetLayout(g.currentBP)
	numColumns := layout.Columns
	columnWidth := float64(layout.Width) / float64(numColumns)
	padding := 20.0

	// Initialize a slice to track the current Y position for each column
	columnY := make([]float64, numColumns)
	for i := 0; i < numColumns; i++ {
		columnY[i] = padding
	}

	for idx, element := range g.elements {
		// Determine column index based on element index
		columnIndex := idx % numColumns

		// X position is based on the column index
		x := float64(columnIndex)*columnWidth + padding

		// Y position is based on the current Y position of the column
		y := columnY[columnIndex]

		// Set the element's position
		element.SetPosition(x, y)

		// Update the current Y position for the column
		columnY[columnIndex] += element.GetHeight() + padding
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	layout := g.layoutSystem.GetLayout(g.currentBP)

	// Draw the layout background and column guides
	layout.DrawLayout(screen)

	// Render each UI element at its position
	for _, element := range g.elements {
		element.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Determine the current breakpoint based on window width
	g.currentBP = g.layoutSystem.DetermineBreakpoint(outsideWidth, outsideHeight)
	currentLayout := g.layoutSystem.GetLayout(g.currentBP)
	return currentLayout.Width, currentLayout.Height
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Ebitengine Responsive Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
