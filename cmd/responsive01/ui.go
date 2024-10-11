package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type UI struct {
	Title         string
	ButtonText    string
	ButtonClicked bool

	TitleX, TitleY            int
	ButtonX, ButtonY          int
	ButtonWidth, ButtonHeight int
}

func NewUI() *UI {
	return &UI{
		Title:      "Responsive Layout",
		ButtonText: "Click Me",
	}
}

func (u *UI) Update(screenWidth, screenHeight int) {

	titleBounds := text.BoundString(basicfont.Face7x13, u.Title)
	u.TitleX = (screenWidth - titleBounds.Dx()) / 2
	u.TitleY = 50

	u.ButtonWidth = int(200 * float64(screenWidth) / 800)
	u.ButtonHeight = int(50 * float64(screenHeight) / 600)
	u.ButtonX = (screenWidth - u.ButtonWidth) / 2
	u.ButtonY = screenHeight - u.ButtonHeight - 50
}

func (u *UI) Draw(screen *ebiten.Image) {

	text.Draw(screen, u.Title, basicfont.Face7x13, u.TitleX, u.TitleY, color.White)

	button := ebiten.NewImage(u.ButtonWidth, u.ButtonHeight)
	if u.ButtonClicked {
		button.Fill(color.RGBA{0xFF, 0x45, 0x00, 0xFF})
	} else {
		button.Fill(color.RGBA{0x00, 0x7A, 0xCC, 0xFF})
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(u.ButtonX), float64(u.ButtonY))
	screen.DrawImage(button, op)

	textBounds := text.BoundString(basicfont.Face7x13, u.ButtonText)
	textX := u.ButtonX + (u.ButtonWidth-textBounds.Dx())/2
	textY := u.ButtonY + (u.ButtonHeight+textBounds.Dy())/2 + 4
	text.Draw(screen, u.ButtonText, basicfont.Face7x13, textX, textY, color.White)
}

func (u *UI) IsButtonClicked(x, y int) bool {
	return x >= u.ButtonX && x <= u.ButtonX+u.ButtonWidth &&
		y >= u.ButtonY && y <= u.ButtonY+u.ButtonHeight
}

func (u *UI) OnButtonClick() {
	u.ButtonClicked = true
	fmt.Println("Button clicked!")

	go func() {

		for i := 0; i < 100000000; i++ {
		}
		u.ButtonClicked = false
	}()
}
