package customPages

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Level02Page struct {
	title string
}

func NewLevel02Page() *Level02Page {
	return &Level02Page{
		title: "Level 02 - The Challenge",
	}
}

func (p *Level02Page) Update() error {

	return nil
}

func (p *Level02Page) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0x3C, 0x2B, 0x1A, 0xFF})

	text.Draw(screen, p.title, basicfont.Face7x13, 50, 50, color.RGBA{255, 255, 0, 255})

}

func (p *Level02Page) HandleInput(x, y int) {

	log.Printf("Level02Page received input at (%d, %d)\n", x, y)
}
