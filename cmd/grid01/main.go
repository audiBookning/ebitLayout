package main

import (
	"fmt"
	"image/color"
	"log"

	"example.com/menu/internals/layout"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	initialRows  = 3
	initialCols  = 4
)

type Game struct {
	grid *layout.Grid03
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.grid.AddRow()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.grid.AddColumn()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < g.grid.Rows; y++ {
		for x := 0; x < g.grid.Cols; x++ {
			cell := g.grid.Cells[y][x]
			ebitenutil.DrawRect(screen, float64(cell.X), float64(cell.Y), float64(cell.Width), float64(cell.Height), color.RGBA{100, 100, 100, 255})
			ebitenutil.DrawRect(screen, float64(cell.X)+1, float64(cell.Y)+1, float64(cell.Width)-2, float64(cell.Height)-2, color.RGBA{200, 200, 200, 255})

			text := fmt.Sprintf("%d,%d", x, y)
			ebitenutil.DebugPrintAt(screen, text, cell.X+5, cell.Y+5)
		}
	}

	ebitenutil.DebugPrint(screen, "Press 'R' to add a row, 'C' to add a column")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Grid Example")

	game := &Game{
		grid: layout.NewGrid03(initialRows, initialCols, screenWidth, screenHeight, 10, 2, 1),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
