package charts

import (
	_ "embed"
	"image/color"
	"strconv"

	"example.com/menu/internals/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Chart03 struct {
	Data        []float64
	XLabel      string
	YLabel      string
	NumXTicks   int
	NumYTicks   int
	GutterWidth float32
	ScreenSize  int
	OffsetX     float64
	OffsetY     float64
	BarColor    color.RGBA
	PointColor  color.RGBA
	AxisColor   color.RGBA
	TextWrapper *textwrapper.TextWrapper
}

func (g *Chart03) DrawPlotline(screen *ebiten.Image) {
	screenWidth := float64(screen.Bounds().Dx())
	screenHeight := float64(screen.Bounds().Dy())

	usableWidth := screenWidth - g.OffsetX*2
	usableHeight := screenHeight - g.OffsetY*2

	//green := color.RGBA{0, 255, 0, 255}

	yMin := g.minValue()
	yMax := g.maxValue()
	if yMin > 0 {
		yMin = 0
	}
	yRange := yMax - yMin
	yScale := usableHeight / yRange

	for x := 0.0; x < usableWidth; x++ {

		dataIndex := int(x / usableWidth * float64(len(g.Data)))
		if dataIndex < len(g.Data) {
			y := g.Data[dataIndex]
			scaledY := (y - yMin) * yScale
			vector.DrawFilledRect(
				screen,
				float32(g.OffsetX+x),
				float32(screenHeight-g.OffsetY-scaledY),
				1*g.GutterWidth,
				1,
				g.PointColor,
				false,
			)
		}
	}

	g.drawAxis(screen, g.minValue(), g.maxValue())

}

func (g *Chart03) DrawBars(screen *ebiten.Image) {
	screenWidth := float64(screen.Bounds().Dx())
	screenHeight := float64(screen.Bounds().Dy())
	usableWidth := screenWidth - g.OffsetX*2
	usableHeight := screenHeight - g.OffsetY*2
	min := g.minValue()
	max := g.maxValue()

	if min > 0 {
		min = 0
	}
	adjustedRange := max - min

	barWidth := usableWidth / float64(len(g.Data))
	for i, y := range g.Data {
		scaledY := ((y - min) / adjustedRange) * usableHeight
		vector.DrawFilledRect(screen,
			float32(g.OffsetX+barWidth*float64(i)),
			float32(screenHeight-g.OffsetY-scaledY),
			float32(barWidth*float64(g.GutterWidth)),
			float32(scaledY),
			g.BarColor,
			false,
		)
	}

	g.drawAxis(screen, min, max)

}

func (g *Chart03) drawAxis(screen *ebiten.Image, min, max float64) {
	screenWidth := float64(screen.Bounds().Dx())
	screenHeight := float64(screen.Bounds().Dy())
	usableWidth := screenWidth - g.OffsetX*2
	usableHeight := screenHeight - g.OffsetY*2

	vector.StrokeLine(screen, float32(g.OffsetX), float32(g.OffsetY), float32(g.OffsetX), float32(screenHeight-g.OffsetY), 1, g.AxisColor, false)
	vector.StrokeLine(screen, float32(g.OffsetX), float32(screenHeight-g.OffsetY), float32(screenWidth-g.OffsetX), float32(screenHeight-g.OffsetY), 1, g.AxisColor, false)

	g.drawAxisLabels(screen)

	for i := 0; i <= g.NumXTicks; i++ {
		tickX := g.OffsetX + (usableWidth/float64(g.NumXTicks))*float64(i)
		vector.StrokeLine(screen, float32(tickX), float32(screenHeight-g.OffsetY), float32(tickX), float32(screenHeight-g.OffsetY+5), 1, g.AxisColor, false)
		labelInt := int((float64(i) / float64(g.NumXTicks)) * float64(len(g.Data)))

		g.TextWrapper.DrawText(
			screen,
			strconv.Itoa(labelInt),
			float64(tickX)-10,
			float64(screenHeight-g.OffsetY+20),
		)
	}

	for i := 0; i <= g.NumYTicks; i++ {
		tickY := screenHeight - g.OffsetY - (usableHeight/float64(g.NumYTicks))*float64(i)
		vector.StrokeLine(screen, float32(g.OffsetX), float32(tickY), float32(g.OffsetX-5), float32(tickY), 1, g.AxisColor, false)
		labelInt := int(min + ((max-min)/float64(g.NumYTicks))*float64(i))

		g.TextWrapper.DrawText(
			screen,
			strconv.Itoa(labelInt),
			float64(g.OffsetX-30),
			float64(tickY+5),
		)
	}
}

func (g *Chart03) drawAxisLabels(screen *ebiten.Image) {
	screenWidth := float64(screen.Bounds().Dx())
	screenHeight := float64(screen.Bounds().Dy())

	g.TextWrapper.DrawText(
		screen,
		g.XLabel,
		float64(screenWidth)-20,
		float64(screenHeight-g.OffsetY+15),
	)

	g.TextWrapper.DrawText(
		screen,
		g.YLabel,
		float64(g.OffsetX-20),
		float64(g.OffsetY-20),
	)
}

func (g *Chart03) minValue() float64 {
	min := g.Data[0]
	for _, v := range g.Data {
		if v < min {
			min = v
		}
	}
	return min
}

func (g *Chart03) maxValue() float64 {
	max := g.Data[0]
	for _, v := range g.Data {
		if v > max {
			max = v
		}
	}
	return max
}
