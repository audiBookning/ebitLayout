package charts

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// Chart01 encapsulates graph-related properties and methods.
type Chart01 struct {
	OffsetX      float64
	OffsetY      float64
	UsableWidth  float64
	UsableHeight float64
	ScreenWidth  float64
	ScreenHeight float64
	YMin         float64
	YMax         float64
	YScale       float64
	Face         *basicfont.Face
	AxesColor    color.Color
	PlotColor    color.Color
}

// NewChart01 initializes and returns a new Graph instance.
func NewChart01(screenWidth, screenHeight float64) *Chart01 {
	offsetX := 50.0
	offsetY := 50.0
	usableWidth := screenWidth - offsetX*2
	usableHeight := screenHeight - offsetY*2
	yMin, yMax := -100.0, 100.0
	yScale := usableHeight / (yMax - yMin)

	return &Chart01{
		OffsetX:      offsetX,
		OffsetY:      offsetY,
		UsableWidth:  usableWidth,
		UsableHeight: usableHeight,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		YMin:         yMin,
		YMax:         yMax,
		YScale:       yScale,
		Face:         basicfont.Face7x13,
		AxesColor:    color.RGBA{255, 255, 255, 255}, // White
		PlotColor:    color.RGBA{0, 255, 0, 255},     // Green
	}
}

// DrawAxes draws the X and Y axes on the screen.
func (g *Chart01) DrawAxes(screen *ebiten.Image) {
	// Draw Y axis
	vector.StrokeLine(screen,
		float32(g.OffsetX),
		float32(g.OffsetY),
		float32(g.OffsetX),
		float32(g.ScreenHeight-g.OffsetY),
		1,
		g.AxesColor,
		false,
	)

	// Draw X axis
	vector.StrokeLine(screen,
		float32(g.OffsetX),
		float32(g.ScreenHeight-g.OffsetY),
		float32(g.ScreenWidth-g.OffsetX),
		float32(g.ScreenHeight-g.OffsetY),
		1,
		g.AxesColor,
		false,
	)

	// Draw axis labels
	text.Draw(screen, "X", g.Face, int(g.ScreenWidth)-40, int(g.ScreenHeight-g.OffsetY+15), g.AxesColor)
	text.Draw(screen, "Y", g.Face, int(g.OffsetX-20), int(g.OffsetY-10), g.AxesColor)

	// Draw origin label
	text.Draw(screen, "(0,0)", g.Face, int(g.OffsetX-25), int(g.ScreenHeight-g.OffsetY+15), g.AxesColor)
}

// PlotSineWave plots a sine wave on the screen.
func (g *Chart01) PlotSineWave(screen *ebiten.Image) {
	for x := 0.0; x < g.UsableWidth; x++ {
		// Compute the sine value
		normalizedX := x / g.UsableWidth * 2 * math.Pi // Normalize X to [0, 2π]
		y := math.Sin(normalizedX) * 100               // Sine wave

		// Adjust Y for scaling and plotting
		scaledY := (y - g.YMin) * g.YScale // Scale Y to fit graph height

		// Plot the point
		vector.DrawFilledRect(screen,
			float32(g.OffsetX+x),
			float32(g.ScreenHeight-g.OffsetY-scaledY),
			1,
			1,
			g.PlotColor,
			false,
		)
	}
}
