package main

import (
	"image/color"
	"log"
	"os"
	"path/filepath"
	"runtime"

	//"example.com/menu/internals/textwrapper02"
	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	textAreaX    = 50
	textAreaY    = 50
	textAreaW    = 500
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
	textPath := getFilePath("cmd/textareaSelection/textStart.txt")

	// Read the text from the textStart.txt file
	textStart, err := os.ReadFile(textPath)
	if err != nil {
		log.Fatalf("Failed to read textStart.txt: %v", err)
	}

	//fontPath := getFilePath("assets/fonts/roboto_regularTTF.ttf")
	fontPath := getFilePath("assets/fonts/Anonymous_Pro.ttf")
	fontSize := 40.0

	//textWrapper := textwrapper.NewTextWrapper(basicfont.Face7x13, fontSize, color.Black)
	textWrapper, err := textwrapper.NewTextWrapper(fontPath, fontSize, false)
	//textWrapper, err := textwrapper02.NewTextWrapper(fontPath, fontSize)
	if err != nil {
		log.Fatalf("Failed to create text wrapper: %v", err)
	}
	textWrapper.Color = color.Black

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text Input with Selection Example")
	game := &Game{
		textarea: widgets.NewTextAreaSelection(textWrapper, textAreaX, textAreaY, textAreaW, textAreaH, string(textStart)),
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
