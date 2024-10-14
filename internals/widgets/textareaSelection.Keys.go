package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (t *TextAreaSelection) checkKeyboardInput() error {

	// Define the keys that support repeat
	repeatKeys := []ebiten.Key{
		ebiten.KeyBackspace,
		ebiten.KeyDelete,
		ebiten.KeyTab,
		ebiten.KeyEnter,
		ebiten.KeyLeft,
		ebiten.KeyRight,
		ebiten.KeyHome,
		ebiten.KeyEnd,
		ebiten.KeyUp,
		ebiten.KeyDown,
		ebiten.KeyC,
		ebiten.KeyX,
		ebiten.KeyV,
		ebiten.KeyZ,
		ebiten.KeyY,
		ebiten.KeyA,
		ebiten.KeyPageUp,
		ebiten.KeyPageDown,
	}

	// Handle repeat keys
	for _, key := range repeatKeys {
		if ebiten.IsKeyPressed(key) {
			// Initialize key state if not present
			if _, exists := t.heldKeys[key]; !exists {
				t.heldKeys[key] = &KeyState{
					InitialPress:    true,
					FramesHeld:      0,
					FramesUntilNext: t.keyRepeatInitialDelay,
				}
			}
			keyState := t.heldKeys[key]
			if keyState.InitialPress {
				// Handle the initial key press
				t.checkKeyPress(key)
				keyState.InitialPress = false
			} else {
				// Increment frames held
				keyState.FramesHeld++
				// Check if it's time to repeat the action
				if keyState.FramesHeld >= keyState.FramesUntilNext {
					// Handle the repeated action
					t.checkKeyPress(key)
					// Reset frames until next action
					keyState.FramesUntilNext = t.keyRepeatInterval
				}
			}
		} else {
			// Remove key from heldKeys when released
			delete(t.heldKeys, key)
		}
	}

	// Handle character input
	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			t.text = t.text[:t.cursorPos] + string(char) + t.text[t.cursorPos:]
			t.cursorPos++
			t.clearSelection()
		}
	}

	t.counter++
	return nil
}

func (t *TextAreaSelection) checkKeyPress(key ebiten.Key) {
	// If there is an active selection and Shift is not pressed,
	// move the cursor to the appropriate end of the selection and clear the selection.
	if t.selectionStart != t.selectionEnd && !t.isShiftPressed() {
		switch key {
		case ebiten.KeyLeft, ebiten.KeyUp, ebiten.KeyHome:
			// Move cursor to the start of the selection
			t.setCursorPos(t.getSelectionBoundsStart())
		case ebiten.KeyRight, ebiten.KeyDown, ebiten.KeyEnd:
			// Move cursor to the end of the selection
			t.setCursorPos(t.getSelectionBoundsEnd())
		case ebiten.KeyBackspace, ebiten.KeyDelete:
			t.handleDelete()
		}

		// Clear the selection
		t.clearSelection()
		return // Exit early to prevent further processing
	}

	switch key {
	case ebiten.KeyTab:
		t.handleTab()
	case ebiten.KeyEnter:
		t.handleEnter()
	case ebiten.KeyLeft:
		if t.isCtrlPressed() {
			if t.isShiftPressed() {
				t.handleCtrlShiftLeftArrow()
			}
			t.handleCtrlLeftArrow()
		}

		t.handleLeftArrow()

	case ebiten.KeyRight:
		if t.isCtrlPressed() {
			if t.isShiftPressed() {
				t.handleCtrlShiftRightArrow()
			}
			t.handleCtrlRightArrow()
		} else {
			t.handleRightArrow()
		}
	case ebiten.KeyHome:
		if t.isCtrlPressed() {
			if t.isShiftPressed() {
				t.handleShiftHome()
			} else {
				t.handleCtrlHome()
			}
		}
		t.handleHome()
	case ebiten.KeyEnd:
		if t.isShiftPressed() {
			t.handleShiftEnd()
		}
		if t.isCtrlPressed() {
			if t.isShiftPressed() {
				t.handleShiftEnd()
			}
			t.handleCtrlEnd()
		}
		t.handleEnd()
	case ebiten.KeyUp:
		if t.isCtrlPressed() {
			if t.isShiftPressed() {
				t.handleCtrlShiftUpArrow()
			}
			t.handleCtrlUpArrow()
		} else if t.isShiftPressed() {
			t.handleShiftUp()
		} else {
			t.handleUpArrow()
		}
	case ebiten.KeyDown:
		if t.isCtrlPressed() {
			if t.isShiftPressed() {
				t.handleCtrlShiftDownArrow()
			}
			t.handleCtrlDownArrow()
		} else if t.isShiftPressed() {
			t.handleShiftDown()
		} else {
			t.handleDownArrow()
		}
	case ebiten.KeyC:
		if t.isCtrlPressed() {
			t.handleCopySelection()
		}
	case ebiten.KeyX:
		if t.isCtrlPressed() {
			t.handleCutSelection()
		}
	case ebiten.KeyV:
		if t.isCtrlPressed() {
			t.handlePasteClipboard()
		}
	case ebiten.KeyZ:
		if t.isCtrlPressed() {
			t.handleUndo()
		}
	case ebiten.KeyY:
		if t.isCtrlPressed() {
			t.handleRedo()
		}
	case ebiten.KeyA:
		if t.isCtrlPressed() {
			t.handleSelectAll()
		}
	case ebiten.KeyBackspace:
		if t.isCtrlPressed() {
			t.handleCtrlBackspace()
		}
		t.handleBackspace()
	case ebiten.KeyDelete:
		if t.isCtrlPressed() {
			t.handleCtrlDelete()
		}
		t.handleDelete()
	case ebiten.KeyPageUp:
		t.handlePageUp()
	case ebiten.KeyPageDown:
		t.handlePageDown()
	}
}
