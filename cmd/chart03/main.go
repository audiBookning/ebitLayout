package main

import (
	_ "embed"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/charts"
	"example.com/menu/internals/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	barGraph      *charts.Chart03
	plotlineGraph *charts.Chart03
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Create sub-images for each graph
	barImage := ebiten.NewImage(screen.Bounds().Dx()/2, screen.Bounds().Dy())
	plotlineImage := ebiten.NewImage(screen.Bounds().Dx()/2, screen.Bounds().Dy())

	// Draw each graph on its respective sub-image
	g.barGraph.DrawBars(barImage)
	g.plotlineGraph.DrawPlotline(plotlineImage)

	// Draw the sub-images on the main screen
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(barImage, op)

	op.GeoM.Translate(float64(screen.Bounds().Dx()/2), 0)
	screen.DrawImage(plotlineImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func getFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

var filePathTxt string

const Assets_Relative_Path = "../../"

func main() {
	// Generate some sine wave sinData for demonstration
	_, filePathTxt, _, _ = runtime.Caller(0)

	ScreenSize := 800
	fontpath := getFilePath("assets/fonts/roboto_regularTTF.ttf")

	textWrapper, err := textwrapper.NewTextWrapper(
		fontpath,
		12,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	textWrapper.Color = color.RGBA{255, 255, 255, 255} // white
	textWrapper.SetFontSize(12)

	// Create the bar graph
	randomData := []float64{4.2, 7.5, 3.8, 6.1, 9.4}
	barGraph := &charts.Chart03{
		Data:        randomData,
		XLabel:      "Categories",
		YLabel:      "Values",
		NumXTicks:   5,
		NumYTicks:   10,
		GutterWidth: 0.8,
		OffsetX:     50,
		OffsetY:     50,
		BarColor:    color.RGBA{0, 0, 255, 255},
		PointColor:  color.RGBA{0, 255, 0, 255},
		AxisColor:   color.RGBA{255, 255, 255, 255},
		ScreenSize:  ScreenSize / 2,
		TextWrapper: textWrapper,
	}

	// Create the plotline graph
	sinData := make([]float64, 100)
	for i := 0; i < len(sinData); i++ {
		normalizedX := float64(i) / float64(len(sinData)) * 2 * math.Pi
		sinData[i] = math.Sin(normalizedX) * 100
	}

	plotlineGraph := &charts.Chart03{
		Data:        sinData,
		XLabel:      "X",
		YLabel:      "Y",
		NumXTicks:   10,
		NumYTicks:   10,
		GutterWidth: 0.8,
		OffsetX:     50,
		OffsetY:     50,
		BarColor:    color.RGBA{0, 0, 255, 255},
		PointColor:  color.RGBA{0, 255, 0, 255},
		AxisColor:   color.RGBA{255, 255, 255, 255},
		ScreenSize:  ScreenSize / 2,
		TextWrapper: textWrapper,
	}

	// Create the game object and run it
	game := &Game{barGraph: barGraph, plotlineGraph: plotlineGraph}
	ebiten.SetWindowSize(ScreenSize, ScreenSize)
	ebiten.SetWindowTitle("Bar and Plotline Graphs")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
