package main

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TextAreaSelection is a simplified text area widget.
type TextAreaSelection struct {
	textWrapper *TextWrapper
	text        string
	cursorPos   int
	counter     int
	hasFocus    bool

	x, y, w, h int
}

// NewTextAreaSelection initializes a new TextAreaSelection.
func NewTextAreaSelection(textWrapper *TextWrapper, x, y, w, h int, startTxt string) *TextAreaSelection {
	return &TextAreaSelection{
		textWrapper: textWrapper,
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		text:        startTxt,
		cursorPos:   0,
		counter:     0,
		hasFocus:    false,
	}
}

// Draw renders the text and the blinking cursor.
func (t *TextAreaSelection) Draw(screen *ebiten.Image) {
	// Draw the background of the text area
	ebitenutil.DrawRect(screen, float64(t.x), float64(t.y), float64(t.w), float64(t.h), color.RGBA{200, 200, 200, 255})

	// Split the text into lines
	lines := strings.Split(t.text, "\n")
	lineHeight := t.textWrapper.FontHeight()

	// Draw each line with appropriate Y-offset
	for i, line := range lines {
		yOffset := float64(t.y) + 15 + float64(i)*lineHeight
		t.textWrapper.DrawText(screen, line, float64(t.x)+5, yOffset)
	}

	// Draw the blinking cursor if focused
	if t.hasFocus {
		// Simple cursor blinking logic
		if (t.counter/30)%2 == 0 {
			cursorX, cursorY := t.getCursorXY()
			ebitenutil.DrawRect(screen, float64(t.x)+cursorX, cursorY-lineHeight+5, 2, lineHeight, color.Black)
		}
	}
}

// getCursorXY returns the cursor's X and Y positions based on cursorPos.
func (t *TextAreaSelection) getCursorXY() (float64, float64) {
	lines := strings.Split(t.text, "\n")
	var currentLine string
	var lineIndex int

	// Determine which line the cursor is on
	for i, line := range lines {
		if t.cursorPos <= len(line) {
			currentLine = line
			lineIndex = i
			break
		}
		t.cursorPos -= len(line) + 1 // +1 for the newline character
	}

	// Measure the width up to the cursor position
	width, _ := t.textWrapper.MeasureString(currentLine[:t.cursorPos])
	xPos := float64(t.x) + 5 + width
	yPos := float64(t.y) + 15 + float64(lineIndex)*t.textWrapper.FontHeight()

	return xPos, yPos
}

// Update handles cursor blinking and mouse click events.
func (t *TextAreaSelection) Update() {
	t.counter++

	// Handle mouse click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			t.hasFocus = true
			clickPos := t.getCursorPosFromClick(x, y)
			t.setCursorPos(clickPos)
		} else {
			t.hasFocus = false
		}
	}
}

// getCursorPosFromClick calculates the cursor position based on the mouse click.
func (t *TextAreaSelection) getCursorPosFromClick(x, y int) int {
	if t.textWrapper == nil {
		return 0
	}

	// Calculate lineHeight and padding
	lineHeight := t.textWrapper.FontHeight()
	paddingY := float64(t.y) + 15

	// Determine which line was clicked
	clickY := float64(y) - paddingY
	if clickY < 0 {
		clickY = 0
	}
	lineIndex := int(clickY / lineHeight)
	lines := strings.Split(t.text, "\n")

	if lineIndex >= len(lines) {
		lineIndex = len(lines) - 1
	}
	if lineIndex < 0 {
		lineIndex = 0
	}

	selectedLine := lines[lineIndex]

	// Calculate the x position within the line
	clickX := float64(x - t.x - 5) // Adjust for padding
	if clickX < 0 {
		clickX = 0
	}

	cursorInLine := 0
	for i := 0; i <= len(selectedLine); i++ {
		substr := selectedLine[:i]
		width, _ := t.textWrapper.MeasureString(substr)
		if width >= clickX {
			cursorInLine = i
			break
		}
		cursorInLine = i
	}

	// Calculate the overall cursor position
	cursorPos := 0
	for i := 0; i < lineIndex; i++ {
		cursorPos += len(lines[i]) + 1 // +1 for newline character
	}
	cursorPos += cursorInLine

	// Ensure cursorPos is within bounds
	if cursorPos > len(t.text) {
		cursorPos = len(t.text)
	}

	return cursorPos
}

// setCursorPos sets the cursor position ensuring it stays within bounds.
func (t *TextAreaSelection) setCursorPos(pos int) {
	if pos < 0 {
		t.cursorPos = 0
	} else if pos > len(t.text) {
		t.cursorPos = len(t.text)
	} else {
		t.cursorPos = pos
	}
}

// Handle inputs (optional for future extensions)
func (t *TextAreaSelection) HandleInputs() {
	// Placeholder for handling keyboard inputs if needed
}
