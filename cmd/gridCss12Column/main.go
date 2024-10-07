package main

import (
	"fmt"
	"image/color"
	"math"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/textwrapper"
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

	// Set row heights to fill the entire screen height
	for i := 0; i <= rows; i++ {
		rowHeights[i] = float32(i) / float32(rows)
	}

	// Set column widths to fill the entire screen width
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
		InitialCellWidth:    0,
		InitialCellHeight:   0,
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
			// Calculate cell width and height based on current window size
			cellWidth := int(float32(currentWidth-g.BodyPadding*2) * (g.ColWidths[x+1] - g.ColWidths[x]))
			cellHeight := (currentHeight - g.BodyPadding*2) / g.Rows

			// Store initial cell dimensions if not already set
			if g.InitialCellWidth == 0 && g.InitialCellHeight == 0 {
				g.InitialCellWidth = cellWidth
				g.InitialCellHeight = cellHeight
			}

			if g.MaintainAspectRatio {
				// Calculate the aspect ratio based on initial dimensions
				//aspectRatio := float64(g.InitialCellWidth) / float64(g.InitialCellHeight)

				// Calculate scaling factors based on initial and current window sizes
				scaleX := float64(currentWidth) / float64(g.InitialTotalWidth)
				scaleY := float64(currentHeight) / float64(g.InitialTotalHeight)

				// Choose the smaller scale to maintain aspect ratio without overflowing
				scale := math.Min(scaleX, scaleY)

				// Apply scaling while maintaining aspect ratio
				scaledWidth := int(float64(g.InitialCellWidth) * scale)
				scaledHeight := int(float64(g.InitialCellHeight) * scale)

				cellWidth = scaledWidth
				cellHeight = scaledHeight
			}

			g.Cells[y][x] = Cell{
				X:      g.BodyPadding + int(float64(x)*float64(cellWidth)),
				Y:      g.BodyPadding + int(float64(y)*float64(cellHeight)),
				Width:  cellWidth,
				Height: cellHeight,
			}
		}
	}
}

type Game struct {
	grid        *Grid
	textWrapper *textwrapper.TextWrapper
	fontSize    float64
}

var filePathTxt string
var Assets_Relative_Path = "../../"

func GetFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
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
	pointColor := color.RGBA{0, 0, 255, 255} // Blue color for the point

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

			// Calculate the center position of the cell
			centerX := float64(cell.X) + float64(cell.Width)/2
			centerY := float64(cell.Y) + float64(cell.Height)/2

			// Measure the dimensions of the text
			textWidth, textHeight := g.textWrapper.MeasureText("42")

			// Calculate the top-left position to center the text
			textX := centerX - float64(textWidth)/2
			textY := centerY - float64(textHeight)/2

			// Draw the number at the centered position
			g.textWrapper.DrawText(screen, "42", textX, textY)

			// Draw a little blue point at the center of the cell
			var pointSize float32 = 2.0
			vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), pointSize, pointColor, true)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Recalculate grid layout based on the new screen size
	g.grid.TotalWidth = outsideWidth
	g.grid.TotalHeight = outsideHeight
	g.grid.InitializeCells(outsideWidth, outsideHeight)

	// Calculate the font size ratio based on the scaling factor
	initialWidth := g.grid.InitialTotalWidth
	initialHeight := g.grid.InitialTotalHeight

	scaleX := float64(outsideWidth) / float64(initialWidth)
	scaleY := float64(outsideHeight) / float64(initialHeight)

	var scale float64
	if g.grid.MaintainAspectRatio {
		scale = math.Min(scaleX, scaleY)
	} else {
		// If not maintaining aspect ratio, you might want to handle scaleX and scaleY separately
		// Here, we choose to scale based on width for simplicity
		scale = scaleX
	}

	// Update the font size based on the scale
	g.textWrapper.SetFontSize(g.fontSize * scale)

	return outsideWidth, outsideHeight
}

func (g *Game) HandleClick(x, y int) {
	for _, row := range g.grid.Cells {
		for _, cell := range row {
			if x >= cell.X && x < cell.X+cell.Width && y >= cell.Y && y < cell.Y+cell.Height {
				rowIndex := y / cell.Height
				colIndex := x / cell.Width
				fmt.Printf("Clicked on cell at row %d, col %d\n", rowIndex, colIndex)
				return
			}
		}
	}
}

func main() {
	_, filePathTxt, _, _ = runtime.Caller(0)
	fontpath := GetFilePath("assets/fonts/roboto_regularTTF.ttf")

	// Define the proportions of rows and columns (like CSS Grid Template)
	rows := 1
	cols := 8

	bodyPadding := 20
	cellMargin := 3
	cellBorderSize := 5

	fontSize := 24.0 // Initial font size

	maintainAspectRatio := false // Set to true to maintain cell aspect ratio

	// Initial window size
	initialWidth := 640
	initialHeight := 480

	grid := NewGrid(rows, cols, initialWidth, initialHeight, bodyPadding, cellMargin, cellBorderSize, maintainAspectRatio)

	game := &Game{
		grid:     grid,
		fontSize: fontSize,
	}

	// Create the text wrapper
	textWrapper, err := textwrapper.NewTextWrapper(fontpath, fontSize, false)
	if err != nil {
		panic(err)
	}
	game.textWrapper = textWrapper
	game.textWrapper.Color = color.RGBA{80, 255, 80, 255}

	// Set the initial window size
	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("Grid Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
