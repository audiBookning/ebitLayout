package layout

import (
	"image/color"

	"example.com/menu/internals/textwrapper"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Element struct {
	X, Y, Width, Height int
	Color               color.Color
	Flex                int
	Text                string
	TextWrapper         *textwrapper.TextWrapper
	TextSize            float64
}

func (e *Element) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(e.X), float32(e.Y),
		float32(e.Width), float32(e.Height),
		e.Color,
		true,
	)

	if e.TextWrapper != nil {

		e.TextWrapper.SetFontSize(e.TextSize)
		textWidth, textHeight := e.TextWrapper.MeasureText(e.Text)

		centerX := float64(e.X) + float64(e.Width)/2
		centerY := float64(e.Y) + float64(e.Height)/2

		textX := centerX - textWidth/2
		textY := centerY - textHeight/2

		e.TextWrapper.DrawText(screen, e.Text, textX, textY)
	}
}

func (e *Element) Update() error {

	return nil
}

func (e *Element) Layout(x, y, width, height int) {
	e.X = x
	e.Y = y
	e.Width = width
	e.Height = height
}

type FlexBox struct {
	Elements       []Element
	Direction      string
	JustifyContent string
	AlignItems     string
}

func (fb *FlexBox) Layout(width, height int) {
	totalFlex := 0
	totalSize := 0

	for _, element := range fb.Elements {
		totalFlex += element.Flex
		if fb.Direction == "row" {
			totalSize += element.Width
		} else {
			totalSize += element.Height
		}
	}

	var offset int
	if fb.Direction == "row" {

		availableSpace := width - totalSize
		for i := range fb.Elements {
			if fb.Elements[i].Flex > 0 {

				flexSize := (availableSpace * fb.Elements[i].Flex) / totalFlex
				fb.Elements[i].Width += flexSize
			}
		}

		offset = 0
		for i := range fb.Elements {
			fb.Elements[i].Layout(offset, (height-fb.Elements[i].Height)/2, fb.Elements[i].Width, fb.Elements[i].Height)
			offset += fb.Elements[i].Width
		}
	} else {

		availableSpace := height - totalSize
		for i := range fb.Elements {
			if fb.Elements[i].Flex > 0 {
				flexSize := (availableSpace * fb.Elements[i].Flex) / totalFlex
				fb.Elements[i].Height += flexSize
			}
		}

		offset = 0
		for i := range fb.Elements {
			fb.Elements[i].Layout((width-fb.Elements[i].Width)/2, offset, fb.Elements[i].Width, fb.Elements[i].Height)
			offset += fb.Elements[i].Height
		}
	}
}

func (fb *FlexBox) Draw(screen *ebiten.Image) {
	for _, element := range fb.Elements {
		element.Draw(screen)
	}
}
