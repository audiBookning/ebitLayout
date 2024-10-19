package responsive

import (
	"example.com/menu/cmd02/more06/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
)

type Title struct {
	Text        string
	X, Y        int
	FontScale   float64
	TextWrapper *textwrapper.TextWrapper
}

func NewTitle(text string, tw *textwrapper.TextWrapper) *Title {
	return &Title{
		Text:        text,
		FontScale:   1.0,
		TextWrapper: tw,
	}
}

func (t *Title) Draw(screen *ebiten.Image) {
	t.TextWrapper.DrawText(screen, t.Text, float64(t.X), float64(t.Y))
}
