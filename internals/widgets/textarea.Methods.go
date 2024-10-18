package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (t *TextArea) isCtrlPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)
}

func (t *TextArea) isShiftPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
}

func (t *TextArea) getCharPosFromLineAndColWithclamp(line, col int) int {
	lines := t.cachedLines
	charPos := 0
	for i := 0; i < line; i++ {
		charPos += len(lines[i]) + 1
	}
	charPos += col
	charPos = clamp(charPos, 0, len(t.text))
	return charPos
}

func (t *TextArea) getCursorLineAndColForPos(pos int) (int, int) {
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

func (t *TextArea) textWidth(str string) float64 {
	width, _ := t.textWrapper.MeasureString(str)
	return width
}

func (t *TextArea) moveToWordStart(pos int) int {
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

func (t *TextArea) moveToWordEnd(pos int) int {
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

// ---------------------

func (t *TextArea) pushUndo() {
	state := TextState{
		Text:      t.text,
		CursorPos: t.cursorPos,
	}
	t.undoStack = append(t.undoStack, state)
	// Clear redoStack whenever a new action is made
	t.redoStack = []TextState{}

}
