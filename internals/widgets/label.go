package widgets

import (
	"image"
	"image/color"

	"example.com/menu/internals/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
)

type Label struct {
	Text        string
	X, Y        float64
	FontSize    int
	FontColor   color.Color
	Align       string
	textWrapper *textwrapper.TextWrapper
}

func NewLabel(
	textWrapper *textwrapper.TextWrapper,
	text string,
	x, y float64,
	fontSize int,
	fontColor color.Color,
	align string,
) (*Label, error) {

	return &Label{
		textWrapper: textWrapper,
		Text:        text,
		X:           x,
		Y:           y,
		FontSize:    fontSize,
		FontColor:   fontColor,
		Align:       align,
	}, nil
}

func (l *Label) Update() error {

	return nil
}

func (l *Label) Draw(screen *ebiten.Image) {
	l.textWrapper.Color = l.FontColor
	l.textWrapper.SetFontSize(float64(l.FontSize))

	textWidth, _ := l.textWrapper.MeasureText(l.Text)
	x := l.X

	switch l.Align {
	case "center":
		x -= float64(textWidth) / 2
	case "right":
		x -= float64(textWidth)
	}

	l.textWrapper.Position = image.Point{X: int(x), Y: int(l.Y)}
	l.textWrapper.DrawText(screen, l.Text, x, l.Y)
}

func (l *Label) Layout(outsideWidth, outsideHeight int) (int, int) {

	return outsideWidth, outsideHeight
}
