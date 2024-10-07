package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	sliderSpeed  = 10
	rectWidth    = 200
	rectHeight   = 300
)

type Game struct {
	contentX     int
	colors       []color.Color
	keyPressed   bool
	mouseX       int
	mouseDragged bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.contentX -= sliderSpeed
		g.keyPressed = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.contentX += sliderSpeed
		g.keyPressed = true
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		if !g.mouseDragged {
			g.mouseDragged = true
			g.mouseX = x
		}
		g.contentX += x - g.mouseX
		g.mouseX = x
	} else {
		if g.mouseDragged {
			g.contentX = roundToNearestPosition(g.contentX)
		}
		g.mouseDragged = false
	}

	// Snap to the nearest fixed position when key is released
	if !ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft) && g.keyPressed {
		g.contentX = roundToNearestPosition(g.contentX)
		g.keyPressed = false
	}

	return nil
}

func roundToNearestPosition(x int) int {
	positionWidth := rectWidth + 20
	// Calculate the nearest position
	return ((x + positionWidth/2) / positionWidth) * positionWidth
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for i, col := range g.colors {
		rect := ebiten.NewImage(rectWidth, rectHeight)
		rect.Fill(col)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.contentX+i*(rectWidth+20)), float64((screenHeight-rectHeight)/2))
		screen.DrawImage(rect, op)
	}

	ebitenutil.DebugPrint(screen, "Use Left/Right arrows or drag with mouse to move the slider")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		contentX: 0,
		colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Horizontal Manual Slider")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
