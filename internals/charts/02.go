package charts

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// Chart02 encapsulates all graph-related functionalities.
type Chart02 struct {
	Screen       *ebiten.Image
	OffsetX      float64
	OffsetY      float64
	UsableWidth  float64
	UsableHeight float64
	White        color.RGBA
	Green        color.RGBA
	Face         *basicfont.Face
	NumXTicks    int
	NumYTicks    int
	YMin         float64
	YMax         float64
	YScale       float64
}

// NewChart02 initializes and returns a new Graph instance.
func NewChart02(screen *ebiten.Image) *Chart02 {
	screenWidth := float64(screen.Bounds().Dx())
	screenHeight := float64(screen.Bounds().Dy())

	offsetX := 50.0
	offsetY := 50.0
	usableWidth := screenWidth - offsetX*2
	usableHeight := screenHeight - offsetY*2
	yMin, yMax := -100.0, 100.0
	yRange := yMax - yMin
	yScale := usableHeight / yRange

	face := basicfont.Face7x13

	return &Chart02{
		Screen:       screen,
		OffsetX:      offsetX,
		OffsetY:      offsetY,
		UsableWidth:  usableWidth,
		UsableHeight: usableHeight,
		White:        color.RGBA{255, 255, 255, 255},
		Green:        color.RGBA{0, 255, 0, 255},
		Face:         face,
		NumXTicks:    10,
		NumYTicks:    10,
		YMin:         yMin,
		YMax:         yMax,
		YScale:       yScale,
	}
}

// DrawAxes draws the X and Y axes.
func (g *Chart02) DrawAxes() {
	screenWidth := float64(g.Screen.Bounds().Dx())
	screenHeight := float64(g.Screen.Bounds().Dy())

	// Y Axis
	vector.StrokeLine(g.Screen, float32(g.OffsetX), float32(g.OffsetY), float32(g.OffsetX), float32(screenHeight-g.OffsetY), 1, g.White, false)
	// X Axis
	vector.StrokeLine(g.Screen, float32(g.OffsetX), float32(screenHeight-g.OffsetY), float32(screenWidth-g.OffsetX), float32(screenHeight-g.OffsetY), 1, g.White, false)
}

// PlotSineWave plots a sine wave on the graph.
func (g *Chart02) PlotSineWave() {
	for x := 0.0; x < g.UsableWidth; x++ {
		normalizedX := x / g.UsableWidth * 2 * math.Pi
		y := math.Sin(normalizedX) * 100
		scaledY := (y - g.YMin) * g.YScale
		vector.DrawFilledRect(g.Screen, float32(g.OffsetX+x), float32(float64(g.Screen.Bounds().Dy())-g.OffsetY-scaledY), 1, 1, g.Green, false)
	}
}

// DrawAxisLabels draws the labels for X and Y axes.
func (g *Chart02) DrawAxisLabels() {
	screenWidth := float64(g.Screen.Bounds().Dx())
	screenHeight := float64(g.Screen.Bounds().Dy())

	// X-axis label
	text.Draw(g.Screen, "X", g.Face, int(screenWidth)-20, int(screenHeight-g.OffsetY+15), g.White)
	// Y-axis label
	text.Draw(g.Screen, "Y", g.Face, int(g.OffsetX-20), int(g.OffsetY-20), g.White)
}

// DrawTicks draws the ticks and their corresponding labels on both axes.
func (g *Chart02) DrawTicks() {
	screenHeight := float64(g.Screen.Bounds().Dy())

	// X-Axis Ticks
	for i := 0; i <= g.NumXTicks; i++ {
		tickX := g.OffsetX + (g.UsableWidth/float64(g.NumXTicks))*float64(i)
		vector.StrokeLine(g.Screen, float32(tickX), float32(screenHeight-g.OffsetY), float32(tickX), float32(screenHeight-g.OffsetY+5), 1, g.White, false)
		label := int((float64(i) / float64(g.NumXTicks)) * 2 * math.Pi * 100)
		text.Draw(g.Screen, strconv.Itoa(label), g.Face, int(tickX)-10, int(screenHeight-g.OffsetY+20), g.White)
	}

	// Y-Axis Ticks
	for i := 0; i <= g.NumYTicks; i++ {
		tickY := screenHeight - g.OffsetY - (g.UsableHeight/float64(g.NumYTicks))*float64(i)
		vector.StrokeLine(g.Screen, float32(g.OffsetX), float32(tickY), float32(g.OffsetX-5), float32(tickY), 1, g.White, false)
		label := int(g.YMin + (float64(g.YMax-g.YMin)/float64(g.NumYTicks))*float64(i))
		text.Draw(g.Screen, strconv.Itoa(label), g.Face, int(g.OffsetX-30), int(tickY+5), g.White)
	}
}

// Render draws the entire graph.
func (g *Chart02) Render() {
	g.DrawAxes()
	g.PlotSineWave()
	g.DrawAxisLabels()
	g.DrawTicks()
}
