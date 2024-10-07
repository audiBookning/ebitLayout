package main

import (
	"fmt"
	"image/color"

	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	grid *layout.Grid02
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.HandleClick(x, y)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	margin := 2 // Margin between cells
	borderColor := color.RGBA{0, 0, 0, 255}

	for _, row := range g.grid.Cells {
		for _, cell := range row {
			x := float32(cell.X)
			y := float32(cell.Y)
			w := float32(cell.Width)
			h := float32(cell.Height)

			// Draw the cell background with a margin
			vector.DrawFilledRect(screen, x+float32(margin), y+float32(margin), w-float32(2*margin), h-float32(2*margin), color.RGBA{255, 255, 255, 255}, true)

			// Draw the cell borders
			// Top border
			vector.DrawFilledRect(screen, x, y, w, float32(margin), borderColor, true)
			// Bottom border
			vector.DrawFilledRect(screen, x, y+h-float32(margin), w, float32(margin), borderColor, true)
			// Left border
			vector.DrawFilledRect(screen, x, y, float32(margin), h, borderColor, true)
			// Right border
			vector.DrawFilledRect(screen, x+w-float32(margin), y, float32(margin), h, borderColor, true)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) HandleClick(x, y int) {
	for _, row := range g.grid.Cells {
		for _, cell := range row {
			if x >= cell.X && x < cell.X+cell.Width && y >= cell.Y && y < cell.Y+cell.Height {
				fmt.Printf("Clicked on cell at row %d, col %d\n", cell.Y/g.grid.CellHeight, cell.X/g.grid.CellWidth)
				return
			}
		}
	}
}

func main() {
	game := &Game{grid: layout.NewGrid02(10, 10, 50, 50)}
	ebiten.RunGame(game)
}
