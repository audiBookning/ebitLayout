package layout

type Cell01 struct {
	X, Y   int
	Width  int
	Height int
}

type Grid01 struct {
	Cells          [][]Cell01
	Rows, Cols     int
	RowHeights     []float32
	ColWidths      []float32
	TotalWidth     int
	TotalHeight    int
	BodyPadding    int
	CellMargin     int
	CellBorderSize int
}

func NewGrid01(rows, cols int, totalWidth, totalHeight, bodyPadding, cellMargin, cellBorderSize int) *Grid01 {
	rowHeights := make([]float32, rows+1)
	colWidths := make([]float32, cols+1)

	// Set row heights to fill the entire screen height
	for i := 0; i <= rows; i++ {
		rowHeights[i] = float32(i) / float32(rows)
	}

	// Set column widths to fill the entire screen width
	for i := 0; i <= cols; i++ {
		colWidths[i] = float32(i) / float32(cols)
	}

	grid := &Grid01{
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

func (g *Grid01) InitializeCells() {
	g.Cells = make([][]Cell01, g.Rows)
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = make([]Cell01, g.Cols)
		for x := 0; x < g.Cols; x++ {
			cellWidth := int(float32(g.TotalWidth-g.BodyPadding*2) * (g.ColWidths[x+1] - g.ColWidths[x]))
			cellHeight := (g.TotalHeight - g.BodyPadding*2) / g.Rows
			g.Cells[y][x] = Cell01{
				X:      g.BodyPadding + int(float32(g.TotalWidth-g.BodyPadding*2)*g.ColWidths[x]),
				Y:      g.BodyPadding + y*cellHeight,
				Width:  cellWidth,
				Height: cellHeight,
			}
		}
	}
}
