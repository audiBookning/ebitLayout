package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth    = 800
	screenHeight   = 600
	contentWidth   = 400
	contentHeight  = 300
	animationSpeed = 0.1 // Speed of the transition
)

type Game struct {
	currentIdx        int
	nextIdx           int
	animating         bool
	animationProgress float64
	contents          []Content
}

func (g *Game) Update() error {
	if g.animating {
		g.animationProgress += animationSpeed
		if g.animationProgress >= 1 {
			g.animating = false
			g.currentIdx = g.nextIdx
			g.animationProgress = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.currentIdx < len(g.contents)-1 {
		g.nextIdx = g.currentIdx + 1
		g.animating = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.currentIdx > 0 {
		g.nextIdx = g.currentIdx - 1
		g.animating = true
	}

	// Update positions of contents based on the current animation state
	for i := range g.contents {
		g.contents[i].X = g.calculatePosition(i)
	}

	return nil
}

func (g *Game) calculatePosition(idx int) float64 {
	if g.animating {
		if g.nextIdx > g.currentIdx {
			if idx == g.currentIdx {
				return float64(screenWidth/2 - contentWidth/2 - contentWidth*g.animationProgress) // Current sliding out
			}
			if idx == g.nextIdx {
				return float64(screenWidth/2 - contentWidth/2 + contentWidth*(1-g.animationProgress)) // Next sliding in
			}
		} else {
			if idx == g.currentIdx {
				return float64(screenWidth/2 - contentWidth/2 + contentWidth*g.animationProgress) // Current sliding out
			}
			if idx == g.nextIdx {
				return float64(screenWidth/2 - contentWidth/2 - contentWidth*(1-g.animationProgress)) // Next sliding in
			}
		}
	}
	if idx == g.currentIdx {
		return float64(screenWidth/2 - contentWidth/2) // Center position for current content
	}
	return float64(-contentWidth) // Off-screen for others
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	for i := range g.contents {
		g.contents[i].Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type Content struct {
	Content string
	Color   color.Color
	X       float64
}

func (c *Content) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(c.X), float32(screenHeight/2-contentHeight/2), contentWidth, contentHeight, c.Color, true)
	ebitenutil.DebugPrintAt(screen, c.Content, int(c.X)+10, screenHeight/2)
}

func main() {
	content1 := Content{
		Content: "Content 1",
		Color:   color.RGBA{255, 0, 0, 255},
	}

	content2 := Content{
		Content: "Content 2",
		Color:   color.RGBA{0, 255, 0, 255},
	}

	content3 := Content{
		Content: "Content 3",
		Color:   color.RGBA{0, 0, 255, 255},
	}

	game := &Game{
		contents: []Content{
			content1, content2, content3,
		},
		currentIdx: 0,
		nextIdx:    0,
		animating:  false,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Horizontal Slider with Animation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
