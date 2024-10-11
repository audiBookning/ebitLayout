package widgets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Slider struct {
	X, Y          float64
	Width, Height float64
	HandlePos     float64
	Dragging      bool
}

func (s *Slider) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		mxf, myf := float64(mx), float64(my)

		if s.Dragging || (mxf >= s.X && mxf <= s.X+s.Width && myf >= s.Y && myf <= s.Y+s.Height) {

			s.Dragging = true

			s.HandlePos = mxf - s.X
			if s.HandlePos < 0 {
				s.HandlePos = 0
			}
			if s.HandlePos > s.Width {
				s.HandlePos = s.Width
			}
		}
	} else {

		s.Dragging = false
	}
}

func (s *Slider) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(s.X), float32(s.Y), float32(s.Width), float32(s.Height), color.RGBA{200, 200, 200, 255}, true)

	handleWidth := 10.0
	handleHeight := s.Height
	handleX := s.X + s.HandlePos - handleWidth/2
	handleY := s.Y
	vector.DrawFilledRect(screen, float32(handleX), float32(handleY), float32(handleWidth), float32(handleHeight), color.RGBA{100, 100, 100, 255}, true)
}
