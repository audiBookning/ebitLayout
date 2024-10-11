package types

import "github.com/hajimehoshi/ebiten/v2"

type Page interface {
	Update() error
	Draw(screen *ebiten.Image)
	DrawBackGround(screen *ebiten.Image)
	HandleInput(x, y int)
	ResetButtonStates()
	Layout(outsideWidth, outsideHeight int) (int, int)
}
