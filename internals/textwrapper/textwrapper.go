package textwrapper

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
)

type TextWrapper struct {
	fontFaceSource   *text.GoTextFaceSource
	GoTextFace       *text.GoTextFace
	Color            color.Color
	fontSize         float64
	Position         image.Point
	isVertical       bool
	textOptions      *text.DrawOptions
	textOptionsDirty bool
	WordWrap         bool
	Font             font.Face
	MaxWidth         int // Optional: Maximum width for text wrapping
	MaxLines         int // Optional: Maximum number of lines for text wrapping
	LineHeight       int // Height of each line
}

func NewTextWrapper(fontPath string, fontSize float64, isVertical bool) (*TextWrapper, error) {
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
	textOptions := &text.DrawOptions{}

	td := &TextWrapper{
		fontFaceSource:   fontFaceSource,
		GoTextFace:       fontFace,
		isVertical:       isVertical,
		Color:            color.White,
		fontSize:         fontSize,
		Position:         image.Point{X: 0, Y: 0},
		textOptions:      textOptions,
		textOptionsDirty: false,
	}

	return td, nil
}

func (tw *TextWrapper) SetFontSize(fontSize float64) {
	tw.GoTextFace.Size = fontSize
	tw.fontSize = fontSize
}

func (tw *TextWrapper) SetGeomScale(x float64, y float64) {
	tw.textOptions.GeoM.Scale(x, y)
	tw.textOptionsDirty = true
}

func (tw *TextWrapper) ResetGeom() {
	tw.textOptions.GeoM.Reset()
}

func (tw *TextWrapper) DrawTextWithWordWrap(screen *ebiten.Image, str string, x, y int) {
	if tw.WordWrap && tw.MaxWidth > 0 {
		lines := tw.wrapText(str)
		for i, line := range lines {
			if tw.MaxLines > 0 && i >= tw.MaxLines {
				break
			}

			tw.DrawText(screen, line, float64(x), float64(y+tw.LineHeight*(i+1)))
		}
	} else {
		tw.DrawText(screen, str, float64(x), float64(y))
	}
}

func (tw *TextWrapper) wrapText(str string) []string {
	words := splitIntoWords(str)
	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word

		if tw.MeasureTextWidth(testLine) > tw.MaxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func splitIntoWords(str string) []string {
	var words []string
	currentWord := ""
	for _, r := range str {
		if r == ' ' || r == '\n' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
			if r == '\n' {
				words = append(words, "\n")
			}
		} else {
			currentWord += string(r)
		}
	}
	if currentWord != "" {
		words = append(words, currentWord)
	}
	return words
}

func (tw *TextWrapper) MeasureTextWidth(str string) int {
	width, _ := text.Measure(str, tw.GoTextFace, 0)
	return int(width)
}

func (tw *TextWrapper) MeasureText(s string) (float64, float64) {
	metrics := tw.GoTextFace.Metrics()
	var lineSpacing float64
	if tw.isVertical {

		lineSpacing = metrics.VAscent + metrics.VDescent + metrics.VLineGap
	} else {

		lineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap
	}
	width, height := text.Measure(s, tw.GoTextFace, lineSpacing)
	return width, height
}

// Add this method to the TextWrapper struct
func (tw *TextWrapper) MeasureString(s string) (float64, float64) {
	return text.Measure(s, tw.GoTextFace, 0)
}

func (tw *TextWrapper) GetFontMetrics() text.Metrics {
	return tw.GoTextFace.Metrics()
}

func (tw *TextWrapper) GetMonospaceWidth() float64 {
	size := tw.GoTextFace.Size
	return size
}

func (tw *TextWrapper) SetColor(color color.Color) {
	tw.Color = color
}

func (tw *TextWrapper) GetTextFace() *text.GoTextFace {
	return tw.GoTextFace
}

func (tw *TextWrapper) DrawText(screen *ebiten.Image, textStr string, x, y float64) {
	//tw.textOptions.GeoM.Reset()
	if !tw.textOptionsDirty {
		tw.textOptions.GeoM.Reset()
	}
	tw.textOptionsDirty = false
	tw.textOptions.GeoM.Translate(x, y)
	tw.textOptions.ColorScale = ebiten.ColorScale{}
	tw.textOptions.ColorScale.ScaleWithColor(tw.Color)
	text.Draw(screen, textStr,
		tw.GoTextFace,
		tw.textOptions)
}
