package page

import (
	"image/color"

	"example.com/menu/internals/textwrapper"
	"example.com/menu/internals/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Page interface for all pages
type Page interface {
	Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error
	Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64)
	DrawBackground(screen *ebiten.Image)
	DrawUIelements(screen *ebiten.Image)
	AddUIelement(element UIElement)
	AddButton(btn PageButton, onNext func())
	GetType() string
	SetCustomDraw(func(screen *ebiten.Image)) // New method to set custom draw function
}

type PageButton struct {
	X, Y  float32
	Label string
}

type UIElement interface {
	Update(offsetX, offsetY float32, isAnimating bool)
	Draw(screen *ebiten.Image)
}

// BasePage struct to hold common page properties
type BasePage struct {
	X, Y                 float32
	Width, Height        float32
	BackgroundColor      color.Color
	Message              string
	UiElements           []UIElement
	TextWrapper          *textwrapper.TextWrapper
	NextPageID           string
	PageArea             *ebiten.Image              // Stores the page image
	DrawCustom           func(screen *ebiten.Image) // Custom drawing function
	DrawBackgroundCustom func(screen *ebiten.Image) // Custom draw background function
	DrawUIElementsCustom func(screen *ebiten.Image) // Custom draw elements function
}

// NewBasePage creates a new BasePage instance
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

// Update for BasePage
func (p *BasePage) Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error {
	for _, element := range p.UiElements {
		element.Update(navigatorOffsetX+p.X, navigatorOffsetY+p.Y, isAnimating)
	}
	return nil
}

// Draw method enhanced to include custom drawing
func (p *BasePage) Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64) {
	// Draw the background and elements onto PageArea
	p.DrawBackground(p.PageArea)
	p.DrawUIelements(p.PageArea)

	// Invoke custom draw function if set
	if p.DrawCustom != nil {
		p.DrawCustom(p.PageArea)
	}

	// Draw the composed PageArea onto the navigatorArea
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X)+offsetX, float64(p.Y)+offsetY)
	navigatorArea.DrawImage(p.PageArea, op)
}

// AddElement for BasePage
func (p *BasePage) AddUIelement(uiElement UIElement) {
	p.UiElements = append(p.UiElements, uiElement)
}

// DrawBackground for BasePage with optional custom draw
func (p *BasePage) DrawBackground(screen *ebiten.Image) {
	if p.DrawBackgroundCustom != nil {
		p.DrawBackgroundCustom(screen)
	} else {
		p.PageArea.Clear()
		p.PageArea.Fill(p.BackgroundColor)
	}
}

// DrawElements for BasePage with optional custom draw
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

// AddButton for BasePage
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

// SetCustomDraw sets the custom drawing function
func (p *BasePage) SetCustomDraw(drawFunc func(screen *ebiten.Image)) {
	p.DrawCustom = drawFunc
}

// SetCustomDrawBackground sets the custom drawing function for the background
func (p *BasePage) SetCustomDrawBackground(drawFunc func(screen *ebiten.Image)) {
	p.DrawBackgroundCustom = drawFunc
}

// SetCustomDrawElements sets the custom drawing function for the elements
func (p *BasePage) SetCustomDrawUIElements(drawFunc func(screen *ebiten.Image)) {
	p.DrawUIElementsCustom = drawFunc
}

// GetType implementation for BasePage
func (p *BasePage) GetType() string {
	return "BasePage"
}
