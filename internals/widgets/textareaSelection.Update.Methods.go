package widgets

import (
	"fmt"
	"strings"
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

func (t *TextAreaSelection) selectWordAt(pos int) {
	if len(t.text) == 0 {
		return
	}

	// Clamp pos to the valid range
	pos = clamp(pos, 0, len(t.text)-1)

	// If pos is at a word separator, do not select
	if isWordSeparator(rune(t.text[pos])) {
		return
	}

	start := pos
	// Move start backwards to the beginning of the word
	for start > 0 && !isWordSeparator(rune(t.text[start-1])) {
		start--
	}

	end := pos
	// Move end forwards to the end of the word
	for end < len(t.text) && !isWordSeparator(rune(t.text[end])) {
		end++
	}

	t.setSelectionStart(start)
	t.setSelectionEnd(end)
	t.setCursorPos(end)

	t.SetIsSelecting(false)

	// Debugging statement to verify selection
	fmt.Printf("Selected word from byte %d to byte %d: %q\n", start, end, t.text[start:end])
}

func (t *TextAreaSelection) selectEntireLineAt(x, y int) {
	charPos := t.getCharPosFromPosition(x, y)
	line, _ := t.getCursorLineAndColForPos(charPos)
	lines := t.cachedLines

	if line < 0 || line >= len(lines) {
		return
	}

	// Calculate the start and end positions of the line
	charStart := t.getCharPosFromLineAndCol(line, 0)
	charEnd := t.getCharPosFromLineAndCol(line, len(lines[line]))

	// Set the selection to the entire line
	t.setSelectionStart(charStart)
	t.setSelectionEnd(charEnd)
	t.setCursorPos(charEnd)

	t.SetIsSelecting(false)
}

func (t *TextAreaSelection) getCharPosFromPosition(x, y int) int {
	// Adjust the line calculation by adding the scrollOffset
	line := float64(y-t.y-t.paddingTop)/t.lineHeight + float64(t.scrollOffset)
	col := float64(x - t.x - t.paddingLeft)

	lines := t.cachedLines
	if line >= float64(len(lines)) {
		line = float64(len(lines)) - 1
	}
	if line < 0 {
		line = 0
	}

	lineInt := int(line)
	if lineInt >= len(lines) {
		lineInt = len(lines) - 1
	}
	if lineInt < 0 {
		lineInt = 0
	}

	lineText := lines[lineInt]
	colIndex := 0
	accumulatedWidth := 0.0

	for i, char := range lineText {
		charWidth := t.textWidth(string(char))

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

	// Ensure colIndex does not exceed line length
	if colIndex > len(lineText) {
		colIndex = len(lineText)
	}

	charPos := t.getCharPosFromLineAndCol(lineInt, colIndex)
	fmt.Printf("Mouse click at (x=%d, y=%d) mapped to byte position %d\n", x, y, charPos)
	return charPos
}
