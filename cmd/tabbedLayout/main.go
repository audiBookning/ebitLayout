package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	tabs     = []string{"Tab 1", "Tab 2", "Tab 3"}
	selected int
)

type Game struct{}

func (g *Game) Update() error {
	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		selected = (selected + 1) % len(tabs)
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		selected = (selected - 1 + len(tabs)) % len(tabs)
	}

	// Handle mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		for i := range tabs {
			if x >= i*100 && x < (i+1)*100 {
				selected = i
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{72, 169, 166, 255})

	// Draw tabs
	for i, tab := range tabs {
		tabColor := color.RGBA{212, 180, 131, 255}
		if i == selected {
			// Highlight selected tab
			tabColor = color.RGBA{228, 223, 218, 255}
		}
		vector.DrawFilledRect(screen, float32(i*100), 0, 100, 30, tabColor, true)
		ebitenutil.DebugPrintAt(screen, tab, i*100+10, 5)
	}

	// Draw content based on selected tab
	content := fmt.Sprintf("Content of %s", tabs[selected])
	ebitenutil.DebugPrintAt(screen, content, 10, 40) // Move content down
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Tabbed Layout Example")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

/*

random palete

115, 29, 216
72, 169, 166
228, 223, 218
212, 180, 131
193, 102, 107


207, 212, 197
238, 207, 212
239, 185, 203
230, 173, 236
194, 135, 232

*/
