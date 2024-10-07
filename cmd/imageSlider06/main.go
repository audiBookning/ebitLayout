package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	if len(n.stack) > 1 {
		n.stack = n.stack[:len(n.stack)-1]
	}
}

func (n *Navigator) CurrentScreen() Screen {
	if len(n.stack) == 0 {
		return nil
	}
	return n.stack[len(n.stack)-1]
}

type ContentScreen struct {
	navigator         *Navigator
	contents          []Content
	currentIdx        int
	prevIdx           int
	animating         bool
	animationProgress float64
	animationSpeed    float64
}

func NewContentScreen(navigator *Navigator, contents []Content, currentIdx int) *ContentScreen {
	return &ContentScreen{
		navigator:      navigator,
		contents:       contents,
		currentIdx:     currentIdx,
		prevIdx:        currentIdx,
		animationSpeed: 0.05, // Adjust animation speed as needed
	}
}

func (s *ContentScreen) Update(navigator *Navigator) error {
	if s.animating {
		s.animationProgress += s.animationSpeed
		if s.animationProgress >= 1 {
			s.animating = false
			s.prevIdx = s.currentIdx
			s.animationProgress = 0
			if s.prevIdx > s.currentIdx {
				navigator.Pop() // Pop the previous screen after the animation finishes
			}
		}
		return nil
	}

	// Handle right arrow key to move forward
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && s.currentIdx < len(s.contents)-1 {
		nextIdx := s.currentIdx + 1
		s.startAnimation(nextIdx)
		nextScreen := NewContentScreen(navigator, s.contents, nextIdx)
		navigator.Push(nextScreen)
	}

	// Handle left arrow key to move backward
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if len(navigator.stack) > 1 {
			prevIdx := s.currentIdx - 1
			s.startAnimation(prevIdx)
		}
	}

	return nil
}

func (s *ContentScreen) startAnimation(nextIdx int) {
	s.prevIdx = s.currentIdx
	s.currentIdx = nextIdx
	s.animating = true
	s.animationProgress = 0
}

func (s *ContentScreen) calculatePosition(idx int) (int, bool) {
	currentContent := s.contents[s.prevIdx]
	nextContent := s.contents[s.currentIdx]

	if s.animating {
		if s.currentIdx > s.prevIdx { // Moving forward
			if idx == s.prevIdx {
				position := 400 - currentContent.width/2 - int(float64(currentContent.width)*s.animationProgress)
				return position, position+currentContent.width > 0
			}
			if idx == s.currentIdx {
				position := 800 - int(float64(nextContent.width)*s.animationProgress)
				return position, position < 800
			}
		} else { // Moving backward
			if idx == s.prevIdx {
				position := int(float64(currentContent.width)*s.animationProgress) - currentContent.width
				return position, position+currentContent.width > 0
			}
			if idx == s.currentIdx {
				position := 400 - nextContent.width/2 - int(float64(nextContent.width)*(1-s.animationProgress))
				return position, position+nextContent.width > 0
			}
		}
	}

	if idx == s.currentIdx {
		position := 400 - currentContent.width/2
		return position, true
	}

	return -s.contents[idx].width, false
}

func (s *ContentScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	for i := range s.contents {
		position, visible := s.calculatePosition(i)
		if visible {
			s.contents[i].X = position
			s.contents[i].Draw(screen, 600)
		}
	}

	ebitenutil.DebugPrintAt(screen, "Use LEFT and RIGHT arrow keys to navigate", 10, 580)
}

type Content struct {
	text    string
	Color   color.Color
	X       int
	Y       int
	width   int
	height  int
	Visible bool
}

func NewContent(text string, color color.Color, width, height int) Content {
	return Content{
		text:    text,
		Color:   color,
		width:   width,
		height:  height,
		Visible: false,
	}
}

func (c *Content) Draw(screen *ebiten.Image, screenHeight int) {
	vector.DrawFilledRect(screen, float32(c.X), float32(screenHeight/2-c.height/2), float32(c.width), float32(c.height), c.Color, true)
	ebitenutil.DebugPrintAt(screen, c.text, c.X+10, screenHeight/2)
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
	return 800, 600
}

func main() {
	navigator := NewNavigator()

	// Define content for screens
	contents := []Content{
		NewContent("Content 1", color.RGBA{255, 0, 0, 255}, 400, 300),
		NewContent("Content 2", color.RGBA{0, 255, 0, 255}, 350, 250),
		NewContent("Content 3", color.RGBA{0, 0, 255, 255}, 450, 350),
	}

	// Push the initial content screen
	initialScreen := NewContentScreen(navigator, contents, 0)
	navigator.Push(initialScreen)

	game := &Game{navigator: navigator}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Horizontal Slider with Animation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
