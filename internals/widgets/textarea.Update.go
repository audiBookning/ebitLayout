package widgets

import (

	//"example.com/menu/internals/textwrapper02"

	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (t *TextArea) Update() error {

	if t.isTextChanged {
		t.cachedLines = strings.Split(t.text, "\n")
		t.isTextChanged = false
	}

	// Single, double, triple, and Shift+Click detection
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if t.isOverScrollbar(x, y) {
			// Clicked on scrollbar
			t.SetIsDraggingThumb(true)
			t.dragOffsetY = float64(y-t.scrollbarY) - t.scrollbarThumbY
			return nil
		}

		if t.isShiftPressed() {
			// Handle Shift + Click: Extend selection
			charPos := t.getCharPosFromPosition(x, y)

			// If there's no existing selection, set the selection start to the current cursor position
			if t.selection.selectionStart == t.selection.selectionEnd {
				poos := clamp(t.cursorPos, 0, len(t.text))
				t.selection.setSelectionStart(poos)

			}

			// Update the selection end and cursor position
			poos := clamp(charPos, 0, len(t.text))
			t.selection.setSelectionEnd(poos)
			t.setCursorPos(charPos)
			t.selection.SetIsSelecting(true)
		} else {
			// Handle single, double, and triple clicks
			t.isMouseLeftPressed = true
			t.clicked = true
			currentFrame := t.counter
			if currentFrame-t.lastClickTime <= t.doubleClickThreshold {
				t.clickCount++
			} else {
				t.clickCount = 1
			}
			t.lastClickTime = currentFrame

			if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
				switch t.clickCount {
				case 1:
					// Single click
					t.hasFocus = true
					charPos := t.getCharPosFromPosition(x, y)
					t.setCursorPos(charPos)
					poos := clamp(charPos, 0, len(t.text))
					// todo: clearSelection with param
					t.selection.setSelectionStart(poos)
					t.selection.setSelectionEnd(poos)

					t.selection.SetIsSelecting(false)
					t.SetIsDraggingThumb(false)
				case 2:
					// Double click
					charPos := t.getCharPosFromPosition(x, y)
					t.selectWordAt(charPos)
					t.doubleClickHandled = true // Flag to prevent further handling in this update
				case 3:
					// Triple click
					t.selectEntireLineAt(x, y)
					t.clickCount = 0 // Reset click count after handling triple click
				}
			} else {
				// Clicked outside text area
				t.hasFocus = false
				t.selection.SetIsSelecting(false)
			}
		}
	}

	// Handle mouse movement while left button is pressed
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Check if a double-click was just handled to prevent unwanted selection changes
		if t.doubleClickHandled {
			// Do not process selection adjustments while a double-click is handled
			// Wait until the mouse button is released to reset the flag
		} else {
			x, y := ebiten.CursorPosition()
			if t.isDraggingThumb {
				t.dragScrollbar(y)
			} else if t.hasFocus && !(t.isShiftPressed() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)) {
				if t.isOverScrollbar(x, y) {
					// Prevent text selection when clicking on scrollbar
				} else {
					charPos := t.getCharPosFromPosition(x, y)
					if !t.selection.isSelecting {
						// Start selection on first movement after click
						// ??? not using charPos
						t.selection.SetIsSelecting(true)
						poos := clamp(t.cursorPos, 0, len(t.text))
						t.selection.setSelectionStart(poos)
					}
					poos := clamp(charPos, 0, len(t.text))
					t.selection.setSelectionEnd(poos)
					t.setCursorPos(charPos)
				}
			}
		}
	}

	// Handle mouse button release (mouse up)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		t.isMouseLeftPressed = false
		if t.isDraggingThumb {
			t.SetIsDraggingThumb(false)
		}
		if t.selection.isSelecting && !t.isDraggingThumb {
			// Finalize selection on mouse release
			t.selection.SetIsSelecting(false)
		}
		if t.doubleClickHandled {
			// Reset the double click handled flag upon mouse release
			t.doubleClickHandled = false
		}
	}

	// Handle keyboard input when focused
	if t.hasFocus {
		t.checkKeyboardInput()
	}

	// Handle mouse wheel scrolling with smooth scrolling
	_, yScroll := ebiten.Wheel()
	if yScroll != 0 {
		const linesPerWheel = 3
		totalLines := len(strings.Split(t.text, "\n"))
		targetScrollOffset := clamp(t.scrollOffset-int(yScroll)*linesPerWheel, 0, max(t.scrollOffset, totalLines-t.maxLines))
		// Implement smooth transition to targetScrollOffset
		scrollSpeed := 1 // Adjust this value for faster or slower scrolling
		if t.scrollOffset < targetScrollOffset {
			t.SetScrollOffset(t.scrollOffset + scrollSpeed)
		} else if t.scrollOffset > targetScrollOffset {
			t.SetScrollOffset(t.scrollOffset - scrollSpeed)
		}
	}

	t.counter++
	return nil
}
