package widgets

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type TextAreaSelection struct {
	text            string
	hasFocus        bool
	cursorPos       int
	counter         int
	selectionStart  int
	selectionEnd    int
	isSelecting     bool
	x, y, w, h      int
	maxLines        int
	cursorBlinkRate int
	tabWidth        int
	lineHeight      int
	font            font.Face
}

func NewTextAreaSelection(x, y, w, h, maxLines int) *TextAreaSelection {
	return &TextAreaSelection{
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		maxLines:        maxLines,
		cursorBlinkRate: 30,
		tabWidth:        4,
		lineHeight:      20,
		font:            basicfont.Face7x13,
	}
}

// Update the Draw method to use the struct fields instead of global variables
func (t *TextAreaSelection) Draw(screen *ebiten.Image) {
	// Draw the textarea background
	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)

	// Draw the text with selection
	lines := strings.Split(t.text, "\n")
	yOffset := t.y // Initial y offset for text drawing

	for i, line := range lines {
		if i >= t.maxLines {
			break
		}

		lineText := line
		lineX := t.x
		lineY := yOffset + t.lineHeight/2 // Center vertically

		// Calculate text width for the current line
		//currentTextWidth := textWidth(basicfont.Face7x13, lineText)

		// Draw the text selection
		if t.selectionStart != t.selectionEnd {
			startLine, startCol := t.getCursorLineAndColForPos(t.selectionStart)
			endLine, endCol := t.getCursorLineAndColForPos(t.selectionEnd)

			if startLine == endLine {
				// Single line selection
				if startCol > len(line) {
					startCol = len(line)
				}
				if endCol > len(line) {
					endCol = len(line)
				}
				startX := t.x + t.textWidth(line[:startCol])
				endX := t.x + t.textWidth(line[:endCol])

				// Draw the selection rectangle for a single line
				vector.DrawFilledRect(screen, float32(startX), float32(yOffset), float32(endX-startX), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
			} else {
				// Multi-line selection

				// Handle the first line
				if startCol > len(line) {
					startCol = len(line)
				}
				startX := t.x + t.textWidth(line[:startCol])
				vector.DrawFilledRect(screen, float32(startX), float32(yOffset), float32(t.x+t.w-startX), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)

				// Handle middle lines
				for j := startLine + 1; j < endLine; j++ {
					if j >= len(lines) {
						break
					}
					vector.DrawFilledRect(screen, float32(t.x), float32(yOffset+t.lineHeight), float32(t.w), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
					yOffset += t.lineHeight
				}

				// Handle the last line
				if endCol > len(line) {
					endCol = len(line)
				}
				endX := t.x + t.textWidth(line[:endCol])
				vector.DrawFilledRect(screen, float32(t.x), float32(yOffset), float32(endX-t.x), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
			}
		}

		// Draw the text itself
		text.Draw(screen, lineText, t.font, lineX, lineY+t.lineHeight/2, color.Black)

		// Move yOffset down for the next line
		yOffset += t.lineHeight
	}

	// Draw the cursor if the text area has focus
	if t.hasFocus {
		cursorLine, cursorCol := t.getCursorLineAndCol()
		cursorX := t.x + t.textWidth(lines[cursorLine][:cursorCol])
		cursorY := t.y + cursorLine*t.lineHeight + t.lineHeight/2

		// Draw the cursor (flashing rectangle)
		if t.counter%(t.cursorBlinkRate*2) < t.cursorBlinkRate {
			vector.DrawFilledRect(screen, float32(cursorX), float32(cursorY-t.lineHeight/2), 2, float32(t.lineHeight), color.RGBA{0, 0, 0, 255}, true)
		}
	}
}

func (t *TextAreaSelection) Update() error {
	// Update focus state based on mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			t.hasFocus = true
			t.startSelectionAtPosition(x, y)
		} else {
			t.hasFocus = false
		}
	}

	// Continue selection if mouse is dragged while pressed
	if t.hasFocus && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		t.isSelecting = true
	}
	if t.isSelecting && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		t.updateSelection(x, y)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		t.isSelecting = false
	}

	// If the text area has focus, handle keyboard input
	if t.hasFocus {
		t.handleKeyboardInput()
	}

	t.counter++
	return nil
}

func (t *TextAreaSelection) handleKeyboardInput() {
	// Handle backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(t.text) > 0 && t.cursorPos > 0 {
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.cursorPos--
		t.clearSelection()
	}

	// Handle Tab key
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if t.isSelecting {
			t.indentSelection()
		} else {
			t.text = t.text[:t.cursorPos] + strings.Repeat(" ", t.tabWidth) + t.text[t.cursorPos:]
			t.cursorPos += t.tabWidth
		}
		t.clearSelection()
	}

	// Handle character input
	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			t.text = t.text[:t.cursorPos] + string(char) + t.text[t.cursorPos:]
			t.cursorPos++
			t.clearSelection()
		}
	}

	// Handle enter key for new line
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
		t.cursorPos++
		t.clearSelection()
	}

	// Handle arrow keys for cursor movement
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && t.cursorPos > 0 {
		t.cursorPos--
		t.updateSelectionWithShiftKey(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && t.cursorPos < len(t.text) {
		t.cursorPos++
		t.updateSelectionWithShiftKey(1)
	}
	// Ensure cursor position is within text bounds
	t.cursorPos = clamp(t.cursorPos, 0, len(t.text))
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (t *TextAreaSelection) indentSelection() {
	lines := strings.Split(t.text, "\n")
	startLine, _ := t.getCursorLineAndColForPos(t.selectionStart)
	endLine, _ := t.getCursorLineAndColForPos(t.selectionEnd)
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	if startLine >= len(lines) {
		startLine = len(lines) - 1
	}

	var indent string
	for i := 0; i < t.tabWidth; i++ {
		indent += " "
	}

	for i := startLine; i <= endLine; i++ {
		lines[i] = indent + lines[i]
	}

	t.text = strings.Join(lines, "\n")
	t.cursorPos = t.selectionEnd + len(indent) // Adjust cursor position after indentation
}

func (t *TextAreaSelection) startSelectionAtPosition(x, y int) {
	line := (y - t.y) / t.lineHeight
	col := x - t.x

	lines := strings.Split(t.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	if line < 0 {
		line = 0
	}

	lineText := lines[line]
	colIndex := 0
	for i := range lineText {
		charWidth := t.textWidth(string(lineText[i]))
		if col < colIndex+charWidth/2 {
			break
		}
		colIndex += charWidth
	}

	charPos := t.getCharPosFromLineAndCol(line, colIndex)
	t.cursorPos = charPos
	t.selectionStart = charPos
	t.selectionEnd = charPos
	t.isSelecting = true
}

func (t *TextAreaSelection) updateSelection(x, y int) {
	line := (y - t.y) / t.lineHeight
	col := x - t.x

	lines := strings.Split(t.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	if line < 0 {
		line = 0
	}

	lineText := lines[line]
	colIndex := 0
	for i := range lineText {
		charWidth := t.textWidth(string(lineText[i]))
		if col < colIndex+charWidth/2 {
			break
		}
		colIndex += charWidth
	}

	charPos := t.getCharPosFromLineAndCol(line, colIndex)
	t.cursorPos = charPos
	t.selectionEnd = charPos
}

func (t *TextAreaSelection) updateSelectionWithShiftKey(offset int) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		newCursorPos := clamp(t.cursorPos+offset, 0, len(t.text))

		if newCursorPos > t.cursorPos {
			// Cursor moving right
			if newCursorPos > t.selectionEnd {
				t.selectionEnd = newCursorPos
			} else {
				t.selectionStart = newCursorPos
			}
		} else {
			// Cursor moving left
			if newCursorPos < t.selectionStart {
				t.selectionStart = newCursorPos
			} else {
				t.selectionEnd = newCursorPos
			}
		}
		t.cursorPos = newCursorPos
	} else {
		// If Shift is not pressed, clear the selection
		t.clearSelection()
	}
}

func (t *TextAreaSelection) clearSelection() {
	t.selectionStart = t.cursorPos
	t.selectionEnd = t.cursorPos
}

func (t *TextAreaSelection) getCharPosFromLineAndCol(line, col int) int {
	lines := strings.Split(t.text, "\n")
	charPos := 0
	for i := 0; i < line; i++ {
		charPos += len(lines[i]) + 1 // +1 for newline character
	}
	charPos += col
	return charPos
}

func (t *TextAreaSelection) getCursorLineAndColForPos(pos int) (int, int) {
	lines := strings.Split(t.text, "\n")
	charCount := 0
	for i, line := range lines {
		if charCount+len(line)+1 > pos { // +1 for newline character
			return i, pos - charCount
		}
		charCount += len(line) + 1 // +1 for newline character
	}
	return len(lines) - 1, len(lines[len(lines)-1])
}

func (t *TextAreaSelection) getCursorLineAndCol() (int, int) {
	return t.getCursorLineAndColForPos(t.cursorPos)
}

func (t *TextAreaSelection) textWidth(str string) int {
	width := 0
	for _, x := range str {
		awidth, _ := t.font.GlyphAdvance(x)
		width += int(awidth >> 6)
	}
	return width
}
