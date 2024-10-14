package textwrapper02

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// TextWrapper is a reusable component for rendering and measuring text.
type TextWrapper struct {
	Font       font.Face
	Color      color.Color
	FontSize   int
	MaxWidth   int // Optional: Maximum width for text wrapping
	WordWrap   bool
	MaxLines   int // Optional: Maximum number of lines for text wrapping
	LineHeight int // Height of each line
}

// NewTextWrapper creates a new instance of TextWrapper.
func NewTextWrapper(fontFace font.Face, fontSize int, clr color.Color) *TextWrapper {
	return &TextWrapper{
		Font:       fontFace,
		Color:      clr,
		FontSize:   fontSize,
		WordWrap:   false,
		MaxWidth:   0,
		MaxLines:   0,
		LineHeight: fontFace.Metrics().Height.Ceil(),
	}
}

// SetWordWrap enables or disables word wrapping with a specified maximum width and lines.
func (tw *TextWrapper) SetWordWrap(maxWidth, maxLines int) {
	tw.WordWrap = true
	tw.MaxWidth = maxWidth
	tw.MaxLines = maxLines
}

// DrawText renders the specified text at the given (x, y) position on the screen.
func (tw *TextWrapper) DrawText(screen *ebiten.Image, str string, x, y int) {
	text.Draw(screen, str, tw.Font, x, y, tw.Color)
}
func (tw *TextWrapper) DrawTextWithWordWrap(screen *ebiten.Image, str string, x, y int) {
	if tw.WordWrap && tw.MaxWidth > 0 {
		lines := tw.wrapText(str)
		for i, line := range lines {
			if tw.MaxLines > 0 && i >= tw.MaxLines {
				break
			}
			text.Draw(screen, line, tw.Font, x, y+tw.LineHeight*i, tw.Color)
		}
	} else {
		text.Draw(screen, str, tw.Font, x, y, tw.Color)
	}
}

// MeasureTextWidth returns the width of the given text string.
func (tw *TextWrapper) MeasureTextWidth(str string) int {
	return text.BoundString(tw.Font, str).Dx()
}

// MeasureTextHeight returns the total height of the given text string, considering word wrapping.
func (tw *TextWrapper) MeasureTextHeight(str string) int {
	if tw.WordWrap && tw.MaxWidth > 0 {
		lines := tw.wrapText(str)
		return tw.LineHeight * len(lines)
	}
	// Single-line height
	return tw.LineHeight
}

// wrapText splits the text into lines based on the maximum width.
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

// splitIntoWords splits a string into words separated by spaces.
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

// splitIntoLines splits a string into lines separated by newline characters.
func splitIntoLines(str string) []string {
	var lines []string
	currentLine := ""
	for _, r := range str {
		if r == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(r)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}
