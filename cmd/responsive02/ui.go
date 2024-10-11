package main

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Breakpoint struct {
	Width      int
	LayoutMode string
}

type Button struct {
	Text          string
	X, Y          int
	Width, Height int
	Clicked       bool
	OnClickFunc   func()
}

func (b *Button) IsClicked(x, y int) bool {
	return x >= b.X && x <= b.X+b.Width &&
		y >= b.Y && y <= b.Y+b.Height
}

func (b *Button) OnClick() {
	if b.OnClickFunc != nil {
		b.Clicked = true
		b.OnClickFunc()

		b.Clicked = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {

	var btnColor color.Color
	if b.Clicked {
		btnColor = color.RGBA{0xFF, 0xA5, 0x00, 0xFF}
	} else {
		btnColor = color.RGBA{0x00, 0x7A, 0xCC, 0xFF}
	}

	button := ebiten.NewImage(b.Width, b.Height)
	button.Fill(btnColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	screen.DrawImage(button, op)

	textBounds := text.BoundString(basicfont.Face7x13, b.Text)
	textX := b.X + (b.Width-textBounds.Dx())/2
	textY := b.Y + (b.Height+textBounds.Dy())/2
	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, color.White)
}

type UI struct {
	Title          string
	TitleX, TitleY int
	Button1        Button
	Button2        Button
	Breakpoints    []Breakpoint
	CurrentMode    string
}

func NewUI(breakpoints []Breakpoint) *UI {

	sort.Slice(breakpoints, func(i, j int) bool {
		return breakpoints[i].Width > breakpoints[j].Width
	})

	ui := &UI{
		Title:       "Responsive Layout",
		Breakpoints: breakpoints,
	}

	ui.Button1 = Button{
		Text: "Button 1",
		OnClickFunc: func() {
			fmt.Println("Button 1 clicked!")
		},
	}
	ui.Button2 = Button{
		Text: "Button 2",
		OnClickFunc: func() {
			fmt.Println("Button 2 clicked!")
		},
	}

	return ui
}

func (u *UI) Update(screenWidth, screenHeight int) {

	titleBounds := text.BoundString(basicfont.Face7x13, u.Title)
	u.TitleX = (screenWidth - titleBounds.Dx()) / 2
	u.TitleY = 50

	u.determineLayoutMode(screenWidth)

	switch u.CurrentMode {
	case "horizontal":
		u.layoutHorizontal(screenWidth, screenHeight)
	case "grid":
		u.layoutGrid(screenWidth, screenHeight)
	case "vertical":
		u.layoutVertical(screenWidth, screenHeight)
	default:
		u.layoutVertical(screenWidth, screenHeight)
	}
}

func (u *UI) determineLayoutMode(screenWidth int) {
	for _, bp := range u.Breakpoints {
		if screenWidth <= bp.Width {
			u.CurrentMode = bp.LayoutMode
		} else {
			continue
		}
	}

	if u.CurrentMode == "" && len(u.Breakpoints) > 0 {
		u.CurrentMode = u.Breakpoints[0].LayoutMode
	}
}

func (u *UI) layoutHorizontal(screenWidth, screenHeight int) {

	btnWidth := 200
	btnHeight := 50

	totalWidth := 2*btnWidth + 50
	startX := (screenWidth - totalWidth) / 2
	yPos := screenHeight - btnHeight - 50

	u.Button1.X = startX
	u.Button1.Y = yPos
	u.Button1.Width = btnWidth
	u.Button1.Height = btnHeight

	u.Button2.X = startX + btnWidth + 50
	u.Button2.Y = yPos
	u.Button2.Width = btnWidth
	u.Button2.Height = btnHeight
}

func (u *UI) layoutVertical(screenWidth, screenHeight int) {

	btnWidth := 150
	btnHeight := 40

	totalHeight := 2*btnHeight + 20
	startY := screenHeight - totalHeight - 50

	u.Button1.X = (screenWidth - btnWidth) / 2
	u.Button1.Y = startY
	u.Button1.Width = btnWidth
	u.Button1.Height = btnHeight

	u.Button2.X = (screenWidth - btnWidth) / 2
	u.Button2.Y = startY + btnHeight + 20
	u.Button2.Width = btnWidth
	u.Button2.Height = btnHeight
}

func (u *UI) layoutGrid(screenWidth, screenHeight int) {

	btnWidth := 180
	btnHeight := 45

	totalWidth := 2*btnWidth + 30
	startX := (screenWidth - totalWidth) / 2
	yPos := screenHeight - btnHeight - 60

	u.Button1.X = startX
	u.Button1.Y = yPos
	u.Button1.Width = btnWidth
	u.Button1.Height = btnHeight

	u.Button2.X = startX + btnWidth + 30
	u.Button2.Y = yPos
	u.Button2.Width = btnWidth
	u.Button2.Height = btnHeight
}

func (u *UI) HandleClick(x, y int) {
	if u.Button1.IsClicked(x, y) {
		u.Button1.OnClick()
	}
	if u.Button2.IsClicked(x, y) {
		u.Button2.OnClick()
	}
}

func (u *UI) Draw(screen *ebiten.Image) {

	text.Draw(screen, u.Title, basicfont.Face7x13, u.TitleX, u.TitleY, color.White)

	u.Button1.Draw(screen)
	u.Button2.Draw(screen)
}
