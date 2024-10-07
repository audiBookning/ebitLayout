package widgets

import (
	"image/color"
	"log"

	"example.com/menu/internals/textwrapper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ButtonStd struct {
	X, Y            float32
	Width, Height   float32
	Text            string
	FontColor       color.Color
	BackgroundColor color.Color
	TextWrapper     *textwrapper.TextWrapper
	FontSize        float64
	OnClick         func()
}

func NewButtonStd(
	x, y,
	width, height float32,
	text string,
	textWrapper *textwrapper.TextWrapper,
	fontColor color.Color,
	backgroundColor color.Color,
	fontSize float64,
	onClick func(),
) *ButtonStd {
	return &ButtonStd{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Text:            text,
		FontColor:       fontColor,
		TextWrapper:     textWrapper,
		FontSize:        fontSize,
		BackgroundColor: backgroundColor,
		OnClick:         onClick,
	}
}

func (b *ButtonStd) Draw(screen *ebiten.Image) {
	// Draw the button background
	vector.DrawFilledRect(screen, b.X, b.Y, b.Width, b.Height, b.BackgroundColor, false)

	// Draw the button text
	b.TextWrapper.Color = b.FontColor
	b.TextWrapper.SetFontSize(b.FontSize)

	// Measure the text dimensions
	textWidth, textHeight := b.TextWrapper.MeasureText(b.Text)

	// Calculate position to center the text within the button
	textX := float64(b.X) + float64(b.Width)/2 - textWidth/2
	textY := float64(b.Y) + float64(b.Height)/2 - textHeight/2

	b.TextWrapper.DrawText(
		screen,
		b.Text,
		textX,
		textY,
	)
}

func (b *ButtonStd) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) {
	if isAnimating {
		return // Ignore clicks during animations
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		globalX := float32(x)
		globalY := float32(y)

		// Adjust global cursor position to local by subtracting navigator offset
		localX := globalX - navigatorOffsetX
		localY := globalY - navigatorOffsetY

		// Debug: Log cursor and button positions
		log.Printf("Cursor Position: (%f, %f)", globalX, globalY)
		log.Printf("Button Position: (%f, %f) with size (%f, %f)", b.X, b.Y, b.Width, b.Height)
		log.Printf("Local Cursor Position: (%f, %f)", localX, localY)

		// Check if the cursor is within the button's bounds
		chek01 := localX >= b.X
		chek02 := localX < b.X+b.Width
		chek03 := localY >= b.Y
		chek04 := localY < b.Y+b.Height

		if chek01 && chek02 && chek03 && chek04 {
			if b.OnClick != nil {
				log.Printf("Button '%s' clicked.", b.Text)
				b.OnClick()
			}
		}
	}
}
