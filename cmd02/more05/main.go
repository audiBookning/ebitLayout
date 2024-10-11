package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	InitializeAssets()
	game := NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Responsive Layout with Ebitengine")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Asset paths
var (
	filePathTxt          string
	Assets_Relative_Path = "../../"
)

func InitializeAssets() {
	_, filePathTxt, _, _ = runtime.Caller(0)
}

func GetFilePath(fileName string) string {
	dir := filepath.Dir(filePathTxt)
	return filepath.Join(dir, Assets_Relative_Path, fileName)
}
