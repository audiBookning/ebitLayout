package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ***** INPUT MANAGER *****
type InputManager struct {
	Clickables     []Clickable
	MouseX, MouseY int
}

func (im *InputManager) Register(c Clickable) {
	im.Clickables = append(im.Clickables, c)
}

func (im *InputManager) Update() {
	im.MouseX, im.MouseY = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, c := range im.Clickables {
			if c.Contains(im.MouseX, im.MouseY) {
				c.OnMouseDown()
				break
			}
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, c := range im.Clickables {
			if c.Contains(im.MouseX, im.MouseY) {
				c.OnClick()
				break
			}
		}
	}

	// Update hover state
	for _, c := range im.Clickables {
		c.SetHovered(c.Contains(im.MouseX, im.MouseY))
	}
}
