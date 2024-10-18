package widgets

import (
	"fmt"
	"strings"

	"golang.design/x/clipboard"
)

func (t *TextArea) handlePageDown() {
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

func (t *TextArea) handlePageUp() {
	t.pushUndo()
	// Calculate the new scroll offset
	newScrollOffset := t.scrollOffset - t.maxLines
	if newScrollOffset < 0 {
		newScrollOffset = 0
	}
	// Update the scroll offset
	t.SetScrollOffset(newScrollOffset)

}

func (t *TextArea) handleCtrlShiftLeftArrow() {
	t.pushUndo()
	newPos := t.moveToWordStart(t.cursorPos)
	t.updateSelection(newPos)
}

func (t *TextArea) handleCtrlShiftRightArrow() {
	t.pushUndo()
	newPos := t.moveToWordEnd(t.cursorPos)
	t.updateSelection(newPos)
}

func (t *TextArea) handleCtrlShiftUpArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndColWithclamp(currentLine, 0)
	t.updateSelection(newPos)
}

func (t *TextArea) handleCtrlShiftDownArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines) {
		newPos := t.getCharPosFromLineAndColWithclamp(currentLine, len(lines[currentLine]))
		t.updateSelection(newPos)
	}
}

func (t *TextArea) handleCtrlUpArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndColWithclamp(currentLine, 0)
	if t.isShiftPressed() {
		t.updateSelection(newPos)
	} else {

		t.selection.ClearSelection(t.cursorPos)
		t.setCursorPos(newPos)
	}
}

func (t *TextArea) handleCtrlDownArrow() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines) {
		newPos := t.getCharPosFromLineAndColWithclamp(currentLine, len(lines[currentLine]))
		if t.isShiftPressed() {
			t.updateSelection(newPos)
		} else {
			t.selection.ClearSelection(t.cursorPos)
			t.setCursorPos(newPos)
		}
	}
}

func (t *TextArea) handleCtrlShiftHome() {
	t.pushUndo()
	// Select from cursor to beginning of text
	t.selection.setSelectionStart(0)
	t.selection.setSelectionEnd(t.selection.selectionStart)
	t.setCursorPos(0)
	// Scroll to the top of the textarea
	t.SetScrollOffset(0)
}

func (t *TextArea) handleCtrlShiftEnd() {
	t.pushUndo()
	// Select from cursor to end of text
	// TODO: not necessary, but just for clarity of the logic
	t.selection.setSelectionStart(t.selection.selectionStart) // Use the existing selection start as the start
	t.selection.setSelectionEnd(len(t.text))
	t.setCursorPos(len(t.text))
	// Scroll to the bottom of the textarea
	maxScrollOffset := len(strings.Split(t.text, "\n")) - t.maxLines
	if maxScrollOffset > 0 {
		t.SetScrollOffset(maxScrollOffset)
	}
}

func (t *TextArea) handleCtrlHome() {
	t.pushUndo()
	// Move cursor to the very beginning of the text
	t.setCursorPos(0)
	if t.isShiftPressed() {
		t.selection.setSelectionEnd(0)

	} else {
		t.selection.ClearSelection(t.cursorPos)
	}
	// Scroll to the top of the textarea
	t.SetScrollOffset(0)
}

func (t *TextArea) handleCtrlEnd() {
	t.pushUndo()
	// Move cursor to the very end of the text
	t.setCursorPos(len(t.text))
	if t.isShiftPressed() {

		t.selection.setSelectionEnd(len(t.text))
	} else {
		t.selection.ClearSelection(t.cursorPos)
	}
	// Scroll to the bottom of the textarea
	maxScrollOffset := len(strings.Split(t.text, "\n")) - t.maxLines
	if maxScrollOffset > 0 {
		t.SetScrollOffset(maxScrollOffset)
	}
}

func (t *TextArea) handleCopySelection() {
	if t.selection.selectionStart == t.selection.selectionEnd {
		// No selection to copy
		fmt.Println("handleCopySelection - No selection to copy.")
		return
	}
	minPos, maxPos := t.selection.getSelectionBounds()
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
func (t *TextArea) handleCutSelection() {
	if t.selection.selectionStart == t.selection.selectionEnd {
		// No selection to cut
		fmt.Println("handleCutSelection - No selection to cut.")
		return
	}
	t.pushUndo()
	minPos, maxPos := t.selection.getSelectionBounds()
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
	t.selection.ClearSelection(t.cursorPos)
	t.isTextChanged = true
}

// handlePasteClipboard inserts text from the OS clipboard into the text area at the current cursor position
func (t *TextArea) handlePasteClipboard() {
	t.pushUndo() // Ensure each paste is undoable individually
	clipboardBytes := clipboard.Read(clipboard.FmtText)
	clipboardText := string(clipboardBytes)
	if t.selection.selectionStart != t.selection.selectionEnd {
		// Replace selected text with clipboard text
		minPos, maxPos := t.selection.getSelectionBounds()
		t.text = t.text[:minPos] + clipboardText + t.text[maxPos:]
		t.cursorPos = minPos + len(clipboardText)
		t.isTextChanged = true
	} else {
		// Insert clipboard text at cursor position
		t.text = t.text[:t.cursorPos] + clipboardText + t.text[t.cursorPos:]
		t.cursorPos += len(clipboardText)
		t.isTextChanged = true
	}
	t.selection.ClearSelection(t.cursorPos)
}

func (t *TextArea) handleBackspace() {
	if t.selection.selectionStart != t.selection.selectionEnd {
		t.pushUndo()
		t.deleteSelection()
		t.isTextChanged = true
	} else if t.cursorPos > 0 {
		t.pushUndo()
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.cursorPos--
		t.isTextChanged = true
		t.selection.ClearSelection(t.cursorPos)
	}
}

func (t *TextArea) handleDelete() {
	if t.selection.selectionStart != t.selection.selectionEnd {
		t.pushUndo()
		t.deleteSelection()
		t.isTextChanged = true
	} else if t.cursorPos < len(t.text) {
		t.pushUndo()
		t.text = t.text[:t.cursorPos] + t.text[t.cursorPos+1:]
		t.isTextChanged = true
		t.selection.ClearSelection(t.cursorPos)
	}
}

func (t *TextArea) handleCtrlBackspace() {
	if t.selection.selectionStart != t.selection.selectionEnd {
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
		t.selection.ClearSelection(t.cursorPos)
	}
}

// handleCtrlDelete deletes text from the current cursor position to the end of the word
func (t *TextArea) handleCtrlDelete() {
	if t.selection.selectionStart != t.selection.selectionEnd {
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
		t.selection.ClearSelection(t.cursorPos)
	}
}

func (t *TextArea) handleTab() {
	t.pushUndo()
	if t.selection.isSelecting {
		t.indentSelection()
	} else {
		t.text = t.text[:t.cursorPos] + strings.Repeat(" ", t.tabWidth) + t.text[t.cursorPos:]
		t.cursorPos += t.tabWidth
	}
	t.selection.ClearSelection(t.cursorPos)
}

func (t *TextArea) handleEnter() {
	t.pushUndo()
	t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
	t.cursorPos++
	t.selection.ClearSelection(t.cursorPos)
}

func (t *TextArea) handleLeftArrow() {
	t.pushUndo()
	if t.selection.selectionStart != t.selection.selectionEnd {
		t.selection.ClearSelection(t.cursorPos)
	} else {
		if t.cursorPos > 0 {
			t.cursorPos--
			t.selection.ClearSelection(t.cursorPos)
		}
	}
}

func (t *TextArea) handleShiftLeftArrow() {
	t.pushUndo()
	t.updateSelectionWithShiftKey(-1)
}

func (t *TextArea) handleShiftRightArrow() {
	t.updateSelectionWithShiftKey(1)
}

func (t *TextArea) handleRightArrow() {
	t.pushUndo()
	if t.selection.selectionStart != t.selection.selectionEnd {
		t.selection.ClearSelection(t.cursorPos)
	} else {
		if t.cursorPos < len(t.text) {
			t.cursorPos++
			t.selection.ClearSelection(t.cursorPos)
		}
	}
}

func (t *TextArea) handleSelectAll() {
	t.pushUndo()
	t.selection.setSelectionStart(0)
	t.selection.setSelectionEnd(len(t.text))
	t.setCursorPos(len(t.text))
}

func (t *TextArea) handleHome() {
	t.pushUndo()
	line, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndColWithclamp(line, 0)
	t.selection.ClearSelection(t.cursorPos)
	t.setCursorPos(newPos)

}

func (t *TextArea) handleEnd() {
	t.pushUndo()
	line, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	newPos := t.getCharPosFromLineAndColWithclamp(line, len(lines[line]))

	t.selection.ClearSelection(t.cursorPos)
	t.setCursorPos(newPos)
}

func (t *TextArea) handleCtrlLeftArrow() {
	t.pushUndo()
	newPos := t.moveToWordStart(t.cursorPos)
	t.selection.ClearSelection(t.cursorPos)
	t.setCursorPos(newPos)
}

func (t *TextArea) handleCtrlRightArrow() {
	t.pushUndo()
	newPos := t.moveToWordEnd(t.cursorPos)

	t.selection.ClearSelection(t.cursorPos)
	t.setCursorPos(newPos)
}

// ---------------------
func (t *TextArea) handleUpArrow() {
	t.pushUndo()
	currentLine, currentCol := t.getCursorLineAndColForPos(t.cursorPos)
	if currentLine > 0 {
		targetLine := currentLine - 1
		lines := strings.Split(t.text, "\n")
		targetCol := currentCol
		if targetCol > len(lines[targetLine]) {
			targetCol = len(lines[targetLine])
		}
		newPos := t.getCharPosFromLineAndColWithclamp(targetLine, targetCol)
		t.setCursorPos(newPos)
		t.selection.ClearSelection(t.cursorPos)
	}
}

func (t *TextArea) handleDownArrow() {
	t.pushUndo()
	currentLine, currentCol := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if currentLine < len(lines)-1 {
		targetLine := currentLine + 1
		targetCol := currentCol
		if targetCol > len(lines[targetLine]) {
			targetCol = len(lines[targetLine])
		}
		newPos := t.getCharPosFromLineAndColWithclamp(targetLine, targetCol)
		t.setCursorPos(newPos)
		t.selection.ClearSelection(t.cursorPos)
	}
}

func (t *TextArea) handleShiftUp() {
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
		newPos := t.getCharPosFromLineAndColWithclamp(targetLine, desiredCol)
		t.updateSelection(newPos)
		t.desiredCursorCol = -1
	}
}

func (t *TextArea) handleShiftDown() {
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
		newPos := t.getCharPosFromLineAndColWithclamp(targetLine, desiredCol)
		t.updateSelection(newPos)
		t.desiredCursorCol = -1
	}
}

// ---------------------
func (t *TextArea) handleShiftHome() {
	if t.selection.selectionStart != t.selection.selectionEnd || !t.hasSelectionStarted() {
		t.pushUndo()
	}

	// Initialize selectionStart if no selection is active
	if t.selection.selectionStart == t.selection.selectionEnd {
		t.selection.setSelectionStart(t.cursorPos)
	}
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndColWithclamp(currentLine, 0)
	t.updateSelection(newPos)
}

func (t *TextArea) handleShiftEnd() {
	// Only push to undo stack if a selection will be changed
	if t.selection.selectionStart != t.selection.selectionEnd || !t.hasSelectionStarted() {
		t.pushUndo()
	}

	// Initialize selectionStart if no selection is active
	if t.selection.selectionStart == t.selection.selectionEnd {
		t.selection.setSelectionStart(t.cursorPos)
	}

	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	newPos := t.getCharPosFromLineAndColWithclamp(currentLine, len(lines[currentLine]))
	t.updateSelection(newPos)
}

// Helper method to check if selection has started
func (t *TextArea) hasSelectionStarted() bool {
	return t.selection.selectionStart != t.selection.selectionEnd
}

func (t *TextArea) handleUndo() {
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
		t.selection.ClearSelection(t.cursorPos)
		t.isTextChanged = true
		t.counter = 0 // Reset blink counter
	}
}

func (t *TextArea) handleRedo() {
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
		t.selection.ClearSelection(t.cursorPos)
		t.isTextChanged = true
		t.counter = 0 // Reset blink counter
	}
}
