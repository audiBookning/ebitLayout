package widgets

import (
	"fmt"
	"sync"

	"example.com/menu/cmd/responsive05/responsive"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type UI struct {
	Title    *Title
	Buttons  []*Button
	manager  *responsive.LayoutManager
	elements []string
	mutex    sync.RWMutex
}

func NewUI(titleText string, breakpoints []responsive.Breakpoint, buttons []*Button) *UI {
	elements := make([]string, len(buttons))
	for i := range buttons {
		elements[i] = fmt.Sprintf("button%d", i+1)
	}

	return &UI{
		Title:    NewTitle(titleText),
		Buttons:  buttons,
		manager:  responsive.NewLayoutManager(breakpoints),
		elements: elements,
	}
}

func (u *UI) Update(screenWidth, screenHeight int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	titleBounds := text.BoundString(basicfont.Face7x13, u.Title.Text)
	u.Title.X = (screenWidth - titleBounds.Dx()) / 2
	u.Title.Y = 50
	//log.Printf("UI.Update: Title position=(%d, %d)\n", u.Title.X, u.Title.Y)

	positions := u.manager.CalculatePositions(screenWidth, screenHeight, u.elements)

	for i, btn := range u.Buttons {
		pos, exists := positions[u.elements[i]]
		if exists {
			btn.Position = pos
			//log.Printf("UI.Update: Button '%s' position=(%d, %d)\n", btn.Text, btn.Position.X, btn.Position.Y)
		}

		btn.Update()
	}
}

func (u *UI) HandleClick(x, y int) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	for _, btn := range u.Buttons {
		if btn.IsClicked(x, y) {
			btn.HandleClick()
		}
	}
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	u.Title.Draw(screen)

	for _, btn := range u.Buttons {
		btn.Draw(screen)
	}
}
