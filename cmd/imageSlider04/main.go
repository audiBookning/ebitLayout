package main

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Content interface defines the methods that each content type must implement.
type Content interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// Navigator struct manages the stack of content and animations between them.
type Navigator struct {
	stack      []Content // Stack of content objects
	transition float64   // Animation progress (0 to 1)
	animating  bool      // Whether an animation is ongoing
	direction  int       // 1 for push, -1 for pop
}

// NewNavigator creates a new Navigator instance.
func NewNavigator() *Navigator {
	return &Navigator{
		stack:      []Content{},
		transition: 1.0,
	}
}

// Push adds a new content to the stack and starts the push animation.
func (n *Navigator) Push(content Content) {
	if len(n.stack) > 0 {
		n.animating = true
		n.transition = 0.0
		n.direction = 1
	}
	n.stack = append(n.stack, content)
}

// Pop removes the top content from the stack if an animation is not in progress.
func (n *Navigator) Pop() {
	if len(n.stack) > 1 && !n.animating {
		n.animating = true
		n.transition = 0.0
		n.direction = -1
	}
}

// Update updates the navigator's state, handling animations and updating the top content.
func (n *Navigator) Update() error {
	if len(n.stack) == 0 {
		return nil
	}

	// Handle animations
	if n.animating {
		n.transition += 0.05 // Adjust this value for speed
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

// Draw renders the current content and manages the animation between transitions.
func (n *Navigator) Draw(screen *ebiten.Image) {
	if len(n.stack) == 0 {
		return
	}

	if n.animating {
		// Calculate offset for sliding animation
		offsetX := 0.0
		if n.direction == 1 {
			offsetX = (1.0 - n.transition) * float64(screen.Bounds().Dx())
		} else if n.direction == -1 {
			offsetX = -n.transition * float64(screen.Bounds().Dx())
		}

		// Draw the current top Content
		topContent := n.stack[len(n.stack)-1]
		tempScreen := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		topContent.Draw(tempScreen)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, 0)
		screen.DrawImage(tempScreen, op)

		// Draw the previous Content (during pop animation)
		if n.direction == -1 && len(n.stack) > 1 {
			previousContent := n.stack[len(n.stack)-2]
			tempScreen := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
			previousContent.Draw(tempScreen)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(screen.Bounds().Dx())+offsetX, 0)
			screen.DrawImage(tempScreen, op)
		}

	} else {
		// No animation, draw normally
		n.stack[len(n.stack)-1].Draw(screen)
	}
}

// SimpleContent struct implements Content interface, allowing for specific drawing logic.
type SimpleContent struct {
	color   color.Color
	message string // A unique message for each content
}

// NewSimpleContent creates a new SimpleContent with the given color and message.
func NewSimpleContent(color color.Color, message string) *SimpleContent {
	return &SimpleContent{
		color:   color,
		message: message,
	}
}

// Update handles updating content state, such as animations or game logic (no-op in this case).
func (c *SimpleContent) Update() error {
	return nil
}

// Draw fills the screen with the content's color and adds debug text to identify the content.
func (c *SimpleContent) Draw(screen *ebiten.Image) {
	screen.Fill(c.color)
	ebitenutil.DebugPrint(screen, c.message) // Add debug print to differentiate contents
}

// Game struct manages the main game loop and interaction with Navigator.
type Game struct {
	navigator *Navigator
}

// Update handles game state updates, including input handling for navigation.
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyP) { // Push with 'P'
		g.navigator.Push(NewSimpleContent(color.RGBA{0, 0, 255, 255}, "Blue Screen"))
	}
	if ebiten.IsKeyPressed(ebiten.KeyO) { // Pop with 'O'
		g.navigator.Pop()
	}
	return g.navigator.Update()
}

// Draw calls the navigator to render the current content.
func (g *Game) Draw(screen *ebiten.Image) {
	g.navigator.Draw(screen)
}

// Layout sets the game's layout dimensions, maintaining the aspect ratio.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// main initializes the game and starts the Ebiten game loop.
func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Navigator Example")

	navigator := NewNavigator()
	navigator.Push(NewSimpleContent(color.RGBA{255, 0, 0, 255}, "Red Screen"))

	game := &Game{navigator: navigator}

	// Example of pushing another content after 2 seconds
	go func() {
		<-time.After(2 * time.Second)
		navigator.Push(NewSimpleContent(color.RGBA{0, 255, 0, 255}, "Green Screen"))
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
