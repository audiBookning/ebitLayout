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
	Cells          [][]Cell
	Rows, Cols     int
	RowHeights     []float32
	ColWidths      []float32
	TotalWidth     int
	TotalHeight    int
	BodyPadding    int
	CellMargin     int
	CellBorderSize int
}

func NewGrid(rows, cols int, totalWidth, totalHeight, bodyPadding, cellMargin, cellBorderSize int) *Grid {
	rowHeights := make([]float32, rows+1)
	colWidths := make([]float32, cols+1)

	// Set dynamic proportions for row heights and column widths
	for i := 0; i <= rows; i++ {
		rowHeights[i] = float32(i) / float32(rows)
	}
	for i := 0; i <= cols; i++ {
		colWidths[i] = float32(i) / float32(cols)
	}

	grid := &Grid{
		Rows:           rows,
		Cols:           cols,
		RowHeights:     rowHeights,
		ColWidths:      colWidths,
		TotalWidth:     totalWidth,
		TotalHeight:    totalHeight,
		BodyPadding:    bodyPadding,
		CellMargin:     cellMargin,
		CellBorderSize: cellBorderSize,
	}
	grid.InitializeCells()
	return grid
}

func (g *Grid) InitializeCells() {
	g.Cells = make([][]Cell, g.Rows)
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = make([]Cell, g.Cols)
		for x := 0; x < g.Cols; x++ {
			cellWidth := int(float32(g.TotalWidth-g.BodyPadding*2) * (g.ColWidths[x+1] - g.ColWidths[x]))
			cellHeight := int(float32(g.TotalHeight-g.BodyPadding*2) * (g.RowHeights[y+1] - g.RowHeights[y]))
			g.Cells[y][x] = Cell{
				X:      g.BodyPadding + int(float32(g.TotalWidth-g.BodyPadding*2)*g.ColWidths[x]),
				Y:      g.BodyPadding + int(float32(g.TotalHeight-g.BodyPadding*2)*g.RowHeights[y]),
				Width:  cellWidth,
				Height: cellHeight,
			}
		}
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.HandleClick(x, y)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	borderColor := color.RGBA{255, 0, 0, 255}

	for _, row := range g.grid.Cells {
		for _, cell := range row {
			// Calculate the effective area for the cell including the margin and border size
			x := float32(cell.X) + float32(g.grid.CellMargin/2)
			y := float32(cell.Y) + float32(g.grid.CellMargin/2)
			w := float32(cell.Width) - float32(g.grid.CellMargin)
			h := float32(cell.Height) - float32(g.grid.CellMargin)

			// Draw the cell background (excluding borders)
			vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{255, 255, 255, 255}, true)

			// Draw the borders inside the cell area
			borderSize := float32(g.grid.CellBorderSize)
			vector.StrokeRect(screen, x+borderSize/2, y+borderSize/2, w-borderSize, h-borderSize, borderSize, borderColor, true)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Recalculate grid layout based on the new screen size
	g.grid.TotalWidth = outsideWidth
	g.grid.TotalHeight = outsideHeight
	g.grid.InitializeCells()
	return outsideWidth, outsideHeight
}

func (g *Game) HandleClick(x, y int) {
	for _, row := range g.grid.Cells {
		for _, cell := range row {
			if x >= cell.X && x < cell.X+cell.Width && y >= cell.Y && y < cell.Y+cell.Height {
				fmt.Printf("Clicked on cell at row %d, col %d\n", cell.Y/cell.Height, cell.X/cell.Width)
				return
			}
		}
	}
}

type Game struct {
	grid *Grid
}

func main() {
	// Define the proportions of rows and columns (like CSS Grid Template)
	rows := 3
	cols := 5

	bodyPadding := 20
	cellMargin := 3
	cellBorderSize := 5 // Customizable border size

	game := &Game{
		grid: NewGrid(rows, cols, 640, 480, bodyPadding, cellMargin, cellBorderSize),
	}

	// ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.RunGame(game)
}
