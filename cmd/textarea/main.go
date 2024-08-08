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
)

type Game struct {
	text      string
	hasFocus  bool
	cursorPos int
	counter   int
}

func (g *Game) Update() error {
	// Update focus state based on mouse input
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= textAreaX && x <= textAreaX+textAreaW && y >= textAreaY && y <= textAreaY+textAreaH {
			g.hasFocus = true
		} else {
			g.hasFocus = false
		}
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
	}

	// Handle character input
	for _, char := range ebiten.InputChars() {
		if char != '\n' && char != '\r' {
			g.text = g.text[:g.cursorPos] + string(char) + g.text[g.cursorPos:]
			g.cursorPos++
		}
	}

	// Handle enter key for new line
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.text = g.text[:g.cursorPos] + "\n" + g.text[g.cursorPos:]
		g.cursorPos++
	}

	// Handle arrow keys for cursor movement
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.cursorPos > 0 {
		g.cursorPos--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.cursorPos < len(g.text) {
		g.cursorPos++
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the textarea background
	vector.DrawFilledRect(screen, float32(textAreaX), float32(textAreaY), float32(textAreaW), float32(textAreaH), color.RGBA{200, 200, 200, 255}, true)

	// Draw the text
	lines := strings.Split(g.text, "\n")
	startY := textAreaY + 20
	for i, line := range lines {
		if i >= maxLines {
			break
		}
		text.Draw(screen, line, basicfont.Face7x13, textAreaX, startY+i*20, color.Black)
	}

	// Draw the cursor if focused and blink state is on
	if g.hasFocus && g.counter/cursorBlinkRate%2 == 0 {
		// Calculate the cursor position in pixels
		line, col := g.getCursorLineAndCol()
		cursorX := textAreaX + textWidth(basicfont.Face7x13, lines[line][:col])
		cursorY := startY + line*20 - 10
		vector.DrawFilledRect(screen, float32(cursorX), float32(cursorY), 2, 20, color.RGBA{0, 0, 0, 255}, true) // Cursor color
	}
}

func (g *Game) getCursorLineAndCol() (int, int) {
	line, col := 0, 0
	for i, char := range g.text[:g.cursorPos] {
		if char == '\n' {
			line++
			col = 0
		} else {
			col++
		}
		if i+1 == g.cursorPos {
			break
		}
	}
	return line, col
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

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Enhanced Text Input Example")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
