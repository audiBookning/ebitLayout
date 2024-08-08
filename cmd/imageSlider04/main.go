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
	screenWidth   = 800
	screenHeight  = 600
	contentWidth  = 400
	contentHeight = 300
)

type Game struct {
	currentIdx int
	contents   []Content
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.currentIdx < len(g.contents)-1 {
			g.currentIdx++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.currentIdx > 0 {
			g.currentIdx--
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Calculate the position for the current content slide
	x := float64(screenWidth/2 - contentWidth/2)

	// Draw the colored rectangle background for the current content
	vector.DrawFilledRect(screen, float32(x), float32(screenHeight/2-contentHeight/2), contentWidth, contentHeight, g.contents[g.currentIdx].Color, true)

	// Draw the content text
	ebitenutil.DebugPrintAt(screen, g.contents[g.currentIdx].Content, int(x)+10, screenHeight/2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type Content struct {
	Content string
	Color   color.Color
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
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Horizontal Slider with Colored Backgrounds")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
