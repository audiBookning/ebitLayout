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
			cursorX, _ := t.getCursorXPosition()
			currentLine := t.getCurrentLine()
			yOffset := float64(t.y) + 15 + float64(currentLine)*lineHeight
			ebitenutil.DrawRect(screen, float64(t.x)+cursorX+5, yOffset-5, 2, lineHeight, color.Black)
		}
	}
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
	// Calculate the line number based on click Y position
	relativeY := y - t.y - 15
	lineHeight := t.textWrapper.FontHeight()
	lineNum := relativeY / int(lineHeight)

	// Clamp the line number to the number of lines
	lines := strings.Split(t.text, "\n")
	if lineNum < 0 {
		lineNum = 0
	} else if lineNum >= len(lines) {
		lineNum = len(lines) - 1
	}

	// Get the specific line
	line := lines[lineNum]

	// Adjust x for padding and calculate cursor position within the line
	relativeX := float64(x - t.x - 5)
	cursorPosInLine := t.getCursorPosInLine(line, relativeX)

	// Calculate the absolute cursor position
	absolutePos := 0
	for i := 0; i < lineNum; i++ {
		absolutePos += len(lines[i]) + 1 // +1 for the newline character
	}
	absolutePos += cursorPosInLine

	return absolutePos
}

// getCursorPosInLine calculates the cursor position within a single line based on X coordinate.
func (t *TextAreaSelection) getCursorPosInLine(line string, clickX float64) int {
	currentPos := 0
	for i := 0; i <= len(line); i++ {
		substr := line[:i]
		width, _ := t.textWrapper.MeasureString(substr)
		if width >= clickX {
			return i
		}
		currentPos = i
	}
	return currentPos
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

// getCurrentLine determines the current line based on cursor position.
func (t *TextAreaSelection) getCurrentLine() int {
	lines := strings.Split(t.text[:t.cursorPos], "\n")
	return len(lines) - 1
}

// getCursorXPosition calculates the X position of the cursor based on the current cursor position.
func (t *TextAreaSelection) getCursorXPosition() (float64, float64) {
	// Find the current line text
	lines := strings.Split(t.text, "\n")
	currentLine := t.getCurrentLine()
	if currentLine >= len(lines) {
		currentLine = len(lines) - 1
	}
	lineText := lines[currentLine]

	// Get the text up to the cursor position within the line
	cursorInLine := t.cursorPos
	for i := 0; i < currentLine; i++ {
		cursorInLine -= len(lines[i]) + 1 // +1 for newline
	}
	if cursorInLine < 0 {
		cursorInLine = 0
	}
	textUpToCursor := ""
	if cursorInLine >= 0 && cursorInLine <= len(lineText) {
		textUpToCursor = lineText[:cursorInLine]
	}

	// Measure the width of the text up to the cursor
	cursorX, _ := t.textWrapper.MeasureString(textUpToCursor)

	return cursorX, float64(t.y)
}

// HandleInputs handles keyboard inputs for text editing (optional for future extensions)
func (t *TextAreaSelection) HandleInputs() {
	// Placeholder for handling keyboard inputs if needed
}
