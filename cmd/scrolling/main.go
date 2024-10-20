package main

import (
	"bytes"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	textHeight   = 24
	linesVisible = 20
)

var (
	textArea   = strings.Repeat("This is a line of text.\n", 50)
	scrollY    = 0
	scrollMax  = (50 - linesVisible) * textHeight
	isDragging = false
	startDragY = 0
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type Game struct{}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	_, dy := ebiten.Wheel()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if isDragging {
			scrollY += (y - startDragY)
			startDragY = y
		} else if x >= screenWidth-20 && x <= screenWidth-10 {
			isDragging = true
			startDragY = y
		}

		if scrollY < 0 {
			scrollY = 0
		}
		if scrollY > scrollMax {
			scrollY = scrollMax
		}
	} else {
		isDragging = false
		scrollY -= int(dy) * 80
		if ebiten.IsKeyPressed(ebiten.KeyUp) && scrollY > 0 {
			scrollY -= 10
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) && scrollY < scrollMax {
			scrollY += 10
		}

		if scrollY < 0 {
			scrollY = 0
		}
		if scrollY > scrollMax {
			scrollY = scrollMax
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{30, 30, 30, 255})

	lines := strings.Split(textArea, "\n")
	for i, line := range lines {

		lineY := i*textHeight - scrollY

		if lineY >= 0 && lineY < screenHeight {
			op := &text.DrawOptions{}
			op.ColorScale.ScaleWithColor(color.White)
			ss := &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   textHeight,
			}

			op.GeoM.Translate(0, float64(lineY))

			text.Draw(screen, line, ss, op)
		}
	}

	barHeight := screenHeight * linesVisible / len(lines)
	barY := scrollY * (screenHeight - barHeight) / scrollMax
	vector.DrawFilledRect(screen, screenWidth-20, float32(barY), 10, float32(barHeight), color.RGBA{0, 0, 0, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Scrolling Text Area with Scrollbar")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
