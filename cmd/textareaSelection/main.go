package main

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

const (
	screenWidth     = 640
	screenHeight    = 480
	textAreaX       = 50
	textAreaY       = 50
	textAreaW       = 540
	textAreaH       = 300
	maxLines        = 10
	cursorBlinkRate = 30 // Frames for cursor blink
	tabWidth        = 4  // Number of spaces for a tab
)

type Game struct {
	text           string
	hasFocus       bool
	cursorPos      int
	counter        int
	selectionStart int
	selectionEnd   int
	isSelecting    bool
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the textarea background
	vector.DrawFilledRect(screen, float32(textAreaX), float32(textAreaY), float32(textAreaW), float32(textAreaH), color.RGBA{200, 200, 200, 255}, true)

	// Draw the text with selection
	lines := strings.Split(g.text, "\n")
	yOffset := textAreaY // Initial y offset for text drawing
	lineHeight := 20     // Height of each line, adjust as needed

	for i, line := range lines {
		if i >= maxLines {
			break
		}

		lineText := line
		lineX := textAreaX
		lineY := yOffset + lineHeight/2 // Center vertically

		// Calculate text width for the current line
		//currentTextWidth := textWidth(basicfont.Face7x13, lineText)

		// Draw the text selection
		if g.selectionStart != g.selectionEnd {
			startLine, startCol := g.getCursorLineAndColForPos(g.selectionStart)
			endLine, endCol := g.getCursorLineAndColForPos(g.selectionEnd)

			if startLine == endLine {
				// Single line selection
				if startCol > len(line) {
					startCol = len(line)
				}
				if endCol > len(line) {
					endCol = len(line)
				}
				startX := textAreaX + textWidth(basicfont.Face7x13, line[:startCol])
				endX := textAreaX + textWidth(basicfont.Face7x13, line[:endCol])

				// Draw the selection rectangle for a single line
				vector.DrawFilledRect(screen, float32(startX), float32(yOffset), float32(endX-startX), float32(lineHeight), color.RGBA{0, 0, 255, 128}, true)
			} else {
				// Multi-line selection

				// Handle the first line
				if startCol > len(line) {
					startCol = len(line)
				}
				startX := textAreaX + textWidth(basicfont.Face7x13, line[:startCol])
				vector.DrawFilledRect(screen, float32(startX), float32(yOffset), float32(textAreaX+textAreaW-startX), float32(lineHeight), color.RGBA{0, 0, 255, 128}, true)

				// Handle middle lines
				for j := startLine + 1; j < endLine; j++ {
					if j >= len(lines) {
						break
					}
					vector.DrawFilledRect(screen, float32(textAreaX), float32(yOffset+lineHeight), float32(textAreaW), float32(lineHeight), color.RGBA{0, 0, 255, 128}, true)
					yOffset += lineHeight
				}

				// Handle the last line
				if endCol > len(line) {
					endCol = len(line)
				}
				endX := textAreaX + textWidth(basicfont.Face7x13, line[:endCol])
				vector.DrawFilledRect(screen, float32(textAreaX), float32(yOffset), float32(endX-textAreaX), float32(lineHeight), color.RGBA{0, 0, 255, 128}, true)
			}
		}

		// Draw the text itself
		text.Draw(screen, lineText, basicfont.Face7x13, lineX, lineY+lineHeight/2, color.Black)

		// Move yOffset down for the next line
		yOffset += lineHeight
	}

	// Draw the cursor if the text area has focus
	if g.hasFocus {
		cursorLine, cursorCol := g.getCursorLineAndCol()
		cursorX := textAreaX + textWidth(basicfont.Face7x13, lines[cursorLine][:cursorCol])
		cursorY := textAreaY + cursorLine*lineHeight + lineHeight/2

		// Draw the cursor (flashing rectangle)
		if g.counter%(cursorBlinkRate*2) < cursorBlinkRate {
			vector.DrawFilledRect(screen, float32(cursorX), float32(cursorY-lineHeight/2), 2, float32(lineHeight), color.RGBA{0, 0, 0, 255}, true)
		}
	}
}

func (g *Game) Update() error {
	// Update focus state based on mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= textAreaX && x <= textAreaX+textAreaW && y >= textAreaY && y <= textAreaY+textAreaH {
			g.hasFocus = true
			g.startSelectionAtPosition(x, y)
		} else {
			g.hasFocus = false
		}
	}

	// Continue selection if mouse is dragged while pressed
	if g.hasFocus && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.isSelecting = true
	}
	if g.isSelecting && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.updateSelection(x, y)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.isSelecting = false
	}

	// If the text area has focus, handle keyboard input
	if g.hasFocus {
		g.handleKeyboardInput()
	}

	g.counter++
	return nil
}

func (g *Game) handleKeyboardInput() {
	// Handle backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(g.text) > 0 && g.cursorPos > 0 {
		g.text = g.text[:g.cursorPos-1] + g.text[g.cursorPos:]
		g.cursorPos--
		g.clearSelection()
	}

	// Handle Tab key
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if g.isSelecting {
			g.indentSelection()
		} else {
			g.text = g.text[:g.cursorPos] + strings.Repeat(" ", tabWidth) + g.text[g.cursorPos:]
			g.cursorPos += tabWidth
		}
		g.clearSelection()
	}

	// Handle character input
	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			g.text = g.text[:g.cursorPos] + string(char) + g.text[g.cursorPos:]
			g.cursorPos++
			g.clearSelection()
		}
	}

	// Handle enter key for new line
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.text = g.text[:g.cursorPos] + "\n" + g.text[g.cursorPos:]
		g.cursorPos++
		g.clearSelection()
	}

	// Handle arrow keys for cursor movement
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.cursorPos > 0 {
		g.cursorPos--
		g.updateSelectionWithShiftKey(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.cursorPos < len(g.text) {
		g.cursorPos++
		g.updateSelectionWithShiftKey(1)
	}
	// Ensure cursor position is within text bounds
	g.cursorPos = clamp(g.cursorPos, 0, len(g.text))
}

func (g *Game) indentSelection() {
	lines := strings.Split(g.text, "\n")
	startLine, _ := g.getCursorLineAndColForPos(g.selectionStart)
	endLine, _ := g.getCursorLineAndColForPos(g.selectionEnd)
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}
	if startLine >= len(lines) {
		startLine = len(lines) - 1
	}

	var indent string
	for i := 0; i < tabWidth; i++ {
		indent += " "
	}

	for i := startLine; i <= endLine; i++ {
		lines[i] = indent + lines[i]
	}

	g.text = strings.Join(lines, "\n")
	g.cursorPos = g.selectionEnd + len(indent) // Adjust cursor position after indentation
}

func (g *Game) startSelectionAtPosition(x, y int) {
	lineHeight := 20
	line := (y - textAreaY) / lineHeight
	col := x - textAreaX

	lines := strings.Split(g.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	if line < 0 {
		line = 0
	}

	lineText := lines[line]
	colIndex := 0
	for i := range lineText {
		charWidth := textWidth(basicfont.Face7x13, string(lineText[i]))
		if col < colIndex+charWidth/2 {
			break
		}
		colIndex += charWidth
	}

	charPos := g.getCharPosFromLineAndCol(line, colIndex)
	g.cursorPos = charPos
	g.selectionStart = charPos
	g.selectionEnd = charPos
	g.isSelecting = true
}

func (g *Game) updateSelection(x, y int) {
	lineHeight := 20
	line := (y - textAreaY) / lineHeight
	col := x - textAreaX

	lines := strings.Split(g.text, "\n")
	if line >= len(lines) {
		line = len(lines) - 1
	}
	if line < 0 {
		line = 0
	}

	lineText := lines[line]
	colIndex := 0
	for i := range lineText {
		charWidth := textWidth(basicfont.Face7x13, string(lineText[i]))
		if col < colIndex+charWidth/2 {
			break
		}
		colIndex += charWidth
	}

	charPos := g.getCharPosFromLineAndCol(line, colIndex)
	g.cursorPos = charPos
	g.selectionEnd = charPos
}

func (g *Game) updateSelectionWithShiftKey(offset int) {
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		newCursorPos := clamp(g.cursorPos+offset, 0, len(g.text))

		if newCursorPos > g.cursorPos {
			// Cursor moving right
			if newCursorPos > g.selectionEnd {
				g.selectionEnd = newCursorPos
			} else {
				g.selectionStart = newCursorPos
			}
		} else {
			// Cursor moving left
			if newCursorPos < g.selectionStart {
				g.selectionStart = newCursorPos
			} else {
				g.selectionEnd = newCursorPos
			}
		}
		g.cursorPos = newCursorPos
	} else {
		// If Shift is not pressed, clear the selection
		g.clearSelection()
	}
}

func (g *Game) clearSelection() {
	g.selectionStart = g.cursorPos
	g.selectionEnd = g.cursorPos
}

func (g *Game) getCharPosFromLineAndCol(line, col int) int {
	lines := strings.Split(g.text, "\n")
	charPos := 0
	for i := 0; i < line; i++ {
		charPos += len(lines[i]) + 1 // +1 for newline character
	}
	charPos += col
	return charPos
}

func (g *Game) getCursorLineAndColForPos(pos int) (int, int) {
	lines := strings.Split(g.text, "\n")
	charCount := 0
	for i, line := range lines {
		if charCount+len(line)+1 > pos { // +1 for newline character
			return i, pos - charCount
		}
		charCount += len(line) + 1 // +1 for newline character
	}
	return len(lines) - 1, len(lines[len(lines)-1])
}

func (g *Game) getCursorLineAndCol() (int, int) {
	return g.getCursorLineAndColForPos(g.cursorPos)
}

func textWidth(face font.Face, str string) int {
	width := 0
	for _, x := range str {
		awidth, _ := face.GlyphAdvance(x)
		width += int(awidth >> 6)
	}
	return width
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text Input with Selection Example")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
