package responsive

import (
	"fmt"
	"sync"

	"example.com/menu/cmd02/more05/textwrapper"
	"example.com/menu/cmd02/more05/types"
	"github.com/hajimehoshi/ebiten/v2"
)

// UI represents a page's UI, containing a title and fields.
type UI struct {
	Title       *Title
	Fields      []types.Element
	manager     *LayoutManager
	elementsIds []string // Identifiers for layout positioning
	mutex       sync.RWMutex
	TextWrapper *textwrapper.TextWrapper
	Alignment   Alignment // Alignment of the UI elements
}

// NewUI creates a new UI instance with the specified alignment.
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

// Update recalculates and sets the positions of all UI elements based on the current screen size.
func (u *UI) Update(screenWidth, screenHeight int) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.manager.DetermineLayout(screenWidth)

	// Pass actual elements to LayoutManager
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

	// Update Title Position (centered at the top)
	titleWidth, _ := u.TextWrapper.MeasureText(u.Title.Text)
	u.Title.X = (screenWidth - int(titleWidth)) / 2
	u.Title.Y = 50
}

// HandleClick processes click events for all fields.
func (u *UI) HandleClick(x, y int) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	for _, field := range u.Fields {
		if field.IsClicked(x, y) {
			field.HandleClick()
		}
	}
}

// Draw renders the UI on the screen.
func (u *UI) Draw(screen *ebiten.Image) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	u.Title.Draw(screen)

	for _, field := range u.Fields {
		field.Draw(screen)
	}
}

// ResetFieldStates resets the state of all fields.
func (u *UI) ResetFieldStates() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, field := range u.Fields {
		field.ResetState()
	}
}
