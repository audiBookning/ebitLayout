package main

import (
	"bytes"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TextWrapper is responsible for rendering text on the screen.
type TextWrapper struct {
	fontFace       *text.GoTextFace
	Color          color.Color
	textOptions    *text.DrawOptions
	textOptionsSet bool
}

// NewTextWrapper initializes a new TextWrapper with the specified font.
func NewTextWrapper(fontPath string, fontSize float64) (*TextWrapper, error) {
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	fontFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		return nil, fmt.Errorf("failed to create font face source: %w", err)
	}

	fontFace := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   fontSize,
	}

	return &TextWrapper{
		fontFace:    fontFace,
		Color:       color.Black,
		textOptions: &text.DrawOptions{},
	}, nil
}

// DrawText renders the given string at the specified (x, y) position on the screen.
func (tw *TextWrapper) DrawText(screen *ebiten.Image, textStr string, x, y float64) {
	tw.textOptions.GeoM.Reset()
	tw.textOptions.GeoM.Translate(x, y)
	tw.textOptions.ColorScale = ebiten.ColorScale{}
	tw.textOptions.ColorScale.ScaleWithColor(tw.Color)
	text.Draw(screen, textStr, tw.fontFace, tw.textOptions)
}

// MeasureString returns the width and height of the given string.
func (tw *TextWrapper) MeasureString(s string) (float64, float64) {
	width, height := text.Measure(s, tw.fontFace, 0)
	return float64(width), float64(height)
}

// FontHeight returns the height of the font using a sample character.
func (tw *TextWrapper) FontHeight() float64 {
	_, height := tw.MeasureString("A")
	return height
}
