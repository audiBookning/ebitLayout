package widgets

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

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
	t.isTextChanged = true // Add this line
}

func (t *TextAreaSelection) clearSelection() {
	t.setSelectionStart(t.cursorPos)
	t.setSelectionEnd(t.cursorPos)
	//fmt.Printf("Selection Cleared: Start=%d, End=%d, CursorPos=%d\n", t.selectionStart, t.selectionEnd, t.cursorPos)
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

func (t *TextAreaSelection) updateSelection(newPos int) {
	t.setSelectionEnd(newPos)
	t.setCursorPos(newPos)
	t.updateSelectionBounds()

	// Debugging statement to verify selection updates
	minPos, maxPos := t.getSelectionBounds()
	fmt.Printf("Selection Updated: Start=%d, End=%d, CursorPos=%d\n", minPos, maxPos, t.cursorPos)
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
