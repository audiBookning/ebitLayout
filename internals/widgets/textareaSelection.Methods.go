package widgets

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (t *TextAreaSelection) isOverScrollbar(x, y int) bool {
	return x >= t.scrollbarX && x <= t.scrollbarX+t.scrollbarWidth && y >= t.y && y <= t.y+t.h
}

func (t *TextAreaSelection) dragScrollbar(mouseY int) {
	// Calculate the new thumb Y position relative to the scrollbar's track
	newThumbY := float64(mouseY-t.scrollbarY) - t.dragOffsetY

	// Clamp the thumb position within the scrollbar track
	maxThumbY := float64(t.scrollbarHeight - int(t.scrollbarThumbH))
	if newThumbY < 0 {
		newThumbY = 0
	}
	if newThumbY > maxThumbY {
		newThumbY = maxThumbY
	}

	t.scrollbarThumbY = newThumbY

	// Calculate the corresponding scrollOffset
	totalLines := len(strings.Split(t.text, "\n"))
	maxScrollOffset := totalLines - t.maxLines
	if maxScrollOffset < 1 {
		maxScrollOffset = 1
	}

	t.SetScrollOffset(int((t.scrollbarThumbY / maxThumbY) * float64(maxScrollOffset)))
	t.SetScrollOffset(clamp(t.scrollOffset, 0, maxScrollOffset))

	//fmt.Printf("Dragging Scrollbar: ThumbY=%.2f, ScrollOffset=%d\n", t.scrollbarThumbY, t.scrollOffset)
}

func (t *TextAreaSelection) drawScrollbar(screen *ebiten.Image, totalLines int) {
	// Initialize scrollbar properties
	t.scrollbarWidth = 10
	t.scrollbarX = t.x + t.w - t.scrollbarWidth
	t.scrollbarY = t.y
	t.scrollbarHeight = t.h

	// Calculate the height of the scrollbar thumb
	visibleRatio := float64(t.maxLines) / float64(totalLines)
	t.scrollbarThumbH = float64(t.scrollbarHeight) * visibleRatio
	if t.scrollbarThumbH < 20 {
		t.scrollbarThumbH = 20 // Minimum thumb height
	}

	// Calculate the Y position of the scrollbar thumb relative to the scrollbar track
	maxScrollOffset := totalLines - t.maxLines
	if maxScrollOffset < 1 {
		maxScrollOffset = 1
	}
	thumbMaxY := float64(t.scrollbarHeight - int(t.scrollbarThumbH))
	t.scrollbarThumbY = float64(t.scrollOffset) / float64(maxScrollOffset) * thumbMaxY

	// Draw the scrollbar track
	vector.DrawFilledRect(screen, float32(t.scrollbarX), float32(t.scrollbarY), float32(t.scrollbarWidth), float32(t.scrollbarHeight), color.RGBA{220, 220, 220, 255}, true)

	// Draw the scrollbar thumb
	vector.DrawFilledRect(
		screen,
		float32(t.scrollbarX),
		float32(t.scrollbarY)+float32(t.scrollbarThumbY),
		float32(t.scrollbarWidth),
		float32(t.scrollbarThumbH),
		color.RGBA{160, 160, 160, 255},
		true)
}

func (t *TextAreaSelection) selectWordAt(pos int) {
	if len(t.text) == 0 {
		return
	}

	runes := []rune(t.text)
	textLen := len(runes)

	pos = clamp(pos, 0, textLen-1)

	if isWordSeparator(runes[pos]) {
		return
	}

	start := pos
	for start > 0 && !isWordSeparator(runes[start-1]) {
		start--
	}

	end := pos
	for end < textLen && !isWordSeparator(runes[end]) {
		end++
	}

	byteStart := runePosToBytePos(t.text, start)
	byteEnd := runePosToBytePos(t.text, end)

	t.setSelectionStart(byteStart)
	t.setSelectionEnd(byteEnd)
	t.setCursorPos(byteEnd)

	t.SetIsSelecting(false)
	/* completeWord := t.text[byteStart:byteEnd]
	fmt.Printf("Word Selected=%s | pos=%d | Byte Start=%d, Byte End=%d \n", completeWord, pos, start, end) */
}

func (t *TextAreaSelection) updateSelectionWithShiftKey(offset int) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		newCursorPos := clamp(t.cursorPos+offset, 0, len(t.text))
		_, currentCol := t.getCursorLineAndColForPos(t.cursorPos)

		if t.desiredCursorCol == -1 {
			t.desiredCursorCol = currentCol
		}

		if offset < 0 { // Moving left
			if newCursorPos > t.cursorPos {
				if newCursorPos > t.selectionEnd {
					t.selectionEnd = newCursorPos
				} else {
					t.setSelectionStart(newCursorPos)
				}
			} else {
				if newCursorPos < t.selectionStart {
					t.setSelectionStart(newCursorPos)
				} else {
					t.setSelectionEnd(newCursorPos)
				}
			}
		} else { // Moving right
			if newCursorPos > t.cursorPos {
				if newCursorPos > t.selectionEnd {
					t.setSelectionEnd(newCursorPos)
				} else {
					t.setSelectionStart(newCursorPos)
				}
			} else {
				if newCursorPos < t.selectionStart {
					t.setSelectionStart(newCursorPos)
				} else {
					t.setSelectionEnd(newCursorPos)
				}
			}
		}

		t.setCursorPos(newCursorPos)
	} else {
		t.clearSelection()
	}

	//fmt.Printf("Selection Updated: Start=%d, End=%d, CursorPos=%d\n", t.selectionStart, t.selectionEnd, t.cursorPos)
}

func (t *TextAreaSelection) isCtrlPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)
}

// getSelectionBoundsStart returns the start position of the current selection
func (t *TextAreaSelection) getSelectionBoundsStart() int {
	minPos, _ := t.getSelectionBounds()
	return minPos
}

// getSelectionBoundsEnd returns the end position of the current selection
func (t *TextAreaSelection) getSelectionBoundsEnd() int {
	_, maxPos := t.getSelectionBounds()
	return maxPos
}

func (t *TextAreaSelection) indentSelection() {
	lines := t.cachedLines
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
	t.setCursorPos(t.selectionEnd + len(indent))
}

func (t *TextAreaSelection) clearSelection() {
	t.setSelectionStart(t.cursorPos)
	t.setSelectionEnd(t.cursorPos)
	//fmt.Printf("Selection Cleared: Start=%d, End=%d, CursorPos=%d\n", t.selectionStart, t.selectionEnd, t.cursorPos)
}

func (t *TextAreaSelection) getCharPosFromLineAndCol(line, col int) int {
	lines := t.cachedLines
	charPos := 0
	for i := 0; i < line; i++ {
		charPos += len(lines[i]) + 1
	}
	charPos += col
	return charPos
}

func (t *TextAreaSelection) getCursorLineAndColForPos(pos int) (int, int) {
	lines := t.cachedLines
	charCount := 0
	for i, line := range lines {
		if charCount+len(line)+1 > pos {
			return i, pos - charCount
		}
		charCount += len(line) + 1
	}
	return len(lines) - 1, len(lines[len(lines)-1])
}

func (t *TextAreaSelection) getCursorLineAndCol() (int, int) {
	return t.getCursorLineAndColForPos(t.cursorPos)
}

func (t *TextAreaSelection) textWidth(str string) float64 {
	//width, _ := t.textWrapper.MeasureText(str)
	width, _ := t.textWrapper.MeasureString(str)
	//width := t.textWrapper.MeasureTextWidth(str)
	return width
}

func (t *TextAreaSelection) getCharPosFromPosition(x, y int) int {
	// Adjust the line calculation by adding the scrollOffset
	line := float64(y-t.y-t.paddingTop)/t.lineHeight + float64(t.scrollOffset)
	//line := (y-t.y)/int(t.lineHeight) + t.scrollOffset
	col := float64(x - t.x - t.paddingLeft)
	//col := x - t.x
	/*
		fmt.Println("X      : ", x, "        Y: ", y)
		fmt.Println("Xpadded: ", t.x+t.paddingLeft, "  YPadded: ", t.y+t.paddingTop)
		fmt.Println("Line   : ", line, "      Col: ", col)
	*/
	lines := t.cachedLines
	if line >= float64(len(lines)) {
		line = float64(len(lines)) - 1
	}
	if line < 0 {
		line = 0
	}

	lineint := int(line)
	lineText := lines[lineint]
	colIndex := 0
	accumulatedWidth := 0.0

	//fmt.Print(" :  char       charWidth     accumulatedWidth     col      colIndex\n")
	for i, char := range lineText {
		charString := string(char)
		charWidth := t.textWidth(charString)

		//fmt.Print(i, ":  ", charString, "             ", charWidth, "                 ", accumulatedWidth, "             ", col, "          ", colIndex, "\n")
		// Check if the click is within the current character's width
		if col < accumulatedWidth+charWidth {
			// Determine if the click is in the first half or second half of the character
			if col < accumulatedWidth+(charWidth/2) {
				colIndex = i
			} else {
				colIndex = i + 1
			}
			break
		}
		accumulatedWidth += charWidth
		colIndex = i + 1
	}
	//fmt.Println("line = ", line, ", colIndex = ", colIndex)
	//fmt.Println("----------------------------------------")

	// Ensure colIndex does not exceed line length
	if colIndex > len(lineText) {
		colIndex = len(lineText)
	}

	return t.getCharPosFromLineAndCol(lineint, colIndex)
}

func (t *TextAreaSelection) getSelectionBounds() (int, int) {
	t.updateSelectionBounds()
	return t.minSelectionPos, t.maxSelectionPos
}

// deleteSelection removes the currently selected text and updates the cursor position
func (t *TextAreaSelection) deleteSelection() {
	minPos, maxPos := t.getSelectionBounds()
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.setCursorPos(minPos)
	t.clearSelection()
}

func (t *TextAreaSelection) moveToWordStart(pos int) int {
	if pos == 0 {
		return pos
	}

	// If cursor is already at the start of a word, move to the previous word
	if pos > 0 && t.text[pos-1] != ' ' && (pos == 0 || t.text[pos-1] == ' ') {
		for pos > 0 && t.text[pos-1] != ' ' && t.text[pos-1] != '\n' {
			pos--
		}
		return pos
	}

	// Otherwise, move to the start of the current word
	for pos > 0 && t.text[pos-1] == ' ' {
		pos--
	}
	for pos > 0 && t.text[pos-1] != ' ' && t.text[pos-1] != '\n' {
		pos--
	}
	return pos
}

func (t *TextAreaSelection) moveToWordEnd(pos int) int {
	textLen := len(t.text)
	if pos >= textLen {
		return pos
	}

	// If cursor is at the end of a word, move to the end of the next word
	if pos < textLen && t.text[pos] != ' ' && t.text[pos] != '\n' {
		for pos < textLen && t.text[pos] != ' ' && t.text[pos] != '\n' {
			pos++
		}
		return pos
	}

	// Otherwise, move to the end of the current word
	for pos < textLen && t.text[pos] == ' ' {
		pos++
	}
	for pos < textLen && t.text[pos] != ' ' && t.text[pos] != '\n' {
		pos++
	}
	return pos
}

func (t *TextAreaSelection) isShiftPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
}

func (t *TextAreaSelection) updateSelection(newPos int) {

	t.setSelectionEnd(newPos)
	t.setCursorPos(newPos)
}

// ---------------------

func (t *TextAreaSelection) pushUndo() {
	state := TextState{
		Text:      t.text,
		CursorPos: t.cursorPos,
	}
	t.undoStack = append(t.undoStack, state)
	// Clear redoStack whenever a new action is made
	t.redoStack = []TextState{}

}

func (t *TextAreaSelection) updateSelectionBounds() {
	if t.selectionStart <= t.selectionEnd {
		t.minSelectionPos = t.selectionStart
		t.maxSelectionPos = t.selectionEnd
	} else {
		t.minSelectionPos = t.selectionEnd
		t.maxSelectionPos = t.selectionStart
	}
}
