package main

import (
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Accordions []widgets.Accordion
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
	accordion1 := widgets.Accordion{
		Sections: []widgets.AccordionSection{
			{Title: "Section 1", Content: "Content 1", Collapsed: true},
			{Title: "Section 2", Content: "Content 2", Collapsed: true},
		},
		Width: 100,
		X:     20,
		Y:     20,
	}

	accordion2 := widgets.Accordion{
		Sections: []widgets.AccordionSection{
			{Title: "Section 3", Content: "Content 3", Collapsed: true},
			{Title: "Section 4", Content: "Content 4", Collapsed: true},
		},
		Width: 100,
		X:     150,
		Y:     20,
	}

	game := &Game{Accordions: []widgets.Accordion{accordion1, accordion2}}
	ebiten.SetWindowSize(680, 480)
	ebiten.SetWindowTitle("Accordion UI")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
