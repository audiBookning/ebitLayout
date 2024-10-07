package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	SliderX      float64
	SliderY      float64
	SliderWidth  float64
	SliderHeight float64
	HandleX      float64
	HandleWidth  float64
	Dragging     bool
}

func (g *Game) Update() error {
	// Get the mouse position
	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Check if the mouse is within the slider bounds
		if float64(mx) >= g.SliderX && float64(mx) <= g.SliderX+g.SliderWidth && float64(my) >= g.SliderY && float64(my) <= g.SliderY+g.SliderHeight {
			// If clicking on the handle
			if float64(mx) >= g.HandleX && float64(mx) <= g.HandleX+g.HandleWidth {
				g.Dragging = true
			} else {
				// Move the handle to where the click occurred
				g.HandleX = float64(mx) - g.HandleWidth/2
				// Ensure the handle stays within the slider bounds
				if g.HandleX < g.SliderX {
					g.HandleX = g.SliderX
				}
				if g.HandleX+g.HandleWidth > g.SliderX+g.SliderWidth {
					g.HandleX = g.SliderX + g.SliderWidth - g.HandleWidth
				}
			}
		}
	}

	if g.Dragging {
		// Update the handle position based on mouse movement
		g.HandleX = float64(mx) - g.HandleWidth/2

		// Ensure the handle stays within the slider bounds
		if g.HandleX < g.SliderX {
			g.HandleX = g.SliderX
		}
		if g.HandleX+g.HandleWidth > g.SliderX+g.SliderWidth {
			g.HandleX = g.SliderX + g.SliderWidth - g.HandleWidth
		}
	}

	// Stop dragging when the mouse button is released
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.Dragging = false
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the slider bar
	vector.DrawFilledRect(screen, float32(g.SliderX), float32(g.SliderY), float32(g.SliderWidth), float32(g.SliderHeight), color.RGBA{0x80, 0x80, 0x80, 0xff}, true)

	// Draw the handle
	vector.DrawFilledRect(screen, float32(g.HandleX), float32(g.SliderY), float32(g.HandleWidth), float32(g.SliderHeight), color.RGBA{0xff, 0x00, 0x00, 0xff}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func main() {
	game := &Game{
		SliderX:      100,
		SliderY:      300,
		SliderWidth:  600,
		SliderHeight: 20,
		HandleX:      100,
		HandleWidth:  40,
		Dragging:     false,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Horizontal Slider Example")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}