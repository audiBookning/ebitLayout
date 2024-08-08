package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Cell struct {
	X, Y   int
	Width  int
	Height int
}

type Grid struct {
	Cells      [][]Cell
	Rows, Cols int
	CellWidth  int
	CellHeight int
}

func NewGrid(rows, cols, cellWidth, cellHeight int) *Grid {
	grid := &Grid{
		Rows:       rows,
		Cols:       cols,
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
	}
	grid.InitializeCells()
	return grid
}

func (g *Grid) InitializeCells() {
	g.Cells = make([][]Cell, g.Rows)
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = make([]Cell, g.Cols)
		for x := 0; x < g.Cols; x++ {
			g.Cells[y][x] = Cell{
				X:      x * g.CellWidth,
				Y:      y * g.CellHeight,
				Width:  g.CellWidth,
				Height: g.CellHeight,
			}
		}
	}
}

func (g *Grid) AddRow() {
	g.Rows++
	row := make([]Cell, g.Cols)
	for x := 0; x < g.Cols; x++ {
		row[x] = Cell{
			X:      x * g.CellWidth,
			Y:      (g.Rows - 1) * g.CellHeight,
			Width:  g.CellWidth,
			Height: g.CellHeight,
		}
	}
	g.Cells = append(g.Cells, row)
}

func (g *Grid) AddColumn() {
	g.Cols++
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = append(g.Cells[y], Cell{
			X:      (g.Cols - 1) * g.CellWidth,
			Y:      y * g.CellHeight,
			Width:  g.CellWidth,
			Height: g.CellHeight,
		})
	}
}

type Game struct {
	grid *Grid
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
	game := &Game{grid: NewGrid(10, 10, 50, 50)}
	ebiten.RunGame(game)
}
