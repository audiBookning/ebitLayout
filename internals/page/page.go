package page

import (
	"image/color"

	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Page interface {
	Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error
	Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64)
	DrawBackground(screen *ebiten.Image)
	DrawUIelements(screen *ebiten.Image)
	AddUIelement(element UIElement)
	AddButton(btn PageButton, onNext func())
	GetType() string
	SetCustomDraw(func(screen *ebiten.Image))
}

type PageButton struct {
	X, Y  float32
	Label string
}

type UIElement interface {
	Update(offsetX, offsetY float32, isAnimating bool)
	Draw(screen *ebiten.Image)
}

type BasePage struct {
	X, Y                 float32
	Width, Height        float32
	BackgroundColor      color.Color
	Message              string
	UiElements           []UIElement
	TextWrapper          *textwrapper.TextWrapper
	NextPageID           string
	PageArea             *ebiten.Image
	DrawCustom           func(screen *ebiten.Image)
	DrawBackgroundCustom func(screen *ebiten.Image)
	DrawUIElementsCustom func(screen *ebiten.Image)
}

func NewBasePage(
	bgColor color.Color,
	message string,
	tw *textwrapper.TextWrapper,
	x, y, width, height float32,
) *BasePage {
	return &BasePage{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		BackgroundColor: bgColor,
		Message:         message,
		UiElements:      make([]UIElement, 0),
		TextWrapper:     tw,
		NextPageID:      "",
		PageArea:        ebiten.NewImage(int(width), int(height)),
	}
}

func (p *BasePage) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error {
	for _, element := range p.UiElements {
		element.Update(navigatorOffsetX+p.X, navigatorOffsetY+p.Y, isAnimating)
	}
	return nil
}

func (p *BasePage) Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64) {

	p.DrawBackground(p.PageArea)
	p.DrawUIelements(p.PageArea)

	if p.DrawCustom != nil {
		p.DrawCustom(p.PageArea)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X)+offsetX, float64(p.Y)+offsetY)
	navigatorArea.DrawImage(p.PageArea, op)
}

func (p *BasePage) AddUIelement(uiElement UIElement) {
	p.UiElements = append(p.UiElements, uiElement)
}

func (p *BasePage) DrawBackground(screen *ebiten.Image) {
	if p.DrawBackgroundCustom != nil {
		p.DrawBackgroundCustom(screen)
	} else {
		p.PageArea.Clear()
		p.PageArea.Fill(p.BackgroundColor)
	}
}

func (p *BasePage) DrawUIelements(screen *ebiten.Image) {
	if p.DrawUIElementsCustom != nil {
		p.DrawUIElementsCustom(screen)
	} else {
		ebitenutil.DebugPrintAt(screen, p.Message, 10, 10)
		for _, uiElement := range p.UiElements {
			uiElement.Draw(screen)
		}
	}
}

func (p *BasePage) AddButton(btn PageButton, onNext func()) {
	button := widgets.NewButtonStd(
		btn.X,
		btn.Y,
		100,
		40,
		btn.Label,
		p.TextWrapper,
		color.RGBA{0, 128, 255, 255},
		color.White,
		16,
		onNext,
	)
	p.AddUIelement(button)
}

func (p *BasePage) SetCustomDraw(drawFunc func(screen *ebiten.Image)) {
	p.DrawCustom = drawFunc
}

func (p *BasePage) SetCustomDrawBackground(drawFunc func(screen *ebiten.Image)) {
	p.DrawBackgroundCustom = drawFunc
}

func (p *BasePage) SetCustomDrawUIElements(drawFunc func(screen *ebiten.Image)) {
	p.DrawUIElementsCustom = drawFunc
}

func (p *BasePage) GetType() string {
	return "BasePage"
}
