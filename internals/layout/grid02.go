package layout

type Cell02 struct {
	X, Y   int
	Width  int
	Height int
}

type Grid02 struct {
	Cells      [][]Cell02
	Rows, Cols int
	CellWidth  int
	CellHeight int
}

func NewGrid02(rows, cols, cellWidth, cellHeight int) *Grid02 {
	grid := &Grid02{
		Rows:       rows,
		Cols:       cols,
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
	}
	grid.InitializeCells()
	return grid
}

func (g *Grid02) InitializeCells() {
	g.Cells = make([][]Cell02, g.Rows)
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = make([]Cell02, g.Cols)
		for x := 0; x < g.Cols; x++ {
			g.Cells[y][x] = Cell02{
				X:      x * g.CellWidth,
				Y:      y * g.CellHeight,
				Width:  g.CellWidth,
				Height: g.CellHeight,
			}
		}
	}
}

func (g *Grid02) AddRow() {
	g.Rows++
	row := make([]Cell02, g.Cols)
	for x := 0; x < g.Cols; x++ {
		row[x] = Cell02{
			X:      x * g.CellWidth,
			Y:      (g.Rows - 1) * g.CellHeight,
			Width:  g.CellWidth,
			Height: g.CellHeight,
		}
	}
	g.Cells = append(g.Cells, row)
}

func (g *Grid02) AddColumn() {
	g.Cols++
	for y := 0; y < g.Rows; y++ {
		g.Cells[y] = append(g.Cells[y], Cell02{
			X:      (g.Cols - 1) * g.CellWidth,
			Y:      y * g.CellHeight,
			Width:  g.CellWidth,
			Height: g.CellHeight,
		})
	}
}
