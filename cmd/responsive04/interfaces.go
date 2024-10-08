package main

import "github.com/hajimehoshi/ebiten/v2"

// Page defines the common interface for all pages.
type Page interface {
	Update() error
	Draw(screen *ebiten.Image)
	HandleInput(x, y int)
}
