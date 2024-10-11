package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	screenWidth       int
	screenHeight      int
	animationSpeed    float64
	currentIdx        int
	nextIdx           int
	animating         bool
	animationProgress float64
	contents          []GeneralScreen
}

func NewGame() *Game {
	return &Game{
		screenWidth:    800,
		screenHeight:   600,
		animationSpeed: 0.1,
		contents: []GeneralScreen{
			NewContent("Content 1", color.RGBA{255, 0, 0, 255}, 400, 300),
			NewContent("Content 2", color.RGBA{0, 255, 0, 255}, 350, 250),
			NewContent("Content 3", color.RGBA{0, 0, 255, 255}, 450, 350),
		},
		currentIdx: 0,
		nextIdx:    0,
		animating:  false,
	}
}

func (g *Game) Update() error {
	if g.animating {
		g.animationProgress += g.animationSpeed
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

	for i := range g.contents {
		position, isVisible := g.calculatePosition(i)

		g.contents[i].X = position
		g.contents[i].Visible = isVisible

	}

	return nil
}

func (g *Game) calculatePosition(idx int) (int, bool) {
	currentContent := g.contents[g.currentIdx]
	nextContent := g.contents[g.nextIdx]

	if g.animating {
		if g.nextIdx > g.currentIdx {
			if idx == g.currentIdx {
				position := g.screenWidth/2 - currentContent.width/2 - int(float64(currentContent.width)*g.animationProgress)
				return position, position+currentContent.width > 0 // Current sliding out
			}
			if idx == g.nextIdx {
				position := g.screenWidth/2 - nextContent.width/2 + int(float64(nextContent.width)*(1-g.animationProgress))
				return position, position < g.screenWidth // Next sliding in
			}
		} else {
			if idx == g.currentIdx {
				position := g.screenWidth/2 - currentContent.width/2 + int(float64(currentContent.width)*g.animationProgress)
				return position, position < g.screenWidth // Current sliding out
			}
			if idx == g.nextIdx {
				position := g.screenWidth/2 - nextContent.width/2 - int(float64(nextContent.width)*(1-g.animationProgress))
				return position, position+nextContent.width > 0 // Next sliding in
			}
		}
	}

	if idx == g.currentIdx {
		position := g.screenWidth/2 - currentContent.width/2 // Center position for current content
		visible := position+currentContent.width > 0 && position < g.screenWidth
		return position, visible
	}

	// Off-screen for others
	return -g.contents[idx].width, false
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	for i := range g.contents {
		if g.contents[i].Visible {
			g.contents[i].Draw(screen, g.screenHeight)
		}
	}

	// Add navigation instructions at the bottom of the screen
	ebitenutil.DebugPrintAt(screen, "Use LEFT and RIGHT arrow keys to navigate", 10, g.screenHeight-20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

type GeneralScreen struct {
	text    string
	Color   color.Color
	X       int
	Y       int
	width   int
	height  int
	Visible bool // New field to store visibility state
}

func NewContent(text string, color color.Color, width, height int) GeneralScreen {
	return GeneralScreen{
		text:    text,
		Color:   color,
		width:   width,
		height:  height,
		Visible: false,
	}
}

func (c *GeneralScreen) Draw(screen *ebiten.Image, screenHeight int) {
	// Draw the content's rectangle
	vector.DrawFilledRect(screen, float32(c.X), float32(screenHeight/2-c.height/2), float32(c.width), float32(c.height), c.Color, true)

	// Draw the content's text
	ebitenutil.DebugPrintAt(screen, c.text, c.X+10, screenHeight/2)
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Horizontal Slider with Animation")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
