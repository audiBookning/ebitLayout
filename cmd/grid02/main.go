package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/layout"
	"example.com/menu/internals/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	grid                *layout.Grid
	textWrapper         *textwrapper.TextWrapper
	fontSize            float64
	selectedCell        *layout.Cell
	maintainAspectRatio bool
}

var (
	filePathTxt          string
	Assets_Relative_Path = "../../"
)

func GetFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.maintainAspectRatio = !g.maintainAspectRatio
		g.grid.MaintainAspectRatio = g.maintainAspectRatio
		g.grid.UpdateCellPositions()
		fmt.Printf("Maintain Aspect Ratio: %v\n", g.maintainAspectRatio)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.grid.AddRow()
		fmt.Printf("Added a row. Total Rows: %d\n", g.grid.Rows)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.grid.AddColumn()
		fmt.Printf("Added a column. Total Columns: %d\n", g.grid.Cols)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.HandleClick(x, y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{30, 30, 30, 255})

	borderColor := color.RGBA{255, 255, 255, 255}
	selectedBorderColor := color.RGBA{255, 215, 0, 255}
	//textColor := g.textWrapper.Color

	for _, row := range g.grid.Cells {
		for _, cell := range row {

			x := float32(cell.X) + float32(g.grid.CellMargin)/2
			y := float32(cell.Y) + float32(g.grid.CellMargin)/2
			w := float32(cell.Width) - float32(g.grid.CellMargin)
			h := float32(cell.Height) - float32(g.grid.CellMargin)

			backgroundColor := color.RGBA{50, 50, 50, 255}
			vector.DrawFilledRect(screen, x, y, w, h, backgroundColor, true)

			borderSize := float32(g.grid.CellBorderSize)
			currentBorderColor := borderColor
			if g.selectedCell != nil && g.selectedCell.X == cell.X && g.selectedCell.Y == cell.Y {
				currentBorderColor = selectedBorderColor
			}
			vector.StrokeRect(screen, x+borderSize/2, y+borderSize/2, w-borderSize, h-borderSize, borderSize, currentBorderColor, true)

			text := fmt.Sprintf("(%d, %d)", cell.X/g.grid.InitialCellWidth, cell.Y/g.grid.InitialCellHeight)
			textWidth, textHeight := g.textWrapper.MeasureText(text)

			textX := float64(cell.X) + float64(cell.Width)/2 - float64(textWidth)/2
			textY := float64(cell.Y) + float64(cell.Height)/2 - float64(textHeight)/2

			g.textWrapper.DrawText(screen, text, textX, textY)
		}
	}

	instructionText := "Press 'A' to Toggle Aspect Ratio | 'R' to Add Row | 'C' to Add Column\nClick on cells to select them."
	g.textWrapper.SetFontSize(16)
	g.textWrapper.Color = color.RGBA{255, 255, 255, 255}
	g.textWrapper.DrawText(screen, instructionText, 10, float64(g.grid.TotalHeight)-50)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.grid.TotalWidth = outsideWidth
	g.grid.TotalHeight = outsideHeight
	g.grid.InitializeCells(outsideWidth, outsideHeight)

	initialWidth := g.grid.InitialTotalWidth
	initialHeight := g.grid.InitialTotalHeight

	scaleX := float64(outsideWidth) / float64(initialWidth)
	scaleY := float64(outsideHeight) / float64(initialHeight)

	var scale float64
	if g.grid.MaintainAspectRatio {
		scale = math.Min(scaleX, scaleY)
	} else {
		scale = scaleX
	}

	g.textWrapper.SetFontSize(g.fontSize * scale)

	return outsideWidth, outsideHeight
}

func (g *Game) HandleClick(x, y int) {
	for _, row := range g.grid.Cells {
		for _, cell := range row {
			if x >= cell.X && x < cell.X+cell.Width && y >= cell.Y && y < cell.Y+cell.Height {
				g.selectedCell = &cell
				fmt.Printf("Selected Cell at X: %d, Y: %d\n", cell.X, cell.Y)
				return
			}
		}
	}

	g.selectedCell = nil
}

func InitializeAssets() {
	_, filePathTxt, _, _ = runtime.Caller(0)
}

func main() {
	InitializeAssets()
	fontPath := GetFilePath("assets/fonts/roboto_regularTTF.ttf")

	initialRows := 5
	initialCols := 5
	initialWidth := 800
	initialHeight := 600
	bodyPadding := 20
	cellMargin := 4
	cellBorderSize := 2
	maintainAspectRatio := true
	fontSize := 24.0

	grid := layout.NewGrid(initialRows, initialCols, initialWidth, initialHeight, bodyPadding, cellMargin, cellBorderSize, maintainAspectRatio)

	textWrapper, err := textwrapper.NewTextWrapper(fontPath, fontSize, false)
	if err != nil {
		log.Fatalf("Failed to create TextWrapper: %v", err)
	}
	textWrapper.Color = color.RGBA{255, 255, 255, 255}

	game := &Game{
		grid:                grid,
		textWrapper:         textWrapper,
		fontSize:            fontSize,
		maintainAspectRatio: maintainAspectRatio,
		selectedCell:        nil,
	}

	ebiten.SetWindowSize(initialWidth, initialHeight)
	ebiten.SetWindowTitle("Comprehensive Grid Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
