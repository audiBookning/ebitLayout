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

type TextArea struct {
	text       string
	hasFocus   bool
	cursorPos  int
	counter    int
	x, y, w, h int
	maxLines   int
	blinkRate  int
	font       font.Face
}

func NewTextArea(x, y, w, h, maxLines int) *TextArea {
	return &TextArea{
		x:         x,
		y:         y,
		w:         w,
		h:         h,
		maxLines:  maxLines,
		blinkRate: 30,
		font:      basicfont.Face7x13,
	}
}

func (t *TextArea) Update() error {
	// Update focus state based on mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= t.x && x <= t.x+t.w && y >= t.y && y <= t.y+t.h {
			t.hasFocus = true
		} else {
			t.hasFocus = false
		}
	}

	// If the text area has focus, handle keyboard input
	if t.hasFocus {
		t.handleKeyboardInput()
	}

	t.counter++
	return nil
}

func (t *TextArea) handleKeyboardInput() {
	// Handle backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(t.text) > 0 && t.cursorPos > 0 {
		t.text = t.text[:t.cursorPos-1] + t.text[t.cursorPos:]
		t.cursorPos--
	}

	// Handle character input
	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			t.text = t.text[:t.cursorPos] + string(char) + t.text[t.cursorPos:]
			t.cursorPos++
		}
	}

	// Handle enter key for new line
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		t.text = t.text[:t.cursorPos] + "\n" + t.text[t.cursorPos:]
		t.cursorPos++
	}

	// Handle arrow keys for cursor movement
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && t.cursorPos > 0 {
		t.cursorPos--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && t.cursorPos < len(t.text) {
		t.cursorPos++
	}
}

func (t *TextArea) Draw(screen *ebiten.Image) {
	// Draw the textarea background
	vector.DrawFilledRect(screen, float32(t.x), float32(t.y), float32(t.w), float32(t.h), color.RGBA{200, 200, 200, 255}, true)

	// Draw the text
	lines := strings.Split(t.text, "\n")
	startY := t.y + 20
	for i, line := range lines {
		if i >= t.maxLines {
			break
		}
		text.Draw(screen, line, t.font, t.x, startY+i*20, color.Black)
	}

	// Draw the cursor if focused and blink state is on
	if t.hasFocus && t.counter/t.blinkRate%2 == 0 {
		// Calculate the cursor position in pixels
		line, col := t.getCursorLineAndCol()
		cursorX := t.x + t.textWidth(lines[line][:col])
		cursorY := startY + line*20 - 10
		vector.DrawFilledRect(screen, float32(cursorX), float32(cursorY), 2, 20, color.RGBA{0, 0, 0, 255}, true)
	}
}

func (t *TextArea) getCursorLineAndCol() (int, int) {
	line, col := 0, 0
	for i, char := range t.text[:t.cursorPos] {
		if char == '\n' {
			line++
			col = 0
		} else {
			col++
		}
		if i+1 == t.cursorPos {
			break
		}
	}
	return line, col
}

func (t *TextArea) textWidth(str string) int {
	width := 0
	for _, x := range str {
		awidth, _ := t.font.GlyphAdvance(x)
		width += int(awidth >> 6)
	}
	return width
}
