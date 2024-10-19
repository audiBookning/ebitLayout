package responsive

import (
	"fmt"
	"sync"

	"example.com/menu/cmd02/more06/textwrapper"
	"example.com/menu/cmd02/more06/types"
	"github.com/hajimehoshi/ebiten/v2"
)

type UI struct {
	Title       *Title
	Fields      []types.Element
	manager     *LayoutManager
	elementsIds []string
	mutex       sync.RWMutex
	TextWrapper *textwrapper.TextWrapper
	Alignment   Alignment
}

func NewUI(titleText string, breakpoints []Breakpoint, fields []types.Element, tw *textwrapper.TextWrapper, alignment Alignment) *UI {
	elements := make([]string, len(fields))
	for i := range fields {
		elements[i] = fmt.Sprintf("field%d", i+1)
	}

	return &UI{
		Title:       NewTitle(titleText, tw),
		Fields:      fields,
		manager:     NewLayoutManager(breakpoints),
		elementsIds: elements,
		TextWrapper: tw,
		Alignment:   alignment,
	}
}

func (u *UI) LayoutUpdate(screenWidth, screenHeight int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.manager.DetermineLayout(screenWidth)

	positions := u.manager.CalculatePositions(screenWidth, screenHeight, u.elementsIds, u.Alignment, u.Fields)

	for i, field := range u.Fields {
		elemID := u.elementsIds[i]
		pos, exists := positions[elemID]
		if exists {
			field.SetPosition(pos)
		} else {
			field.SetPosition(types.Position{X: 0, Y: i * 50, Width: 100, Height: 40})
		}
		field.Update()
	}

	titleWidth, _ := u.TextWrapper.MeasureText(u.Title.Text)
	u.Title.X = (screenWidth - int(titleWidth)) / 2
	u.Title.Y = 50
}

func (u *UI) HandleClick(x, y int) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	for _, field := range u.Fields {
		if field.IsClicked(x, y) {
			field.HandleClick()
		}
	}
}

func (u *UI) Draw(screen *ebiten.Image) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	u.Title.Draw(screen)

	for _, field := range u.Fields {
		field.Draw(screen)
	}
}

func (u *UI) ResetFieldStates() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, field := range u.Fields {
		field.ResetState()
	}
}
