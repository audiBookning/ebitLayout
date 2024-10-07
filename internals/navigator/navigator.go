package navigator

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Page interface for all pages
type Page interface {
	Update(navigatorOffsetX, navigatorOffsetY float32, isAnimating bool) error
	Draw(navigatorArea *ebiten.Image, offsetX, offsetY float64)
	GetType() string
}

// Navigator manages a stack of navigable objects with animation capabilities.
type Navigator struct {
	Stack      []Page
	Animating  bool
	Transition float64 // Animation progress (0 to 1)
	Direction  int     // 1 for push, -1 for pop
}

// NewNavigator initializes a new Navigator instance with animation support.
func NewNavigator() *Navigator {
	return &Navigator{
		Stack:      []Page{},
		Animating:  false,
		Transition: 1.0, // Start with no transition
	}
}

// Push pushes a new navigable object onto the stack and starts the push animation.
func (n *Navigator) Push(nav Page) {
	if len(n.Stack) > 0 {
		n.Animating = true
		n.Transition = 0.0
		n.Direction = 1 // Push direction
	}
	n.Stack = append(n.Stack, nav)
}

// Pop pops the top navigable object from the stack and starts the pop animation.
func (n *Navigator) Pop() {
	if len(n.Stack) > 1 && !n.Animating {
		n.Animating = true
		n.Transition = 0.0
		n.Direction = -1 // Pop direction
	}
}

// Update updates the navigator's state and handles animations.
func (n *Navigator) Update(navigatorOffsetX, navigatorOffsetY float32) (bool, error) {
	if len(n.Stack) == 0 {
		return false, nil
	}

	// Handle animation transitions
	if n.Animating {
		n.Transition += 0.05 // Adjust this value for animation speed
		if n.Transition >= 1.0 {
			n.Transition = 1.0
			n.Animating = false
			if n.Direction == -1 {
				// Complete the pop after animation
				n.Stack = n.Stack[:len(n.Stack)-1]
			}
		}
	}

	// Update the current page with navigator offsets and animation state
	currentPage := n.Stack[len(n.Stack)-1]
	err := currentPage.Update(navigatorOffsetX, navigatorOffsetY, n.Animating)
	return n.Animating, err
}

// Draw renders the navigator area and manages page animations.
func (n *Navigator) Draw(screen *ebiten.Image, navigatorAreaRect image.Rectangle) {
	// Create an off-screen image for the navigator area
	navigatorArea := ebiten.NewImage(navigatorAreaRect.Dx(), navigatorAreaRect.Dy())
	navigatorArea.Fill(color.RGBA{30, 30, 30, 255}) // Optional: Background color for navigator area

	if n.Animating && len(n.Stack) > 1 {
		var prevOffsetX, currentOffsetX float64

		if n.Direction == 1 { // Push
			// Current Page slides to the left
			prevOffsetX = -n.Transition * float64(navigatorAreaRect.Dx())
			// New Page slides in from the right
			currentOffsetX = float64(navigatorAreaRect.Dx()) * (1.0 - n.Transition)
		} else if n.Direction == -1 { // Pop
			// Current Page slides to the right
			currentOffsetX = float64(navigatorAreaRect.Dx()) * n.Transition
			// Previous Page slides in from the left
			prevOffsetX = -float64(navigatorAreaRect.Dx()) * (1.0 - n.Transition)
		}

		// Debugging: Log animation state
		//fmt.Printf("PrevOffsetX: %f, CurrentOffsetX: %f\n", prevOffsetX, currentOffsetX)

		// Draw the Previous Page with its own Y position
		previousPage := n.Stack[len(n.Stack)-2]
		previousPage.Draw(navigatorArea, prevOffsetX, 0)

		// Draw the Current (New) Page with its own Y position
		currentPage := n.Stack[len(n.Stack)-1]
		currentPage.Draw(navigatorArea, currentOffsetX, 0)
	} else {
		// No animation, draw the top page normally
		if len(n.Stack) > 0 {
			currentPage := n.Stack[len(n.Stack)-1]
			currentPage.Draw(navigatorArea, 0, 0)
		}
	}

	// Draw the navigator area onto the main screen within the defined rectangle
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(navigatorAreaRect.Min.X), float64(navigatorAreaRect.Min.Y))
	screen.DrawImage(navigatorArea, op)
}

// CurrentPage retrieves the current active page.
func (n *Navigator) CurrentPage() Page {
	if len(n.Stack) == 0 {
		return nil
	}
	return n.Stack[len(n.Stack)-1]
}
