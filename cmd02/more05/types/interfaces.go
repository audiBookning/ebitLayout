package types

import "github.com/hajimehoshi/ebiten/v2"

type Page interface {
	Update() error
	Draw(screen *ebiten.Image)
	DrawBackGround(screen *ebiten.Image)
	HandleInput(x, y int)
	ResetFieldStates()
	Layout(outsideWidth, outsideHeight int) (int, int)
}

type Element interface {
	GetPosition() Position
	SetPosition(Position)
	Draw(screen *ebiten.Image)
	Update()
	IsClicked(x, y int) bool
	HandleClick()
	ResetState()
	GetSize() (int, int)
}
