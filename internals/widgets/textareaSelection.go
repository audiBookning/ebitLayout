package widgets

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"golang.design/x/clipboard"
)

const (
	// keyRepeatInitialDelay is the number of frames to wait before repeating the key action
	keyRepeatInitialDelay = 30 // Approximately 0.5 seconds at 60 TPS

	// keyRepeatInterval is the number of frames between each repeated action
	keyRepeatInterval = 5 // Approximately 10 times per second at 60 TPS
)

// KeyState tracks the repeat state of a specific key
type KeyState struct {
	InitialPress    bool // Indicates if the initial press has been handled
	FramesHeld      int  // Number of frames the key has been held down
	FramesUntilNext int  // Frames remaining until the next action
}

type TextState struct {
	Text      string
	CursorPos int
}

type TextAreaSelection struct {
	text                 string
	hasFocus             bool
	cursorPos            int
	counter              int
	selectionStart       int
	selectionEnd         int
	isSelecting          bool
	x, y, w, h           int
	maxLines             int
	cursorBlinkRate      int
	tabWidth             int
	lineHeight           int
	font                 font.Face
	heldKeys             map[ebiten.Key]*KeyState
	undoStack            []TextState
	redoStack            []TextState
	desiredCursorCol     int
	lastClickTime        int  // Frame count of the last click
	clickCount           int  // Number of consecutive clicks
	doubleClickThreshold int  // Threshold frames to consider as a double-click
	doubleClickHandled   bool // Indicates if a double-click has just been handled
	scrollOffset         int

	// Scrollbar Fields
	scrollbarWidth  int     // Width of the scrollbar
	scrollbarX      int     // X position of the scrollbar
	scrollbarY      int     // Y position of the scrollbar
	scrollbarHeight int     // Height of the scrollbar
	scrollbarThumbY float64 // Y position of the scrollbar thumb
	scrollbarThumbH float64 // Height of the scrollbar thumb
	isDraggingThumb bool    // Indicates if the scrollbar thumb is being dragged
	dragOffsetY     float64 // Offset between mouse position and thumb position during drag
}

func (t *TextAreaSelection) setCursorPos(pos int) {
	t.cursorPos = clamp(pos, 0, len(t.text))
}

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
	t.scrollOffset = newScrollOffset

}

func (t *TextAreaSelection) handlePageUp() {
	t.pushUndo()
	// Calculate the new scroll offset
	newScrollOffset := t.scrollOffset - t.maxLines
	if newScrollOffset < 0 {
		newScrollOffset = 0
	}

	// Update the scroll offset
	t.scrollOffset = newScrollOffset

}

func (t *TextAreaSelection) handleKeyboardInput() error {

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
	}

	// Handle repeat keys
	for _, key := range repeatKeys {
		if ebiten.IsKeyPressed(key) {
			// Initialize key state if not present
			if _, exists := t.heldKeys[key]; !exists {
				t.heldKeys[key] = &KeyState{
					InitialPress:    true,
					FramesHeld:      0,
					FramesUntilNext: keyRepeatInitialDelay,
				}
			}
			keyState := t.heldKeys[key]
			if keyState.InitialPress {
				// Handle the initial key press
				t.handleKeyPress(key)
				keyState.InitialPress = false
			} else {
				// Increment frames held
				keyState.FramesHeld++
				// Check if it's time to repeat the action
				if keyState.FramesHeld >= keyState.FramesUntilNext {
					// Handle the repeated action
					t.handleKeyPress(key)
					// Reset frames until next action
					keyState.FramesUntilNext = keyRepeatInterval
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

	// ------------------
	// Ctrl+C, Ctrl+X, Ctrl+V, Ctrl+Z, Ctrl+Y, Ctrl+A,
	// Ctrl+Home, Ctrl+End, Ctrl+Backspace, Ctrl+Delete
	if t.isCtrlPressed() {
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			t.handleCopySelection()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyX) {
			t.handleCutSelection()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			t.handlePasteClipboard()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
			t.handleUndo()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyY) {
			t.handleRedo()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			t.handleSelectAll()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			t.handleCtrlBackspace()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
			t.handleCtrlDelete()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
			t.handleCtrlHome()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
			t.handleCtrlEnd()
		}
	} else {
		// Handle Page Up and Page Down keys
		if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
			t.pushUndo()
			t.handlePageUp()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
			t.pushUndo()
			t.handlePageDown()
		}
	}

	// ------------------
	// Clamp cursor position to valid range
	//t.setCursorPos(clamp(t.cursorPos, 0, len(t.text)))

	t.counter++
	return nil
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
	t.scrollOffset = 0
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
		t.scrollOffset = maxScrollOffset
	}
}

func (t *TextAreaSelection) setSelectionStart(pos int) {
	t.selectionStart = pos
	fmt.Printf("Selection Start Updated: %d\n", pos)
}

func (t *TextAreaSelection) setSelectionEnd(pos int) {
	t.selectionEnd = pos
	fmt.Printf("Selection End Updated: %d\n", pos)
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

	fmt.Printf("Selection Updated: Start=%d, End=%d, CursorPos=%d\n", t.selectionStart, t.selectionEnd, t.cursorPos)
}

func NewTextAreaSelection(x, y, w, h, maxLines int, startTxt string) *TextAreaSelection {
	err := clipboard.Init()
	if err != nil {
		return nil
	}

	return &TextAreaSelection{
		x:                    x,
		y:                    y,
		w:                    w,
		h:                    h,
		maxLines:             maxLines,
		cursorBlinkRate:      30,
		tabWidth:             4,
		lineHeight:           20,
		font:                 basicfont.Face7x13,
		heldKeys:             make(map[ebiten.Key]*KeyState),
		desiredCursorCol:     -1,
		lastClickTime:        0,     // Frame count of the last click
		clickCount:           0,     // Number of consecutive clicks
		doubleClickThreshold: 30,    // Threshold frames to consider as a double-click
		doubleClickHandled:   false, // Indicates if a double-click has just been handled
		scrollOffset:         0,
		text:                 startTxt, // Default text added here
	}
}

func (t *TextAreaSelection) Draw(screen *ebiten.Image) {

	// Draw the background of the text area
	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)

	lines := strings.Split(t.text, "\n")
	yOffset := t.y

	// Apply scroll offset
	startLine := t.scrollOffset
	endLine := clamp(t.scrollOffset+t.maxLines, 0, len(lines))

	// Retrieve normalized selection bounds
	minPos, maxPos := t.getSelectionBounds()

	for i := startLine; i < endLine; i++ {
		line := lines[i]

		lineText := line
		lineX := t.x
		lineY := yOffset + t.lineHeight/2

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
				clampedYOffset := clamp(yOffset, t.y, t.y+t.h-t.lineHeight)

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

		// Draw the actual text
		text.Draw(screen, lineText, t.font, lineX, lineY+t.lineHeight/2, color.Black)

		yOffset += t.lineHeight
	}

	// Draw the scrollbar if content exceeds maxLines
	if len(lines) > t.maxLines {
		t.drawScrollbar(screen, len(lines))
	}

	// Draw the cursor if the text area has focus and the cursor is within the visible lines
	if t.hasFocus {
		cursorLine, cursorCol := t.getCursorLineAndCol()
		if cursorLine >= t.scrollOffset && cursorLine < t.scrollOffset+t.maxLines {
			// Calculate cursorX and cursorY
			cursorX := t.x + t.textWidth(strings.Split(t.text, "\n")[cursorLine][:cursorCol])
			cursorY := t.y + (cursorLine-t.scrollOffset)*t.lineHeight

			// Clamp the cursor position within textarea bounds
			cursorX = clamp(cursorX, t.x, t.x+t.w)
			cursorY = clamp(cursorY, t.y, t.y+t.h-t.lineHeight)

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
}

func (t *TextAreaSelection) Update() error {

	// Handle mouse button just pressed (mouse down)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		currentFrame := t.counter
		if currentFrame-t.lastClickTime <= t.doubleClickThreshold {
			t.clickCount++
		} else {
			t.clickCount = 1
			t.doubleClickHandled = false
		}
		t.lastClickTime = currentFrame

		x, y := ebiten.CursorPosition()
		if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			if t.clickCount == 2 {
				// Double-click detected
				charPos := t.getCharPosFromPosition(x, y)
				t.selectWordAt(charPos)
				t.doubleClickHandled = true // Set the flag
			} else {
				// Single click
				t.hasFocus = true
				charPos := t.getCharPosFromPosition(x, y)
				t.setCursorPos(charPos)
				t.setSelectionStart(charPos)
				t.setSelectionEnd(charPos)
				// Initially, no selection is active
				t.isSelecting = false
			}
		} else if t.isOverScrollbar(x, y) {
			// Clicked on scrollbar
			t.isDraggingThumb = true
			t.dragOffsetY = float64(y) - t.scrollbarThumbY
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
		} else if t.hasFocus && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !t.doubleClickHandled {
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
				t.setCursorPos(charPos)
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
		t.handleKeyboardInput()
	}

	// Handle mouse wheel scrolling
	_, yScroll := ebiten.Wheel()
	if yScroll != 0 {
		t.scrollOffset = clamp(t.scrollOffset-int(yScroll), 0, len(strings.Split(t.text, "\n"))-t.maxLines)
	}

	t.counter++
	return nil
}

func (t *TextAreaSelection) isOverScrollbar(x, y int) bool {
	return x >= t.scrollbarX && x <= t.scrollbarX+t.scrollbarWidth && y >= t.y && y <= t.y+t.h
}

func (t *TextAreaSelection) dragScrollbar(mouseY int) {
	// Calculate the new thumb Y position
	newThumbY := float64(mouseY) - t.dragOffsetY

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
	maxScrollOffset := len(strings.Split(t.text, "\n")) - t.maxLines
	if maxScrollOffset < 1 {
		maxScrollOffset = 1
	}

	t.scrollOffset = int((t.scrollbarThumbY / maxThumbY) * float64(maxScrollOffset))
	t.scrollOffset = clamp(t.scrollOffset, 0, maxScrollOffset)
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

	// Calculate the Y position of the scrollbar thumb
	maxScrollOffset := len(strings.Split(t.text, "\n")) - t.maxLines
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
		float32(t.scrollbarX), float32(t.scrollbarY)+float32(t.scrollbarThumbY),
		float32(t.scrollbarWidth), float32(t.scrollbarThumbH),
		color.RGBA{160, 160, 160, 255}, true)
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
	t.isSelecting = false // Ensure no ongoing selection
	completeWord := t.text[byteStart:byteEnd]
	fmt.Printf("Word Selected=%s | pos=%d | Byte Start=%d, Byte End=%d \n", completeWord, pos, start, end)
}

// Helper function to determine if a character is a word separator
func isWordSeparator(r rune) bool {
	separators := " \n\t.,;:!?'\"()-"
	return strings.ContainsRune(separators, r)
}

// Converts rune position to byte position in the original string
func runePosToBytePos(s string, runePos int) int {
	runes := []rune(s)
	if runePos > len(runes) {
		runePos = len(runes)
	}
	return len(string(runes[:runePos]))
}

// Helper function to clamp a value within a range

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// clampFloat limits a float64 value between min and max
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
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

func (t *TextAreaSelection) handleKeyPress(key ebiten.Key) {
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
	case ebiten.KeyBackspace:
		t.handleBackspace()
	case ebiten.KeyDelete:
		t.handleDelete()
	case ebiten.KeyTab:
		t.handleTab()
	case ebiten.KeyEnter:
		t.handleEnter()
	case ebiten.KeyLeft:
		if t.isCtrlPressed() {
			t.handleCtrlLeftArrow()
		} else {
			t.handleLeftArrow()
		}
	case ebiten.KeyRight:
		if t.isCtrlPressed() {
			t.handleCtrlRightArrow()
		} else {
			t.handleRightArrow()
		}
	case ebiten.KeyHome:
		if t.isShiftPressed() {
			t.handleShiftHome()
		} else {
			t.handleHome()
		}
	case ebiten.KeyEnd:
		if t.isShiftPressed() {
			t.handleShiftEnd()
		} else {
			t.handleEnd()
		}
	case ebiten.KeyUp:
		if t.isShiftPressed() {
			t.handleShiftUp()
		} else {
			t.handleUpArrow()
		}
	case ebiten.KeyDown:
		if t.isShiftPressed() {
			t.handleShiftDown()
		} else {
			t.handleDownArrow()
		}
	}
}

func (t *TextAreaSelection) indentSelection() {
	lines := strings.Split(t.text, "\n")
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
	fmt.Printf("Selection Cleared: Start=%d, End=%d, CursorPos=%d\n", t.selectionStart, t.selectionEnd, t.cursorPos)
}

func (t *TextAreaSelection) getCharPosFromLineAndCol(line, col int) int {
	lines := strings.Split(t.text, "\n")
	charPos := 0
	for i := 0; i < line; i++ {
		charPos += len(lines[i]) + 1
	}
	charPos += col
	return charPos
}

func (t *TextAreaSelection) getCursorLineAndColForPos(pos int) (int, int) {
	lines := strings.Split(t.text, "\n")
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

func (t *TextAreaSelection) textWidth(str string) int {
	width := 0
	for _, x := range str {
		awidth, _ := t.font.GlyphAdvance(x)
		width += int(awidth >> 6)
	}
	return width
}

func (t *TextAreaSelection) getCharPosFromPosition(x, y int) int {
	// Adjust the line calculation by adding the scrollOffset
	line := (y-t.y)/t.lineHeight + t.scrollOffset
	col := x - t.x

	lines := strings.Split(t.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	if line < 0 {
		line = 0
	}

	lineText := lines[line]
	colIndex := 0
	accumulatedWidth := 0
	for i, char := range lineText {
		charWidth := t.textWidth(string(char))
		if col < accumulatedWidth+charWidth/2 {
			colIndex = i
			break
		}
		accumulatedWidth += charWidth
		colIndex = i + 1
	}

	// Ensure colIndex does not exceed line length
	if colIndex > len(lineText) {
		colIndex = len(lineText)
	}

	return t.getCharPosFromLineAndCol(line, colIndex)
}

// getSelectionBounds returns the minimum and maximum positions of the current selection
func (t *TextAreaSelection) getSelectionBounds() (int, int) {
	if t.selectionStart <= t.selectionEnd {
		return t.selectionStart, t.selectionEnd
	}
	return t.selectionEnd, t.selectionStart
}

// handleCopySelection copies the selected text to the OS clipboard

// deleteSelection removes the currently selected text and updates the cursor position
func (t *TextAreaSelection) deleteSelection() {
	minPos, maxPos := t.getSelectionBounds()
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.setCursorPos(minPos)
	t.clearSelection()
}

// internals/widgets/textareaSelection01.go
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

func (t *TextAreaSelection) cursorColumn() int {
	_, col := t.getCursorLineAndCol()
	return col
}

func (t *TextAreaSelection) isShiftPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)
}

func (t *TextAreaSelection) updateSelection(newPos int) {
	/* if t.selectionStart > t.selectionEnd {
		t.selectionStart = newPos
	} else {
	} */
	t.setSelectionEnd(newPos)
	t.setCursorPos(newPos)
}

/*
Handle key presses
*/

func (t *TextAreaSelection) handleCopySelection() {
	t.pushUndo()
	if t.selectionStart == t.selectionEnd {
		// No selection to copy
		return
	}

	minPos, maxPos := t.getSelectionBounds()
	selectedText := t.text[minPos:maxPos]

	// Write to clipboard using golang-design/clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(selectedText))
	if err != nil {
		fmt.Println("handleCopySelection - Error writing to clipboard:", err)
	}
}

// handleCutSelection copies the selected text to the OS clipboard and removes it from the text area
func (t *TextAreaSelection) handleCutSelection() {
	t.pushUndo()
	if t.selectionStart == t.selectionEnd {
		// No selection to cut
		return
	}

	minPos, maxPos := t.getSelectionBounds()
	selectedText := t.text[minPos:maxPos]

	// Write to clipboard using golang-design/clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(selectedText))
	if err != nil {
		fmt.Println("handleCutSelection - Error writing to clipboard:", err)
	}

	// Remove the selected text from the text area
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.setCursorPos(minPos)
	t.clearSelection()
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
		t.setCursorPos(minPos + len(clipboardText))
	} else {
		// Insert clipboard text at cursor position
		t.text = t.text[:t.cursorPos] + clipboardText + t.text[t.cursorPos:]
		t.setCursorPos(t.cursorPos + len(clipboardText))
	}

	t.clearSelection()
}

func (t *TextAreaSelection) handleBackspace() {
	if t.selectionStart != t.selectionEnd {
		t.pushUndo()
		t.deleteSelection()
	} else if t.cursorPos > 0 {
		t.pushUndo()
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.setCursorPos(t.cursorPos - 1)
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleDelete() {
	if t.selectionStart != t.selectionEnd {
		t.pushUndo()
		t.deleteSelection()
	} else if t.cursorPos < len(t.text) {
		t.pushUndo()
		t.text = t.text[:t.cursorPos] + t.text[t.cursorPos+1:]
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleCtrlBackspace() {
	t.pushUndo()
	if t.selectionStart != t.selectionEnd {
		t.deleteSelection()
	} else {
		newPos := t.moveToWordStart(t.cursorPos)
		t.text = t.text[:newPos] + t.text[t.cursorPos:]
		t.setCursorPos(newPos)
		t.clearSelection()
	}
}

func (t *TextAreaSelection) handleCtrlDelete() {
	t.pushUndo()
	if t.selectionStart != t.selectionEnd {
		t.deleteSelection()
	} else {
		newPos := t.moveToWordEnd(t.cursorPos)
		// Prevent deleting the newline if cursor is at the end of a line
		if newPos > t.cursorPos && t.text[newPos-1] == '\n' {
			newPos--
		}
		t.text = t.text[:t.cursorPos] + t.text[newPos:]
		t.clearSelection()
	}
}

// handleTab processes the Tab key press
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

// handleEnter processes the Enter key press
func (t *TextAreaSelection) handleEnter() {
	t.pushUndo()
	t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
	t.cursorPos++
	t.clearSelection()
}

// handleLeftArrow processes the Left Arrow key press
func (t *TextAreaSelection) handleLeftArrow() {
	t.pushUndo()
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		t.updateSelectionWithShiftKey(-1)
	} else {
		if t.selectionStart != t.selectionEnd {
			t.clearSelection()
		} else {
			if t.cursorPos > 0 {
				t.cursorPos--
				t.clearSelection()
			}
		}
	}
}

// handleRightArrow processes the Right Arrow key press
func (t *TextAreaSelection) handleRightArrow() {
	t.pushUndo()
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		t.updateSelectionWithShiftKey(1)
	} else {
		if t.selectionStart != t.selectionEnd {
			t.clearSelection()
		} else {
			if t.cursorPos < len(t.text) {
				t.cursorPos++
				t.clearSelection()
			}
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
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		t.updateSelection(newPos)
	} else {
		t.clearSelection()
		t.setCursorPos(newPos)
	}
}

func (t *TextAreaSelection) handleEnd() {
	t.pushUndo()
	line, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	newPos := t.getCharPosFromLineAndCol(line, len(lines[line]))
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		t.updateSelection(newPos)
	} else {
		t.clearSelection()
		t.setCursorPos(newPos)
	}
}

func (t *TextAreaSelection) handleCtrlLeftArrow() {
	t.pushUndo()
	newPos := t.moveToWordStart(t.cursorPos)
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		t.updateSelection(newPos)
	} else {
		t.clearSelection()
		t.setCursorPos(newPos)
	}
}

func (t *TextAreaSelection) handleCtrlRightArrow() {
	t.pushUndo()
	newPos := t.moveToWordEnd(t.cursorPos)
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		t.updateSelection(newPos)
	} else {
		t.clearSelection()
		t.setCursorPos(newPos)
	}
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

// internals/widgets/textareaSelection01.go
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
// internals/widgets/textareaSelection01.go
func (t *TextAreaSelection) handleShiftHome() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	newPos := t.getCharPosFromLineAndCol(currentLine, 0)
	t.updateSelection(newPos)
}

func (t *TextAreaSelection) handleShiftEnd() {
	t.pushUndo()
	currentLine, _ := t.getCursorLineAndColForPos(t.cursorPos)
	lines := strings.Split(t.text, "\n")
	newPos := t.getCharPosFromLineAndCol(currentLine, len(lines[currentLine]))
	t.updateSelection(newPos)
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

// handleUndo processes the Undo action
// internals/widgets/textareaSelection01.go
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
		t.counter = 0 // Reset blink counter
	}
}
