package layout

import "log"

type Cell03 struct {
	X, Y   int
	Width  int
	Height int
}

type Grid03 struct {
	Cells          [][]Cell03
	Rows, Cols     int
	RowHeights     []float32
	ColWidths      []float32
	TotalWidth     int
	TotalHeight    int
	BodyPadding    int
	CellMargin     int
	CellBorderSize int
	CellWidth      int
	CellHeight     int
}

func NewGrid03(rows, cols, totalWidth, totalHeight, bodyPadding, cellMargin, cellBorderSize int) *Grid03 {
	rowHeights := make([]float32, rows+1)
	colWidths := make([]float32, cols+1)

	for i := 0; i <= rows; i++ {
		rowHeights[i] = float32(i) / float32(rows)
	}

	for i := 0; i <= cols; i++ {
		colWidths[i] = float32(i) / float32(cols)
	}

	cellWidth := (totalWidth - bodyPadding*2) / cols
	cellHeight := (totalHeight - bodyPadding*2) / rows

	grid := &Grid03{
		Rows:           rows,
		Cols:           cols,
		RowHeights:     rowHeights,
		ColWidths:      colWidths,
		TotalWidth:     totalWidth,
		TotalHeight:    totalHeight,
		BodyPadding:    bodyPadding,
		CellMargin:     cellMargin,
		CellBorderSize: cellBorderSize,
		CellWidth:      cellWidth,
		CellHeight:     cellHeight,
	}
	grid.InitializeCells()
	return grid
}

func (g *Grid03) InitializeCells() {
	g.Cells = make([][]Cell03, g.Rows)
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = make([]Cell03, g.Cols)
		for x := 0; x < g.Cols; x++ {
			cellWidth := int(float32(g.TotalWidth-g.BodyPadding*2) * (g.ColWidths[x+1] - g.ColWidths[x]))
			cellHeight := int(float32(g.TotalHeight-g.BodyPadding*2) * (g.RowHeights[y+1] - g.RowHeights[y]))
			g.Cells[y][x] = Cell03{
				X:      g.BodyPadding + int(float32(g.TotalWidth-g.BodyPadding*2)*g.ColWidths[x]),
				Y:      g.BodyPadding + int(float32(g.TotalHeight-g.BodyPadding*2)*g.RowHeights[y]),
				Width:  cellWidth,
				Height: cellHeight,
			}
		}
	}
}

func (g *Grid03) AddRow() {
	g.Rows++
	updatedRowHeights := make([]float32, g.Rows+1)
	for i := 0; i <= g.Rows; i++ {
		updatedRowHeights[i] = float32(i) / float32(g.Rows)
	}
	g.RowHeights = updatedRowHeights

	newRow := make([]Cell03, g.Cols)
	for x := 0; x < g.Cols; x++ {
		cellWidth := int(float32(g.TotalWidth-g.BodyPadding*2) * (g.ColWidths[x+1] - g.ColWidths[x]))
		cellHeight := int(float32(g.TotalHeight-g.BodyPadding*2) * (g.RowHeights[g.Rows] - g.RowHeights[g.Rows-1]))
		newRow[x] = Cell03{
			X:      g.BodyPadding + int(float32(g.TotalWidth-g.BodyPadding*2)*g.ColWidths[x]),
			Y:      g.BodyPadding + int(float32(g.TotalHeight-g.BodyPadding*2)*g.RowHeights[g.Rows-1]),
			Width:  cellWidth,
			Height: cellHeight,
		}
	}
	g.Cells = append(g.Cells, newRow)
	g.UpdateCellPositions()
}

func (g *Grid03) AddColumn() {
	g.Cols++
	updatedColWidths := make([]float32, g.Cols+1)
	for i := 0; i <= g.Cols; i++ {
		updatedColWidths[i] = float32(i) / float32(g.Cols)
	}
	g.ColWidths = updatedColWidths

	for y := 0; y < g.Rows; y++ {
		cellWidth := int(float32(g.TotalWidth-g.BodyPadding*2) * (g.ColWidths[g.Cols] - g.ColWidths[g.Cols-1]))
		cellHeight := int(float32(g.TotalHeight-g.BodyPadding*2) * (g.RowHeights[y+1] - g.RowHeights[y]))
		newCell := Cell03{
			X:      g.BodyPadding + int(float32(g.TotalWidth-g.BodyPadding*2)*g.ColWidths[g.Cols-1]),
			Y:      g.BodyPadding + int(float32(g.TotalHeight-g.BodyPadding*2)*g.RowHeights[y]),
			Width:  cellWidth,
			Height: cellHeight,
		}
		g.Cells[y] = append(g.Cells[y], newCell)
	}
	g.UpdateCellPositions()
}

func (g *Grid03) UpdateCellPositions() {
	for y := 0; y < g.Rows; y++ {
		for x := 0; x < g.Cols; x++ {
			cellWidth := int(float32(g.TotalWidth-g.BodyPadding*2) * (g.ColWidths[x+1] - g.ColWidths[x]))
			cellHeight := int(float32(g.TotalHeight-g.BodyPadding*2) * (g.RowHeights[y+1] - g.RowHeights[y]))
			g.Cells[y][x] = Cell03{
				X:      g.BodyPadding + int(float32(g.TotalWidth-g.BodyPadding*2)*g.ColWidths[x]),
				Y:      g.BodyPadding + int(float32(g.TotalHeight-g.BodyPadding*2)*g.RowHeights[y]),
				Width:  cellWidth,
				Height: cellHeight,
			}
			log.Printf("Cell[%d][%d] - X: %d, Y: %d, Width: %d, Height: %d", y, x, g.Cells[y][x].X, g.Cells[y][x].Y, g.Cells[y][x].Width, g.Cells[y][x].Height)
		}
	}
}
