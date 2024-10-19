package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Title struct {
	Text      string
	X, Y      int
	FontScale float64
}

func NewTitle(text string) *Title {
	return &Title{
		Text:      text,
		FontScale: 1.0,
	}
}

func (t *Title) Draw(screen *ebiten.Image) {

	text.Draw(screen, t.Text, basicfont.Face7x13, t.X, t.Y, color.White)
}
