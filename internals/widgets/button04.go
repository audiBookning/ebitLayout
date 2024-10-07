package widgets

import (
	"image"
	"image/color"
	"math"

	"example.com/menu/internals/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const animationDuration = 15 // Number of frames for the animation

type ToggleButton04 struct {
	X, Y                 int
	Width, Height        int
	OnLabel              string
	OffLabel             string
	DefaultColor         color.Color
	CurrentColor         color.Color
	ToggleColor          color.Color
	IsToggled            bool
	knobX                float64
	tx                   *textwrapper.TextWrapper
	animationProgress    int
	cachedOnLabelBounds  image.Rectangle
	cachedOffLabelBounds image.Rectangle
	OnClickFunc          func()
}

func (b *ToggleButton04) OnMouseDown() {
	// not implemented
	b.IsToggled = !b.IsToggled
	if b.IsToggled {
		b.CurrentColor = b.ToggleColor
	} else {
		b.CurrentColor = b.DefaultColor
	}
}

func (b *ToggleButton04) SetHovered(isHovered bool) {
	// not implemented
}

func NewToggleButton04(
	x, y,
	width, height int,
	onLabel, offLabel string,
	defaultColor, toggleColor color.Color,
	tx *textwrapper.TextWrapper,
	onClick func(),
) *ToggleButton04 {
	b := &ToggleButton04{
		X: x, Y: y,
		Width: width, Height: height,
		OnLabel: onLabel, OffLabel: offLabel,
		DefaultColor: defaultColor, ToggleColor: toggleColor,
		tx:          tx,
		OnClickFunc: onClick,
	}
	b.knobX = float64(x)
	onWidth, onHeight := b.tx.MeasureText(onLabel)
	b.cachedOnLabelBounds = image.Rect(0, 0, int(onWidth), int(onHeight))

	offWidth, offHeight := b.tx.MeasureText(offLabel)
	b.cachedOffLabelBounds = image.Rect(0, 0, int(offWidth), int(offHeight))
	return b
}

func (b *ToggleButton04) OnClick() {
	b.IsToggled = !b.IsToggled
	b.animationProgress = 0 // Reset animation progress
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
}

func (b *ToggleButton04) Update() {
	if b.animationProgress < animationDuration {
		b.animationProgress++
	}

	startX := float64(b.X)
	endX := float64(b.X + b.Width - b.Height)

	progress := float64(b.animationProgress) / animationDuration

	// Use a smooth easing function
	easedProgress := easeInOutCubic(progress)

	if !b.IsToggled {
		// Swap start and end for the off state
		startX, endX = endX, startX
	}

	b.knobX = startX + (endX-startX)*easedProgress
}

// easeInOutCubic provides a smooth acceleration and deceleration
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

func (b *ToggleButton04) Draw(screen *ebiten.Image) {
	// Draw background
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), color.RGBA{200, 200, 200, 255}, true)

	// Draw knob
	knobSize := float32(b.Height)
	vector.DrawFilledRect(
		screen,
		float32(b.knobX), float32(b.Y),
		knobSize, knobSize,
		b.CurrentColor, true,
	)

	// Draw label
	label := b.OffLabel
	bounds := b.cachedOffLabelBounds
	if b.IsToggled {
		label = b.OnLabel
		bounds = b.cachedOnLabelBounds
	}

	// Adjust text position
	x := b.X + b.Width + 10             // 10 pixels padding to the right of the button
	y := b.Y + (b.Height-bounds.Dy())/2 // Center vertically

	b.tx.Position = image.Point{X: x, Y: y}
	b.tx.DrawText(screen, label, float64(x), float64(y))

	// Debug: Draw a red dot at the text position
	vector.DrawFilledCircle(screen, float32(x), float32(y), 2, color.RGBA{255, 0, 0, 255}, true)

}

func (b *ToggleButton04) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}
