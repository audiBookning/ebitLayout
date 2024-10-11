package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"example.com/menu/internals/layout"
	"example.com/menu/internals/textwrapper"
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	flexbox              *layout.FlexBox
	initialWindowWidth   int
	initialWindowHeight  int
	initialFontSize      float64
	initialElementHeight int
}

func Create12ColumnLayout(
	textWrapper *textwrapper.TextWrapper,
	initialElementHeight int) *layout.FlexBox {

	columns := make([]layout.Element, 8)
	for i := range columns {
		columns[i] = layout.Element{
			Width:       0,
			Height:      initialElementHeight,
			Color:       color.RGBA{uint8(20 * i), uint8(255 - 20*i), uint8(128 + 10*i), 255},
			Flex:        1,
			Text:        fmt.Sprintf("%d", i+1),
			TextWrapper: textWrapper,
			TextSize:    24,
		}
	}

	return &layout.FlexBox{
		Elements:       columns,
		Direction:      "row",
		JustifyContent: "space-between",
		AlignItems:     "stretch",
	}
}

var filePathTxt string
var Assets_Relative_Path = "../../"

func GetFilePath(fileName string) string {

	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}

func NewGame() *Game {
	_, filePathTxt, _, _ = runtime.Caller(0)
	fontpath := GetFilePath("assets/fonts/roboto_regularTTF.ttf")
	log.Println("fontpath", fontpath)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("current working directory", dir)
	initialWindowWidth, initialWindowHeight := 1200, 600
	initialFontSize := 24.0
	initialElementHeight := 100

	textWrapper, err := textwrapper.NewTextWrapper(
		fontpath, initialFontSize, false)
	if err != nil {
		log.Fatal(err)
	}
	textWrapper.Color = color.White

	flexbox := Create12ColumnLayout(
		textWrapper, initialElementHeight)
	return &Game{
		flexbox:              flexbox,
		initialWindowWidth:   initialWindowWidth,
		initialWindowHeight:  initialWindowHeight,
		initialFontSize:      initialFontSize,
		initialElementHeight: initialElementHeight,
	}
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	scaleFactor := float64(outsideWidth) / float64(g.initialWindowWidth)

	newFontSize := g.initialFontSize * scaleFactor
	newElementHeight := int(float64(g.initialElementHeight) * scaleFactor)

	for i := range g.flexbox.Elements {
		g.flexbox.Elements[i].TextSize = newFontSize
		g.flexbox.Elements[i].Height = newElementHeight
	}

	g.flexbox.Layout(outsideWidth, outsideHeight)

	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.flexbox.Draw(screen)
}

func main() {
	ebiten.SetWindowSize(1200, 600)
	ebiten.SetWindowTitle("12-Column Layout Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
