package charts

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

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
		AxesColor:    color.RGBA{255, 255, 255, 255},
		PlotColor:    color.RGBA{0, 255, 0, 255},
	}
}

func (g *Chart01) DrawAxes(screen *ebiten.Image) {

	vector.StrokeLine(screen,
		float32(g.OffsetX),
		float32(g.OffsetY),
		float32(g.OffsetX),
		float32(g.ScreenHeight-g.OffsetY),
		1,
		g.AxesColor,
		false,
	)

	vector.StrokeLine(screen,
		float32(g.OffsetX),
		float32(g.ScreenHeight-g.OffsetY),
		float32(g.ScreenWidth-g.OffsetX),
		float32(g.ScreenHeight-g.OffsetY),
		1,
		g.AxesColor,
		false,
	)

	text.Draw(screen, "X", g.Face, int(g.ScreenWidth)-40, int(g.ScreenHeight-g.OffsetY+15), g.AxesColor)
	text.Draw(screen, "Y", g.Face, int(g.OffsetX-20), int(g.OffsetY-10), g.AxesColor)

	text.Draw(screen, "(0,0)", g.Face, int(g.OffsetX-25), int(g.ScreenHeight-g.OffsetY+15), g.AxesColor)
}

func (g *Chart01) PlotSineWave(screen *ebiten.Image) {
	for x := 0.0; x < g.UsableWidth; x++ {

		normalizedX := x / g.UsableWidth * 2 * math.Pi
		y := math.Sin(normalizedX) * 100

		scaledY := (y - g.YMin) * g.YScale

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
