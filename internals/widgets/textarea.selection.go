package widgets

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func (t *TextArea) updateSelectionWithShiftKey(offset int) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		newCursorPos := clamp(t.cursorPos+offset, 0, len(t.text))
		_, currentCol := t.getCursorLineAndColForPos(t.cursorPos)

		if t.desiredCursorCol == -1 {
			t.desiredCursorCol = currentCol
		}

		if offset < 0 { // Moving left
			if newCursorPos > t.cursorPos {
				if newCursorPos > t.selection.selectionEnd {
					t.selection.selectionEnd = newCursorPos
				} else {
					t.selection.setSelectionStart(newCursorPos)
				}
			} else {
				if newCursorPos < t.selection.selectionStart {
					t.selection.setSelectionStart(newCursorPos)
				} else {
					t.selection.setSelectionEnd(newCursorPos)
				}
			}
		} else { // Moving right
			if newCursorPos > t.cursorPos {
				if newCursorPos > t.selection.selectionEnd {
					t.selection.setSelectionEnd(newCursorPos)
				} else {
					t.selection.setSelectionStart(newCursorPos)
				}
			} else {
				if newCursorPos < t.selection.selectionStart {
					t.selection.setSelectionStart(newCursorPos)
				} else {
					t.selection.setSelectionEnd(newCursorPos)
				}
			}
		}

		t.setCursorPos(newCursorPos)
	} else {
		t.selection.ClearSelection(t.cursorPos)
	}

	//fmt.Printf("Selection Updated: Start=%d, End=%d, CursorPos=%d\n", t.selectionStart, t.selectionEnd, t.cursorPos)
}

// getSelectionBoundsStart returns the start position of the current selection

func (t *TextArea) indentSelection() {
	lines := t.cachedLines
	startLine, _ := t.getCursorLineAndColForPos(t.selection.selectionStart)
	endLine, _ := t.getCursorLineAndColForPos(t.selection.selectionEnd)
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
	t.setCursorPos(t.selection.selectionEnd + len(indent))
	t.isTextChanged = true // Add this line
}

// deleteSelection removes the currently selected text and updates the cursor position
func (t *TextArea) deleteSelection() {
	minPos, maxPos := t.selection.getSelectionBounds()
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.setCursorPos(minPos)
	t.selection.ClearSelection(t.cursorPos)
}

func (t *TextArea) updateSelection(newPos int) {
	t.selection.setSelectionEnd(newPos)
	t.setCursorPos(newPos)
	t.selection.updateSelectionBounds()

	// Debugging statement to verify selection updates
	minPos, maxPos := t.selection.getSelectionBounds()
	fmt.Printf("Selection Updated: Start=%d, End=%d, CursorPos=%d\n", minPos, maxPos, t.cursorPos)
}
