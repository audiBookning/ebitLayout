package dd

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"golang.design/x/clipboard" // Updated clipboard package
)

const (
	// keyRepeatInitialDelay02 is the number of frames to wait before repeating the key action
	keyRepeatInitialDelay02 = 30 // Approximately 0.5 seconds at 60 TPS

	// keyRepeatInterval02 is the number of frames between each repeated action
	keyRepeatInterval02 = 5 // Approximately 10 times per second at 60 TPS
)

// KeyState02 tracks the repeat state of a specific key
type KeyState02 struct {
	InitialPress    bool // Indicates if the initial press has been handled
	FramesHeld      int  // Number of frames the key has been held down
	FramesUntilNext int  // Frames remaining until the next action
}

type TextAreaSelection02 struct {
	text            string
	hasFocus        bool
	cursorPos       int
	counter         int
	selectionStart  int
	selectionEnd    int
	isSelecting     bool
	x, y, w, h      int
	maxLines        int
	cursorBlinkRate int
	tabWidth        int
	lineHeight      int
	font            font.Face
	heldKeys        map[ebiten.Key]*KeyState02
}

func NewTextAreaSelection02(x, y, w, h, maxLines int) *TextAreaSelection02 {
	err := clipboard.Init()
	if err != nil {
		return nil
	}

	return &TextAreaSelection02{
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		maxLines:        maxLines,
		cursorBlinkRate: 30,
		tabWidth:        4,
		lineHeight:      20,
		font:            basicfont.Face7x13,
		heldKeys:        make(map[ebiten.Key]*KeyState02),
	}
}

func (t *TextAreaSelection02) Draw(screen *ebiten.Image) {

	// Draw the background of the text area
	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)

	lines := strings.Split(t.text, "\n")
	yOffset := t.y

	// Retrieve normalized selection bounds
	minPos, maxPos := t.getSelectionBounds()

	for i, line := range lines {
		if i >= t.maxLines {
			break
		}

		lineText := line
		lineX := t.x
		lineY := yOffset + t.lineHeight/2

		// Draw selection if active and within this line
		if minPos != maxPos {
			startLine, startCol := t.getCursorLineAndColForPos(minPos)
			endLine, endCol := t.getCursorLineAndColForPos(maxPos)

			if i < startLine || i > endLine {
				// No selection in this line
			} else {
				// Determine selection bounds within the current line
				var selStart, selEnd int
				if i == startLine {
					selStart = startCol
				} else {
					selStart = 0
				}
				if i == endLine {
					selEnd = endCol
				} else {
					selEnd = len(line)
				}

				// Calculate x positions based on character widths
				selectionXStart := t.x + t.textWidth(lineText[:selStart])
				selectionXEnd := t.x + t.textWidth(lineText[:selEnd])

				// Draw the selection rectangle
				vector.DrawFilledRect(screen, float32(selectionXStart), float32(yOffset), float32(selectionXEnd-selectionXStart), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
			}
		}

		// Draw the actual text
		text.Draw(screen, lineText, t.font, lineX, lineY+t.lineHeight/2, color.Black)

		yOffset += t.lineHeight
	}

	// Draw the cursor if the text area has focus
	if t.hasFocus {
		cursorLine, cursorCol := t.getCursorLineAndCol()
		cursorX := t.x + t.textWidth(lines[cursorLine][:cursorCol])
		cursorY := t.y + cursorLine*t.lineHeight

		// Render the blinking cursor
		if t.counter%(t.cursorBlinkRate*2) < t.cursorBlinkRate {
			vector.DrawFilledRect(screen, float32(cursorX), float32(cursorY), 2, float32(t.lineHeight), color.RGBA{0, 0, 0, 255}, true)
		}
	}
}

func (t *TextAreaSelection02) Update() error {
	// Handle mouse button just pressed (mouse down)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			t.hasFocus = true
			charPos := t.getCharPosFromPosition(x, y)
			t.cursorPos = charPos
			t.selectionStart = charPos
			t.selectionEnd = charPos
			// Initially, no selection is active
			t.isSelecting = false
		} else {
			t.hasFocus = false
			t.isSelecting = false
		}
	}

	// Handle mouse movement while the left button is pressed (dragging)
	if t.hasFocus && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		charPos := t.getCharPosFromPosition(x, y)
		if !t.isSelecting {
			// Start selection on first movement after click
			t.isSelecting = true
			t.selectionStart = t.cursorPos
		}
		t.selectionEnd = charPos
		t.cursorPos = charPos
	}

	// Handle mouse button release (mouse up)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if t.isSelecting {
			// Finalize selection on mouse release
			t.isSelecting = false
		}
	}

	// Handle keyboard input when focused
	if t.hasFocus {
		t.handleKeyboardInput()
	}

	t.counter++
	return nil
}

func (t *TextAreaSelection02) handleKeyboardInput() {
	// Define the keys that support repeat (Backspace and Delete)
	repeatKeys := []ebiten.Key{ebiten.KeyBackspace, ebiten.KeyDelete}

	// Handle Backspace and Delete with support for repeat
	for _, key := range repeatKeys {
		if ebiten.IsKeyPressed(key) {
			// Initialize key state if not present
			if _, exists := t.heldKeys[key]; !exists {
				t.heldKeys[key] = &KeyState02{
					InitialPress:    true,
					FramesHeld:      0,
					FramesUntilNext: keyRepeatInitialDelay02,
				}
			}

			keyState := t.heldKeys[key]

			if keyState.InitialPress {
				// Handle the initial key press
				if key == ebiten.KeyBackspace {
					t.handleBackspace()
				} else if key == ebiten.KeyDelete {
					t.handleDelete()
				}
				keyState.InitialPress = false
			} else {
				// Increment frames held
				keyState.FramesHeld++

				// Check if it's time to repeat the action
				if keyState.FramesHeld >= keyState.FramesUntilNext {
					if key == ebiten.KeyBackspace {
						t.handleBackspace()
					} else if key == ebiten.KeyDelete {
						t.handleDelete()
					}
					// Reset frames until next action
					keyState.FramesUntilNext = keyRepeatInterval02
				}
			}
		} else {
			// Key is not pressed; remove from heldKeys
			delete(t.heldKeys, key)
		}
	}

	// Handle Tab
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if t.isSelecting {
			t.indentSelection()
		} else {
			t.text = t.text[:t.cursorPos] + strings.Repeat(" ", t.tabWidth) + t.text[t.cursorPos:]
			t.cursorPos += t.tabWidth
		}
		t.clearSelection()
	}

	// Handle character input
	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			t.text = t.text[:t.cursorPos] + string(char) + t.text[t.cursorPos:]
			t.cursorPos++
			t.clearSelection()
		}
	}

	// Handle Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
		t.cursorPos++
		t.clearSelection()
	}

	// Handle Left Arrow
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
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

	// Handle Right Arrow
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
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

	// Handle Copy (Ctrl+C), Cut (Ctrl+X), and Paste (Ctrl+V)
	if ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight) {
		// Copy (Ctrl+C)
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			t.copySelection()
		}

		// Cut (Ctrl+X)
		if inpututil.IsKeyJustPressed(ebiten.KeyX) {
			t.cutSelection()
		}

		// Paste (Ctrl+V)
		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			t.pasteClipboard()
		}
	}

	// Clamp cursor position to valid range
	t.cursorPos = clamp(t.cursorPos, 0, len(t.text))
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (t *TextAreaSelection02) indentSelection() {
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
	t.cursorPos = t.selectionEnd + len(indent)
}

func (t *TextAreaSelection02) updateSelectionWithShiftKey(offset int) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		newCursorPos := clamp(t.cursorPos+offset, 0, len(t.text))

		if newCursorPos > t.cursorPos {
			if newCursorPos > t.selectionEnd {
				t.selectionEnd = newCursorPos
			} else {
				t.selectionStart = newCursorPos
			}
		} else {
			if newCursorPos < t.selectionStart {
				t.selectionStart = newCursorPos
			} else {
				t.selectionEnd = newCursorPos
			}
		}
		t.cursorPos = newCursorPos
	} else {
		t.clearSelection()
	}
}

func (t *TextAreaSelection02) clearSelection() {
	t.selectionStart = t.cursorPos
	t.selectionEnd = t.cursorPos
}

func (t *TextAreaSelection02) getCharPosFromLineAndCol(line, col int) int {
	lines := strings.Split(t.text, "\n")
	charPos := 0
	for i := 0; i < line; i++ {
		charPos += len(lines[i]) + 1
	}
	charPos += col
	return charPos
}

func (t *TextAreaSelection02) getCursorLineAndColForPos(pos int) (int, int) {
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

func (t *TextAreaSelection02) getCursorLineAndCol() (int, int) {
	return t.getCursorLineAndColForPos(t.cursorPos)
}

func (t *TextAreaSelection02) textWidth(str string) int {
	width := 0
	for _, x := range str {
		awidth, _ := t.font.GlyphAdvance(x)
		width += int(awidth >> 6)
	}
	return width
}

func (t *TextAreaSelection02) getCharPosFromPosition(x, y int) int {
	line := (y - t.y) / t.lineHeight
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
func (t *TextAreaSelection02) getSelectionBounds() (int, int) {
	if t.selectionStart <= t.selectionEnd {
		return t.selectionStart, t.selectionEnd
	}
	return t.selectionEnd, t.selectionStart
}

// copySelection copies the selected text to the OS clipboard
func (t *TextAreaSelection02) copySelection() {
	if t.selectionStart == t.selectionEnd {
		// No selection to copy
		return
	}

	minPos, maxPos := t.getSelectionBounds()
	selectedText := t.text[minPos:maxPos]

	// Write to clipboard using golang-design/clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(selectedText))
	if err != nil {
		// Handle the error appropriately (e.g., log it)
		// For simplicity, we'll ignore it here
	}
}

// cutSelection copies the selected text to the OS clipboard and removes it from the text area
func (t *TextAreaSelection02) cutSelection() {
	if t.selectionStart == t.selectionEnd {
		// No selection to cut
		return
	}

	minPos, maxPos := t.getSelectionBounds()
	selectedText := t.text[minPos:maxPos]

	// Write to clipboard using golang-design/clipboard
	err := clipboard.Write(clipboard.FmtText, []byte(selectedText))
	if err != nil {
		// Handle the error appropriately (e.g., log it)
		// For simplicity, we'll ignore it here
	}

	// Remove the selected text from the text area
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.cursorPos = minPos
	t.clearSelection()
}

// pasteClipboard inserts text from the OS clipboard into the text area at the current cursor position
func (t *TextAreaSelection02) pasteClipboard() {
	clipboardBytes := clipboard.Read(clipboard.FmtText)

	clipboardText := string(clipboardBytes)

	if t.selectionStart != t.selectionEnd {
		// If there's an active selection, replace it with the clipboard text
		minPos, maxPos := t.getSelectionBounds()
		t.text = t.text[:minPos] + clipboardText + t.text[maxPos:]
		t.cursorPos = minPos + len(clipboardText)
	} else {
		// Insert the clipboard text at the current cursor position
		t.text = t.text[:t.cursorPos] + clipboardText + t.text[t.cursorPos:]
		t.cursorPos += len(clipboardText)
	}

	t.clearSelection()
}

// deleteSelection removes the currently selected text and updates the cursor position
func (t *TextAreaSelection02) deleteSelection() {
	minPos, maxPos := t.getSelectionBounds()
	t.text = t.text[:minPos] + t.text[maxPos:]
	t.cursorPos = minPos
	t.clearSelection()
}

// handleBackspace processes the Backspace key press
func (t *TextAreaSelection02) handleBackspace() {
	if t.selectionStart != t.selectionEnd {
		// **Active Selection**: Delete the entire selection
		t.deleteSelection()
	} else if t.cursorPos > 0 {
		// **No Selection**: Delete character before cursor
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.cursorPos--
		t.clearSelection()
	}
}

// handleDelete processes the Delete key press
func (t *TextAreaSelection02) handleDelete() {
	if t.selectionStart != t.selectionEnd {
		// **Active Selection**: Delete the entire selection
		t.deleteSelection()
	} else if t.cursorPos < len(t.text) {
		// **No Selection**: Delete character after cursor
		t.text = t.text[:t.cursorPos] + t.text[t.cursorPos+1:]
		t.clearSelection()
	}
}
