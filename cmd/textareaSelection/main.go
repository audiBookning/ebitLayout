package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	textAreaX    = 50
	textAreaY    = 50
	textAreaW    = 540
	textAreaH    = 300
)

type Game struct {
	textarea *widgets.TextAreaSelection
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.textarea.Draw(screen)
}

func (g *Game) Update() error {
	return g.textarea.Update()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func getFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

var filePathTxt string

const Assets_Relative_Path = "../../"

func main() {
	_, filePathTxt, _, _ = runtime.Caller(0)
	textPath := getFilePath("internals/widgets/textStart.txt")

	// Read the text from the textStart.txt file
	textStart, err := os.ReadFile(textPath)
	if err != nil {
		log.Fatalf("Failed to read textStart.txt: %v", err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text Input with Selection Example")
	game := &Game{
		textarea: widgets.NewTextAreaSelection(textAreaX, textAreaY, textAreaW, textAreaH, 14, string(textStart)),
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
