package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	textAreaX    = 50
	textAreaY    = 50
	textAreaW    = 540
	textAreaH    = 430
	fontSize     = 14
)

type Game struct {
	textarea *TextAreaSelection
}

func (g *Game) Update() error {
	g.textarea.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.textarea.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// getFilePath constructs the absolute path for the given file.
func getFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

var filePathTxt string

const Assets_Relative_Path = "../../"

func main() {

	_, filePathTxt, _, _ = runtime.Caller(0)
	textPath := getFilePath("cmd/textareaSelection/textStart.txt")

	// Read the initial text
	textBytes, err := os.ReadFile(textPath)
	if err != nil {
		log.Fatalf("Failed to read textStart.txt: %v", err)
	}
	startText := string(textBytes)

	// Specify the path to the font file
	fontPath := getFilePath("assets/fonts/roboto_regularTTF.ttf")
	fontSize := 14.0

	// Initialize TextWrapper
	textWrapper, err := NewTextWrapper(fontPath, fontSize)
	if err != nil {
		log.Fatalf("Failed to create TextWrapper: %v", err)
	}

	// Initialize TextAreaSelection
	textarea := NewTextAreaSelection(textWrapper, textAreaX, textAreaY, textAreaW, textAreaH, startText)

	// Initialize the game
	game := &Game{
		textarea: textarea,
	}

	// Set up Ebiten window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic TextAreaSelection Example")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
