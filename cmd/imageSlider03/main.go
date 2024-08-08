package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth   = 800
	screenHeight  = 600
	rectWidth     = 200
	rectHeight    = 300
	sliderSpeed   = 10
	snapStrength  = 0.1
	dragFactor    = 0.9
	minVelocity   = 0.5
	positionWidth = rectWidth + 20
)

type Game struct {
	contentX     float64
	velocity     float64
	colors       []color.Color
	mouseX       int
	mouseDragged bool
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.velocity += sliderSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.velocity -= sliderSpeed
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		if !g.mouseDragged {
			g.mouseDragged = true
			g.mouseX = x
		}
		g.velocity = float64(x - g.mouseX)
		g.mouseX = x
	} else {
		g.mouseDragged = false
	}

	// Apply snap force towards the nearest fixed position
	nearestPosition := roundToNearestPosition(g.contentX)
	distance := nearestPosition - g.contentX
	force := distance * snapStrength
	g.velocity += force

	// Update position based on velocity
	g.contentX += g.velocity

	// Apply drag to the velocity
	g.velocity *= dragFactor

	// Stop the movement if the velocity is very low
	if math.Abs(g.velocity) < minVelocity {
		g.velocity = 0
		g.contentX = nearestPosition
	}

	return nil
}

func roundToNearestPosition(x float64) float64 {
	return math.Round(x/positionWidth) * positionWidth
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for i, col := range g.colors {
		rect := ebiten.NewImage(rectWidth, rectHeight)
		rect.Fill(col)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.contentX+float64(i*positionWidth), float64((screenHeight-rectHeight)/2))
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
		velocity: 0,
		colors:   []color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 255, 255}},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Horizontal Manual Slider with Inertia")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
