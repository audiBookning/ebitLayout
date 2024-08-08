package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AccordionSection struct {
	Title     string
	Content   string
	Collapsed bool
}

type Accordion struct {
	X, Y     int
	Sections []AccordionSection
	Width    int
}

func (a *Accordion) Update() {
	x, y := ebiten.CursorPosition()
	mouseClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	yOffset := a.Y
	for i := range a.Sections {
		headerY := yOffset

		// Check if the mouse is within the header of this section
		if mouseClicked && y >= headerY && y < headerY+30 && x >= a.X && x < a.X+a.Width {
			a.Sections[i].Collapsed = !a.Sections[i].Collapsed
		}

		// Update the yOffset for the next section
		yOffset += 30
		if !a.Sections[i].Collapsed {
			yOffset += 50
		}
	}
}

func (a *Accordion) Draw(screen *ebiten.Image) {
	yOffset := a.Y
	for _, section := range a.Sections {
		// Draw header
		vector.DrawFilledRect(screen, float32(a.X), float32(yOffset), float32(a.Width), 30, color.RGBA{100, 100, 100, 255}, true)
		ebitenutil.DebugPrintAt(screen, section.Title, a.X+10, yOffset+10)

		yOffset += 30

		// Draw content if not collapsed
		if !section.Collapsed {
			vector.DrawFilledRect(screen, float32(a.X), float32(yOffset), float32(a.Width), 50, color.RGBA{100, 50, 50, 255}, true)
			ebitenutil.DebugPrintAt(screen, section.Content, a.X+10, yOffset+10)
			yOffset += 50
		}
	}
}

type Game struct {
	Accordions []Accordion
}

func (g *Game) Update() error {
	for i := range g.Accordions {
		g.Accordions[i].Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, accordion := range g.Accordions {
		accordion.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	accordion1 := Accordion{
		Sections: []AccordionSection{
			{Title: "Section 1", Content: "Content 1", Collapsed: true},
			{Title: "Section 2", Content: "Content 2", Collapsed: true},
		},
		Width: 100,
		X:     20,
		Y:     20,
	}

	accordion2 := Accordion{
		Sections: []AccordionSection{
			{Title: "Section 3", Content: "Content 3", Collapsed: true},
			{Title: "Section 4", Content: "Content 4", Collapsed: true},
		},
		Width: 100,
		X:     150,
		Y:     20,
	}

	game := &Game{Accordions: []Accordion{accordion1, accordion2}}
	ebiten.SetWindowSize(680, 480)
	ebiten.SetWindowTitle("Accordion UI")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
