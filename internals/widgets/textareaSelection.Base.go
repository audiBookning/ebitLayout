package widgets

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (t *TextAreaSelection) drawCursor(screen *ebiten.Image) {
	cursorLine, cursorCol := t.getCursorLineAndCol()
	if cursorLine >= t.scrollOffset && cursorLine < t.scrollOffset+t.maxLines {

		cursorX := t.x + t.textWidth(strings.Split(t.text, "\n")[cursorLine][:cursorCol])

		cursorY := float64(t.y) + float64(cursorLine-t.scrollOffset)*t.lineHeight

		// Clamp the cursor position within textarea bounds
		cursorX = clamp(cursorX, t.x, t.x+t.w)
		cursorY = clampFloat(cursorY, float64(t.y), float64(t.y+t.h)-t.lineHeight)

		// Render the blinking cursor
		if t.counter%(t.cursorBlinkRate*2) < t.cursorBlinkRate {
			vector.DrawFilledRect(screen,
				float32(cursorX),
				float32(cursorY),
				2,
				float32(t.lineHeight),
				color.RGBA{0, 0, 0, 255},
				true)
		}
	}
}

func (t *TextAreaSelection) Draw(screen *ebiten.Image) {

	// Split lines only if text has changed
	if t.isTextChanged {
		t.cachedLines = strings.Split(t.text, "\n")
		t.isTextChanged = false
	}
	lines := t.cachedLines

	// Draw the background of the text area
	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)

	yOffset := float64(t.y)
	// Apply scroll offset
	startLine := t.scrollOffset
	endLine := clamp(t.scrollOffset+t.maxLines, 0, len(lines))

	// Retrieve normalized selection bounds
	minPos, maxPos := t.getSelectionBounds()

	for i := startLine; i < endLine; i++ {
		line := lines[i]

		lineText := line
		lineX := t.x + t.paddingLeft
		lineY := int(yOffset+t.lineHeight) + t.paddingTop

		// Draw selection if active and within this line
		if minPos != maxPos {
			startLineSel, startCol := t.getCursorLineAndColForPos(minPos)
			endLineSel, endCol := t.getCursorLineAndColForPos(maxPos)

			if i >= startLineSel && i <= endLineSel {
				// Determine selection bounds within the current line
				var selStart, selEnd int
				if i == startLineSel {
					selStart = startCol
				} else {
					selStart = 0
				}
				if i == endLineSel {
					selEnd = endCol
				} else {
					selEnd = len(line)
				}

				// Calculate x positions based on character widths
				selectionXStart := t.x + t.textWidth(lineText[:selStart])
				selectionXEnd := t.x + t.textWidth(lineText[:selEnd])

				// Clamp the selection rectangle within textarea bounds
				selectionXStart = clamp(selectionXStart, t.x, t.x+t.w)
				selectionXEnd = clamp(selectionXEnd, t.x, t.x+t.w)

				// Clamp the yOffset within textarea bounds
				clampedYOffset := clampFloat(yOffset, float64(t.y), float64(t.y+t.h)-t.lineHeight)

				// Draw the selection rectangle
				vector.DrawFilledRect(screen,
					float32(selectionXStart),
					float32(clampedYOffset),
					float32(selectionXEnd-selectionXStart),
					float32(t.lineHeight),
					color.RGBA{0, 0, 255, 128},
					true)
			}
		}

		t.textWrapper.DrawText(screen, lineText, lineX, lineY)

		yOffset += t.lineHeight
	}

	// Draw the scrollbar if content exceeds maxLines
	if len(lines) > t.maxLines {
		t.drawScrollbar(screen, len(lines))
	}

	// Draw the cursor if the text area has focus and the cursor is within the visible lines
	if t.hasFocus {
		t.drawCursor(screen)
	}
}

func (t *TextAreaSelection) Update() error {

	// Handle mouse button just pressed (mouse down)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		t.clicked = true
		currentFrame := t.counter
		if currentFrame-t.lastClickTime <= t.doubleClickThreshold {
			t.clickCount++
		} else {
			t.clickCount = 1
			t.doubleClickHandled = false
		}
		t.lastClickTime = currentFrame

		x, y := ebiten.CursorPosition()

		if t.isOverScrollbar(x, y) {
			// Clicked on scrollbar
			t.isDraggingThumb = true
			t.dragOffsetY = float64(y-t.scrollbarY) - t.scrollbarThumbY // Adjust offset calculation
			return nil                                                  // Exit early to prevent further processing
		} else if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			if t.clickCount == 2 {
				// Double-click detected
				charPos := t.getCharPosFromPosition(x, y)
				t.selectWordAt(charPos)
				t.doubleClickHandled = true // Set the flag

			} else {
				// Single click
				t.hasFocus = true
				charPos := t.getCharPosFromPosition(x, y)
				t.setCursorPos(charPos) //
				t.setSelectionStart(charPos)
				t.setSelectionEnd(charPos)
				// Initially, no selection is active
				t.isSelecting = false
			}
		} else {
			t.hasFocus = false
			t.isSelecting = false
		}
	}

	// Handle mouse movement
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if t.isDraggingThumb {
			t.dragScrollbar(y)
		} else if t.hasFocus && !t.doubleClickHandled {
			if t.isOverScrollbar(x, y) {
				// Prevent text selection when clicking on scrollbar
			} else {
				charPos := t.getCharPosFromPosition(x, y)
				if !t.isSelecting {
					// Start selection on first movement after click
					t.isSelecting = true
					t.setSelectionStart(t.cursorPos)
				}
				t.setSelectionEnd(charPos)
				t.setCursorPos(charPos) // called 2 times.
			}
		}
	}

	// Handle mouse button release (mouse up)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if t.isDraggingThumb {
			t.isDraggingThumb = false
		}
		if t.isSelecting && !t.isDraggingThumb {
			// Finalize selection on mouse release
			t.isSelecting = false
		}
	}

	// Handle keyboard input when focused
	if t.hasFocus {
		t.checkKeyboardInput()
	}

	// Handle mouse wheel scrolling
	_, yScroll := ebiten.Wheel()
	if yScroll != 0 {
		t.SetScrollOffset(clamp(t.scrollOffset-int(yScroll), 0, len(strings.Split(t.text, "\n"))-t.maxLines))
	}

	t.counter++
	return nil
}
