package layout

import (
	"log"
	"math"
)

type Cell struct {
	X, Y          int
	Width, Height int
}

type Grid struct {
	Cells               [][]Cell
	Rows, Cols          int
	RowHeights          []float32
	ColWidths           []float32
	TotalWidth          int
	TotalHeight         int
	BodyPadding         int
	CellMargin          int
	CellBorderSize      int
	MaintainAspectRatio bool
	InitialCellWidth    int
	InitialCellHeight   int
	InitialTotalWidth   int
	InitialTotalHeight  int
}

func NewGrid(rows, cols int, totalWidth, totalHeight, bodyPadding, cellMargin, cellBorderSize int, maintainAspectRatio bool) *Grid {
	rowHeights := make([]float32, rows+1)
	colWidths := make([]float32, cols+1)

	for i := 0; i <= rows; i++ {
		rowHeights[i] = float32(i) / float32(rows)
	}

	for i := 0; i <= cols; i++ {
		colWidths[i] = float32(i) / float32(cols)
	}

	grid := &Grid{
		Rows:                rows,
		Cols:                cols,
		RowHeights:          rowHeights,
		ColWidths:           colWidths,
		TotalWidth:          totalWidth,
		TotalHeight:         totalHeight,
		BodyPadding:         bodyPadding,
		CellMargin:          cellMargin,
		CellBorderSize:      cellBorderSize,
		MaintainAspectRatio: maintainAspectRatio,
		InitialCellWidth:    (totalWidth - bodyPadding*2) / cols,
		InitialCellHeight:   (totalHeight - bodyPadding*2) / rows,
		InitialTotalWidth:   totalWidth,
		InitialTotalHeight:  totalHeight,
	}
	grid.InitializeCells(totalWidth, totalHeight)
	return grid
}

func (g *Grid) InitializeCells(currentWidth, currentHeight int) {
	g.Cells = make([][]Cell, g.Rows)
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = make([]Cell, g.Cols)
		for x := 0; x < g.Cols; x++ {

			cellWidth := (g.InitialTotalWidth - g.BodyPadding*2) / g.Cols
			cellHeight := (g.InitialTotalHeight - g.BodyPadding*2) / g.Rows

			if g.MaintainAspectRatio {
				scaleX := float64(currentWidth) / float64(g.InitialTotalWidth)
				scaleY := float64(currentHeight) / float64(g.InitialTotalHeight)
				scale := math.Min(scaleX, scaleY)

				cellWidth = int(float64(cellWidth) * scale)
				cellHeight = int(float64(cellHeight) * scale)
			}

			g.Cells[y][x] = Cell{
				X:      g.BodyPadding + x*cellWidth,
				Y:      g.BodyPadding + y*cellHeight,
				Width:  cellWidth - g.CellMargin,
				Height: cellHeight - g.CellMargin,
			}
		}
	}
}

func (g *Grid) AddRow() {
	g.Rows++
	updatedRowHeights := make([]float32, g.Rows+1)
	for i := 0; i <= g.Rows; i++ {
		updatedRowHeights[i] = float32(i) / float32(g.Rows)
	}
	g.RowHeights = updatedRowHeights

	newRow := make([]Cell, g.Cols)
	for x := 0; x < g.Cols; x++ {

		cellWidth := (g.InitialTotalWidth - g.BodyPadding*2) / g.Cols
		cellHeight := (g.InitialTotalHeight - g.BodyPadding*2) / g.Rows

		if g.MaintainAspectRatio {
			scaleX := float64(g.TotalWidth) / float64(g.InitialTotalWidth)
			scaleY := float64(g.TotalHeight) / float64(g.InitialTotalHeight)
			scale := math.Min(scaleX, scaleY)

			cellWidth = int(float64(cellWidth) * scale)
			cellHeight = int(float64(cellHeight) * scale)
		}

		newRow[x] = Cell{
			X:      g.BodyPadding + x*cellWidth,
			Y:      g.BodyPadding + (g.Rows-1)*cellHeight,
			Width:  cellWidth - g.CellMargin,
			Height: cellHeight - g.CellMargin,
		}
	}
	g.Cells = append(g.Cells, newRow)
	g.UpdateCellPositions()
}

func (g *Grid) AddColumn() {
	g.Cols++
	updatedColWidths := make([]float32, g.Cols+1)
	for i := 0; i <= g.Cols; i++ {
		updatedColWidths[i] = float32(i) / float32(g.Cols)
	}
	g.ColWidths = updatedColWidths

	for y := 0; y < g.Rows; y++ {

		cellWidth := (g.InitialTotalWidth - g.BodyPadding*2) / g.Cols
		cellHeight := (g.InitialTotalHeight - g.BodyPadding*2) / g.Rows

		if g.MaintainAspectRatio {
			scaleX := float64(g.TotalWidth) / float64(g.InitialTotalWidth)
			scaleY := float64(g.TotalHeight) / float64(g.InitialTotalHeight)
			scale := math.Min(scaleX, scaleY)

			cellWidth = int(float64(cellWidth) * scale)
			cellHeight = int(float64(cellHeight) * scale)
		}

		newCell := Cell{
			X:      g.BodyPadding + (g.Cols-1)*cellWidth,
			Y:      g.BodyPadding + y*cellHeight,
			Width:  cellWidth - g.CellMargin,
			Height: cellHeight - g.CellMargin,
		}
		g.Cells[y] = append(g.Cells[y], newCell)
	}
	g.UpdateCellPositions()
}

func (g *Grid) UpdateCellPositions() {
	for y := 0; y < g.Rows; y++ {
		for x := 0; x < g.Cols; x++ {

			cellWidth := (g.InitialTotalWidth - g.BodyPadding*2) / g.Cols
			cellHeight := (g.InitialTotalHeight - g.BodyPadding*2) / g.Rows

			if g.MaintainAspectRatio {
				scaleX := float64(g.TotalWidth) / float64(g.InitialTotalWidth)
				scaleY := float64(g.TotalHeight) / float64(g.InitialTotalHeight)
				scale := math.Min(scaleX, scaleY)

				cellWidth = int(float64(cellWidth) * scale)
				cellHeight = int(float64(cellHeight) * scale)
			}

			g.Cells[y][x] = Cell{
				X:      g.BodyPadding + x*cellWidth,
				Y:      g.BodyPadding + y*cellHeight,
				Width:  cellWidth - g.CellMargin,
				Height: cellHeight - g.CellMargin,
			}

			log.Printf("Cell[%d][%d] - X: %d, Y: %d, Width: %d, Height: %d",
				y, x, g.Cells[y][x].X, g.Cells[y][x].Y, g.Cells[y][x].Width, g.Cells[y][x].Height)
		}
	}
}
