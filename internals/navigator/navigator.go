package navigator

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Page interface {
	Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error
	Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64)
	GetType() string
}

type Navigator struct {
	Stack      []Page
	Animating  bool
	Transition float64
	Direction  int
}

func NewNavigator() *Navigator {
	return &Navigator{
		Stack:      []Page{},
		Animating:  false,
		Transition: 1.0,
	}
}

func (n *Navigator) Push(nav Page) {
	if len(n.Stack) > 0 {
		n.Animating = true
		n.Transition = 0.0
		n.Direction = 1
	}
	n.Stack = append(n.Stack, nav)
}

func (n *Navigator) Pop() {
	if len(n.Stack) > 1 && !n.Animating {
		n.Animating = true
		n.Transition = 0.0
		n.Direction = -1
	}
}

func (n *Navigator) Update(navigatorOffsetX, navigatorOffsetY float32) (bool, error) {
	if len(n.Stack) == 0 {
		return false, nil
	}

	if n.Animating {
		n.Transition += 0.05
		if n.Transition >= 1.0 {
			n.Transition = 1.0
			n.Animating = false
			if n.Direction == -1 {

				n.Stack = n.Stack[:len(n.Stack)-1]
			}
		}
	}

	currentPage := n.Stack[len(n.Stack)-1]
	err := currentPage.Update(navigatorOffsetX, navigatorOffsetY, n.Animating)
	return n.Animating, err
}

func (n *Navigator) Draw(screen *ebiten.Image, navigatorAreaRect image.Rectangle) {

	navigatorArea := ebiten.NewImage(navigatorAreaRect.Dx(), navigatorAreaRect.Dy())
	navigatorArea.Fill(color.RGBA{30, 30, 30, 255})

	if n.Animating && len(n.Stack) > 1 {
		var prevOffsetX, currentOffsetX float64

		if n.Direction == 1 {

			prevOffsetX = -n.Transition * float64(navigatorAreaRect.Dx())

			currentOffsetX = float64(navigatorAreaRect.Dx()) * (1.0 - n.Transition)
		} else if n.Direction == -1 {

			currentOffsetX = float64(navigatorAreaRect.Dx()) * n.Transition

			prevOffsetX = -float64(navigatorAreaRect.Dx()) * (1.0 - n.Transition)
		}

		//fmt.Printf("PrevOffsetX: %f, CurrentOffsetX: %f\n", prevOffsetX, currentOffsetX)

		previousPage := n.Stack[len(n.Stack)-2]
		previousPage.Draw(navigatorArea, prevOffsetX, 0)

		currentPage := n.Stack[len(n.Stack)-1]
		currentPage.Draw(navigatorArea, currentOffsetX, 0)
	} else {

		if len(n.Stack) > 0 {
			currentPage := n.Stack[len(n.Stack)-1]
			currentPage.Draw(navigatorArea, 0, 0)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(navigatorAreaRect.Min.X), float64(navigatorAreaRect.Min.Y))
	screen.DrawImage(navigatorArea, op)
}

func (n *Navigator) CurrentPage() Page {
	if len(n.Stack) == 0 {
		return nil
	}
	return n.Stack[len(n.Stack)-1]
}
