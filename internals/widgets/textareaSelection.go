package widgets

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type TextAreaSelection struct {
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
}

func NewTextAreaSelection(x, y, w, h, maxLines int) *TextAreaSelection {
	return &TextAreaSelection{
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		maxLines:        maxLines,
		cursorBlinkRate: 30,
		tabWidth:        4,
		lineHeight:      20,
		font:            basicfont.Face7x13,
	}
}

func (t *TextAreaSelection) Draw(screen *ebiten.Image) {

	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)

	lines := strings.Split(t.text, "\n")
	yOffset := t.y

	for i, line := range lines {
		if i >= t.maxLines {
			break
		}

		lineText := line
		lineX := t.x
		lineY := yOffset + t.lineHeight/2

		if t.selectionStart != t.selectionEnd {
			startLine, startCol := t.getCursorLineAndColForPos(t.selectionStart)
			endLine, endCol := t.getCursorLineAndColForPos(t.selectionEnd)

			if startLine == endLine {

				if startCol > len(line) {
					startCol = len(line)
				}
				if endCol > len(line) {
					endCol = len(line)
				}
				startX := t.x + t.textWidth(line[:startCol])
				endX := t.x + t.textWidth(line[:endCol])

				vector.DrawFilledRect(screen, float32(startX), float32(yOffset), float32(endX-startX), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
			} else {

				if startCol > len(line) {
					startCol = len(line)
				}
				startX := t.x + t.textWidth(line[:startCol])
				vector.DrawFilledRect(screen, float32(startX), float32(yOffset), float32(t.x+t.w-startX), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)

				for j := startLine + 1; j < endLine; j++ {
					if j >= len(lines) {
						break
					}
					vector.DrawFilledRect(screen, float32(t.x), float32(yOffset+t.lineHeight), float32(t.w), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
					yOffset += t.lineHeight
				}

				if endCol > len(line) {
					endCol = len(line)
				}
				endX := t.x + t.textWidth(line[:endCol])
				vector.DrawFilledRect(screen, float32(t.x), float32(yOffset), float32(endX-t.x), float32(t.lineHeight), color.RGBA{0, 0, 255, 128}, true)
			}
		}

		text.Draw(screen, lineText, t.font, lineX, lineY+t.lineHeight/2, color.Black)

		yOffset += t.lineHeight
	}

	if t.hasFocus {
		cursorLine, cursorCol := t.getCursorLineAndCol()
		cursorX := t.x + t.textWidth(lines[cursorLine][:cursorCol])
		cursorY := t.y + cursorLine*t.lineHeight + t.lineHeight/2

		if t.counter%(t.cursorBlinkRate*2) < t.cursorBlinkRate {
			vector.DrawFilledRect(screen, float32(cursorX), float32(cursorY-t.lineHeight/2), 2, float32(t.lineHeight), color.RGBA{0, 0, 0, 255}, true)
		}
	}
}

func (t *TextAreaSelection) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			t.hasFocus = true
			t.startSelectionAtPosition(x, y)
		} else {
			t.hasFocus = false
		}
	}

	if t.hasFocus && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		t.isSelecting = true
	}
	if t.isSelecting && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		t.updateSelection(x, y)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		t.isSelecting = false
	}

	if t.hasFocus {
		t.handleKeyboardInput()
	}

	t.counter++
	return nil
}

func (t *TextAreaSelection) handleKeyboardInput() {

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(t.text) > 0 && t.cursorPos > 0 {
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.cursorPos--
		t.clearSelection()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if t.isSelecting {
			t.indentSelection()
		} else {
			t.text = t.text[:t.cursorPos] + strings.Repeat(" ", t.tabWidth) + t.text[t.cursorPos:]
			t.cursorPos += t.tabWidth
		}
		t.clearSelection()
	}

	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			t.text = t.text[:t.cursorPos] + string(char) + t.text[t.cursorPos:]
			t.cursorPos++
			t.clearSelection()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
		t.cursorPos++
		t.clearSelection()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && t.cursorPos > 0 {
		t.cursorPos--
		t.updateSelectionWithShiftKey(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && t.cursorPos < len(t.text) {
		t.cursorPos++
		t.updateSelectionWithShiftKey(1)
	}

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
	t.cursorPos = t.selectionEnd + len(indent)
}

func (t *TextAreaSelection) startSelectionAtPosition(x, y int) {
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
	for i := range lineText {
		charWidth := t.textWidth(string(lineText[i]))
		if col < colIndex+charWidth/2 {
			break
		}
		colIndex += charWidth
	}

	charPos := t.getCharPosFromLineAndCol(line, colIndex)
	t.cursorPos = charPos
	t.selectionStart = charPos
	t.selectionEnd = charPos
	t.isSelecting = true
}

func (t *TextAreaSelection) updateSelection(x, y int) {
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
	for i := range lineText {
		charWidth := t.textWidth(string(lineText[i]))
		if col < colIndex+charWidth/2 {
			break
		}
		colIndex += charWidth
	}

	charPos := t.getCharPosFromLineAndCol(line, colIndex)
	t.cursorPos = charPos
	t.selectionEnd = charPos
}

func (t *TextAreaSelection) updateSelectionWithShiftKey(offset int) {
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

func (t *TextAreaSelection) clearSelection() {
	t.selectionStart = t.cursorPos
	t.selectionEnd = t.cursorPos
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
