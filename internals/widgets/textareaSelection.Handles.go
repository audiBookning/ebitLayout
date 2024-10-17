package widgets

import (
	"fmt"
	"strings"

	"golang.design/x/clipboard"
)

func (t *TextAreaSelection) handlePageDown() {
	t.pushUndo()
	lines := strings.Split(t.text, "\n")
	totalLines := len(lines)
	// Calculate the new scroll offset
	newScrollOffset := t.scrollOffset + t.maxLines
	if newScrollOffset > totalLines-t.maxLines {
		newScrollOffset = totalLines - t.maxLines
	}
	if newScrollOffset < 0 {
		newScrollOffset = 0
	}
	// Update the scroll offset
	t.SetScrollOffset(newScrollOffset)

}

func (t *TextAreaSelection) handlePageUp() {
	t.pushUndo()
	// Calculate the new scroll offset
	newScrollOffset := t.scrollOffset - t.maxLines
	if newScrollOffset < 0 {
		newScrollOffset = 0
	}
	// Update the scroll offset
	t.SetScrollOffset(newScrollOffset)

}

func (t *TextAreaSelection) handleCtrlShiftLeftArrow() {
	t.pushUndo()
	newPos := t.moveToWordStart(t.cursorPos)
	t.updateSelection(newPos)
}

func (t *TextAreaSelection) handleCtrlShiftRightArrow() {
	t.pushUndo()
	newPos := t.moveToWordEnd(t.cursorPos)
	t.updateSelection(newPos)
}

func (t *TextAreaSelection) handleCtrlShiftUpArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndCol(currentLine, 0)
	t.updateSelection(newPos)
}

func (t *TextAreaSelection) handleCtrlShiftDownArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines) {
		newPos := t.getCharPosFromLineAndCol(currentLine, len(lines[currentLine]))
		t.updateSelection(newPos)
	}
}

func (t *TextAreaSelection) handleCtrlUpArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndCol(currentLine, 0)
	if t.isShiftPressed() {
		t.updateSelection(newPos)
	} else {
		t.clearSelection()
		t.setCursorPos(newPos)
	}
}

func (t *TextAreaSelection) handleCtrlDownArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines) {
		newPos := t.getCharPosFromLineAndCol(currentLine, len(lines[currentLine]))
		if t.isShiftPressed() {
			t.updateSelection(newPos)
		} else {
			t.clearSelection()
			t.setCursorPos(newPos)
		}
	}
}

func (t *TextAreaSelection) handleCtrlShiftHome() {
	t.pushUndo()
	// Select from cursor to beginning of text
	t.setSelectionStart(0)
	t.setSelectionEnd(t.selectionStart) // Use the existing selection start as the end
	t.setCursorPos(0)
	// Scroll to the top of the textarea
	t.SetScrollOffset(0)
}

func (t *TextAreaSelection) handleCtrlShiftEnd() {
	t.pushUndo()
	// Select from cursor to end of text
	// TODO: not necessary, but just for clarity of the logic
	t.setSelectionStart(t.selectionStart) // Use the existing selection start as the start
	t.setSelectionEnd(len(t.text))
	t.setCursorPos(len(t.text))
	// Scroll to the bottom of the textarea
	maxScrollOffset := len(strings.Split(t.text, "\n")) - t.maxLines
	if maxScrollOffset > 0 {
		t.SetScrollOffset(maxScrollOffset)
	}
}

func (t *TextAreaSelection) handleCtrlHome() {
	t.pushUndo()
	// Move cursor to the very beginning of the text
	t.setCursorPos(0)
	if t.isShiftPressed() {
		t.setSelectionEnd(0)
	} else {
		t.clearSelection()
	}
	// Scroll to the top of the textarea
	t.SetScrollOffset(0)
}

func (t *TextAreaSelection) handleCtrlEnd() {
	t.pushUndo()
	// Move cursor to the very end of the text
	t.setCursorPos(len(t.text))
	if t.isShiftPressed() {
		t.setSelectionEnd(len(t.text))
	} else {
		t.clearSelection()
	}
	// Scroll to the bottom of the textarea
	maxScrollOffset := len(strings.Split(t.text, "\n")) - t.maxLines
	if maxScrollOffset > 0 {
		t.SetScrollOffset(maxScrollOffset)
	}
}

func (t *TextAreaSelection) handleCopySelection() {
	if t.selectionStart == t.selectionEnd {
		// No selection to copy
		fmt.Println("handleCopySelection - No selection to copy.")
		return
	}
	minPos, maxPos := t.getSelectionBounds()
	selectedText := t.text[minPos:maxPos]
	fmt.Printf("handleCopySelection - Copying text from %d to %d: %q\n", minPos, maxPos, selectedText)
	// Write to clipboard using golang-design/clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(selectedText))
	if err != nil {
		fmt.Println("handleCopySelection - Error writing to clipboard:", err)
	} else {
		fmt.Println("handleCopySelection - Successfully copied to clipboard.")
	}
}

// handleCutSelection copies the selected text to the OS clipboard and removes it from the text area
func (t *TextAreaSelection) handleCutSelection() {
	if t.selectionStart == t.selectionEnd {
		// No selection to cut
		fmt.Println("handleCutSelection - No selection to cut.")
		return
	}
	t.pushUndo()
	minPos, maxPos := t.getSelectionBounds()
	selectedText := t.text[minPos:maxPos]
	fmt.Printf("handleCutSelection - Cutting text from %d to %d: %q\n", minPos, maxPos, selectedText)
	// Write to clipboard using golang-design/clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(selectedText))
	if err != nil {
		fmt.Println("handleCutSelection - Error writing to clipboard:", err)
	} else {
		fmt.Println("handleCutSelection - Successfully cut to clipboard.")
	}
	// Remove the selected text from the text area
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.setCursorPos(minPos)
	t.clearSelection()
	t.isTextChanged = true
}

// handlePasteClipboard inserts text from the OS clipboard into the text area at the current cursor position
func (t *TextAreaSelection) handlePasteClipboard() {
	t.pushUndo() // Ensure each paste is undoable individually
	clipboardBytes := clipboard.Read(clipboard.FmtText)
	clipboardText := string(clipboardBytes)
	if t.selectionStart != t.selectionEnd {
		// Replace selected text with clipboard text
		minPos, maxPos := t.getSelectionBounds()
		t.text = t.text[:minPos] + clipboardText + t.text[maxPos:]
		t.cursorPos = minPos + len(clipboardText)
		t.isTextChanged = true
	} else {
		// Insert clipboard text at cursor position
		t.text = t.text[:t.cursorPos] + clipboardText + t.text[t.cursorPos:]
		t.cursorPos += len(clipboardText)
		t.isTextChanged = true
	}
	t.clearSelection()
}

func (t *TextAreaSelection) handleBackspace() {
	if t.selectionStart != t.selectionEnd {
		t.pushUndo()
		t.deleteSelection()
		t.isTextChanged = true
	} else if t.cursorPos > 0 {
		t.pushUndo()
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.cursorPos--
		t.isTextChanged = true
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleDelete() {
	if t.selectionStart != t.selectionEnd {
		t.pushUndo()
		t.deleteSelection()
		t.isTextChanged = true
	} else if t.cursorPos < len(t.text) {
		t.pushUndo()
		t.text = t.text[:t.cursorPos] + t.text[t.cursorPos+1:]
		t.isTextChanged = true
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleCtrlBackspace() {
	if t.selectionStart != t.selectionEnd {
		// If there's an active selection, delete the selected text
		t.pushUndo()
		t.deleteSelection()
		t.isTextChanged = true
	} else {
		// Delete from the cursor to the beginning of the word
		t.pushUndo()
		newPos := t.moveToWordStart(t.cursorPos)
		t.text = t.text[:newPos] + t.text[t.cursorPos:]
		t.setCursorPos(newPos)
		t.isTextChanged = true
		t.clearSelection()
	}
}

// handleCtrlDelete deletes text from the current cursor position to the end of the word
func (t *TextAreaSelection) handleCtrlDelete() {
	if t.selectionStart != t.selectionEnd {
		// If there's an active selection, delete the selected text
		t.pushUndo()
		t.deleteSelection()
		t.isTextChanged = true
	} else {
		// Delete from the cursor to the end of the word
		t.pushUndo()
		newPos := t.moveToWordEnd(t.cursorPos)
		// Prevent deleting the newline if cursor is at the end of a line
		if newPos > t.cursorPos && t.text[newPos-1] == '\n' {
			newPos--
		}
		t.text = t.text[:t.cursorPos] + t.text[newPos:]
		t.isTextChanged = true
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleTab() {
	t.pushUndo()
	if t.isSelecting {
		t.indentSelection()
	} else {
		t.text = t.text[:t.cursorPos] + strings.Repeat(" ", t.tabWidth) + t.text[t.cursorPos:]
		t.cursorPos += t.tabWidth
	}
	t.clearSelection()
}

func (t *TextAreaSelection) handleEnter() {
	t.pushUndo()
	t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
	t.cursorPos++
	t.clearSelection()
}

func (t *TextAreaSelection) handleLeftArrow() {
	t.pushUndo()
	if t.selectionStart != t.selectionEnd {
		t.clearSelection()
	} else {
		if t.cursorPos > 0 {
			t.cursorPos--
			t.clearSelection()
		}
	}
}

func (t *TextAreaSelection) handleShiftLeftArrow() {
	t.pushUndo()
	t.updateSelectionWithShiftKey(-1)
}

func (t *TextAreaSelection) handleShiftRightArrow() {
	t.updateSelectionWithShiftKey(1)
}

func (t *TextAreaSelection) handleRightArrow() {
	t.pushUndo()
	if t.selectionStart != t.selectionEnd {
		t.clearSelection()
	} else {
		if t.cursorPos < len(t.text) {
			t.cursorPos++
			t.clearSelection()
		}
	}
}

func (t *TextAreaSelection) handleSelectAll() {
	t.pushUndo()
	t.setSelectionStart(0)
	t.setSelectionEnd(len(t.text))
	t.setCursorPos(len(t.text))
}

func (t *TextAreaSelection) handleHome() {
	t.pushUndo()
	line, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndCol(line, 0)
	t.clearSelection()
	t.setCursorPos(newPos)

}

func (t *TextAreaSelection) handleEnd() {
	t.pushUndo()
	line, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	newPos := t.getCharPosFromLineAndCol(line, len(lines[line]))

	t.clearSelection()
	t.setCursorPos(newPos)
}

func (t *TextAreaSelection) handleCtrlLeftArrow() {
	t.pushUndo()
	newPos := t.moveToWordStart(t.cursorPos)
	t.clearSelection()
	t.setCursorPos(newPos)
}

func (t *TextAreaSelection) handleCtrlRightArrow() {
	t.pushUndo()
	newPos := t.moveToWordEnd(t.cursorPos)

	t.clearSelection()
	t.setCursorPos(newPos)
}

// ---------------------
func (t *TextAreaSelection) handleUpArrow() {
	t.pushUndo()
	currentLine, currentCol := t.getCursorLineAndColForPos(t.cursorPos)
	if currentLine > 0 {
		targetLine := currentLine - 1
		lines := strings.Split(t.text, "\n")
		targetCol := currentCol
		if targetCol > len(lines[targetLine]) {
			targetCol = len(lines[targetLine])
		}
		newPos := t.getCharPosFromLineAndCol(targetLine, targetCol)
		t.setCursorPos(newPos)
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleDownArrow() {
	t.pushUndo()
	currentLine, currentCol := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines)-1 {
		targetLine := currentLine + 1
		targetCol := currentCol
		if targetCol > len(lines[targetLine]) {
			targetCol = len(lines[targetLine])
		}
		newPos := t.getCharPosFromLineAndCol(targetLine, targetCol)
		t.setCursorPos(newPos)
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleShiftUp() {
	t.pushUndo()
	currentLine, currentCol := t.getCursorLineAndColForPos(t.cursorPos)
	if currentLine > 0 {
		targetLine := currentLine - 1
		lines := strings.Split(t.text, "\n")
		desiredCol := t.desiredCursorCol
		if desiredCol == -1 {
			desiredCol = currentCol
			t.desiredCursorCol = desiredCol
		}
		if desiredCol > len(lines[targetLine]) {
			desiredCol = len(lines[targetLine])
		}
		newPos := t.getCharPosFromLineAndCol(targetLine, desiredCol)
		t.updateSelection(newPos)
		t.desiredCursorCol = -1
	}
}

func (t *TextAreaSelection) handleShiftDown() {
	t.pushUndo()
	currentLine, currentCol := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines)-1 {
		targetLine := currentLine + 1
		desiredCol := t.desiredCursorCol
		if desiredCol == -1 {
			desiredCol = currentCol
			t.desiredCursorCol = desiredCol
		}
		if desiredCol > len(lines[targetLine]) {
			desiredCol = len(lines[targetLine])
		}
		newPos := t.getCharPosFromLineAndCol(targetLine, desiredCol)
		t.updateSelection(newPos)
		t.desiredCursorCol = -1
	}
}

// ---------------------
func (t *TextAreaSelection) handleShiftHome() {
	if t.selectionStart != t.selectionEnd || !t.hasSelectionStarted() {
		t.pushUndo()
	}

	// Initialize selectionStart if no selection is active
	if t.selectionStart == t.selectionEnd {
		t.setSelectionStart(t.cursorPos)
	}
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndCol(currentLine, 0)
	t.updateSelection(newPos)
}

func (t *TextAreaSelection) handleShiftEnd() {
	// Only push to undo stack if a selection will be changed
	if t.selectionStart != t.selectionEnd || !t.hasSelectionStarted() {
		t.pushUndo()
	}

	// Initialize selectionStart if no selection is active
	if t.selectionStart == t.selectionEnd {
		t.setSelectionStart(t.cursorPos)
	}

	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	newPos := t.getCharPosFromLineAndCol(currentLine, len(lines[currentLine]))
	t.updateSelection(newPos)
}

// Helper method to check if selection has started
func (t *TextAreaSelection) hasSelectionStarted() bool {
	return t.selectionStart != t.selectionEnd
}

func (t *TextAreaSelection) handleUndo() {
	if len(t.undoStack) > 0 {
		// Pop the last state from undoStack
		lastState := t.undoStack[len(t.undoStack)-1]
		t.undoStack = t.undoStack[:len(t.undoStack)-1]
		// Push the current state to redoStack
		currentState := TextState{
			Text:      t.text,
			CursorPos: t.cursorPos,
		}
		t.redoStack = append(t.redoStack, currentState)
		// Restore the last state
		t.text = lastState.Text
		t.setCursorPos(lastState.CursorPos)
		t.clearSelection()
		t.isTextChanged = true
		t.counter = 0 // Reset blink counter
	}
}

func (t *TextAreaSelection) handleRedo() {
	if len(t.redoStack) > 0 {
		// Pop the last state from redoStack
		lastState := t.redoStack[len(t.redoStack)-1]
		t.redoStack = t.redoStack[:len(t.redoStack)-1]
		// Push the current state to undoStack
		currentState := TextState{
			Text:      t.text,
			CursorPos: t.cursorPos,
		}
		t.undoStack = append(t.undoStack, currentState)
		// Restore the last state
		t.text = lastState.Text
		t.setCursorPos(lastState.CursorPos)
		t.clearSelection()
		t.isTextChanged = true
		t.counter = 0 // Reset blink counter
	}
}
