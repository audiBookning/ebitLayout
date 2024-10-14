package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// TextWrapper is responsible for rendering text on the screen.
type TextWrapper struct {
	fontFace font.Face
	Color    color.Color
}

// NewTextWrapper initializes a new TextWrapper with the specified font.
func NewTextWrapper(fontPath string, fontSize float64) (*TextWrapper, error) {
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	tt, err := truetype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	fontFace := truetype.NewFace(tt, &truetype.Options{
		Size: fontSize,
	})

	return &TextWrapper{
		fontFace: fontFace,
		Color:    color.Black,
	}, nil
}

// DrawText renders the given string at the specified (x, y) position on the screen.
func (tw *TextWrapper) DrawText(screen *ebiten.Image, textStr string, x, y float64) {
	text.Draw(screen, textStr, tw.fontFace, int(x), int(y), tw.Color)
}

// MeasureString returns the width and height of the given string.
func (tw *TextWrapper) MeasureString(s string) (float64, float64) {
	if len(s) == 0 {
		return 0, tw.FontHeight()
	}

	var widthFixed fixed.Int26_6
	for _, r := range s {
		advance, ok := tw.fontFace.GlyphAdvance(r)
		if !ok {
			// Handle missing glyphs if necessary
			continue
		}
		widthFixed += advance
	}

	width := float64(widthFixed) / 64.0 // Convert from fixed.Int26_6 to float64
	height := tw.FontHeight()

	return width, height
}

// FontHeight returns the height of the font.
func (tw *TextWrapper) FontHeight() float64 {
	return float64(tw.fontFace.Metrics().Height) / 64.0
}
