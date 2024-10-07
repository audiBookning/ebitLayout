package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"example.com/menu/internals/layout"
)

const (
	initialScreenWidth  = 640
	initialScreenHeight = 480
	initialRows         = 3
	initialCols         = 4
	cellSize            = 100
)

type Game struct {
	grid         *layout.Grid03
	screenWidth  int
	screenHeight int
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.grid.AddRow()
		g.screenHeight += cellSize
		g.updateGridSize()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.grid.AddColumn()
		g.screenWidth += cellSize
		g.updateGridSize()
	}
	return nil
}

func (g *Game) updateGridSize() {
	newGrid := layout.NewGrid03(g.grid.Rows, g.grid.Cols, g.screenWidth, g.screenHeight, 10, 2, 1)
	g.grid = newGrid
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
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
	return g.screenWidth, g.screenHeight
}

func main() {
	game := &Game{
		screenWidth:  initialScreenWidth,
		screenHeight: initialScreenHeight,
	}
	game.grid = layout.NewGrid03(initialRows, initialCols, game.screenWidth, game.screenHeight, 10, 2, 1)

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Grid Example")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
