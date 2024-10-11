package responsive

import (
	"log"
	"sort"
	"sync"

	"example.com/menu/cmd02/more04/types"
)

type LayoutMode string

const (
	LayoutHorizontal LayoutMode = "horizontal"
	LayoutVertical   LayoutMode = "vertical"
	LayoutGrid       LayoutMode = "grid"
)

type Breakpoint struct {
	Width      int
	LayoutMode LayoutMode
}

type LayoutManager struct {
	breakpoints []Breakpoint
	currentMode LayoutMode
	mutex       sync.RWMutex
}

func NewLayoutManager(breakpoints []Breakpoint) *LayoutManager {

	sort.Slice(breakpoints, func(i, j int) bool {
		return breakpoints[i].Width > breakpoints[j].Width
	})

	var initialMode LayoutMode
	if len(breakpoints) > 0 {
		initialMode = breakpoints[len(breakpoints)-1].LayoutMode
	}

	return &LayoutManager{
		breakpoints: breakpoints,
		currentMode: initialMode,
	}
}

func (lm *LayoutManager) DetermineLayout(screenWidth int) LayoutMode {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	for _, bp := range lm.breakpoints {
		if screenWidth >= bp.Width {
			if bp.LayoutMode != lm.currentMode {
				lm.currentMode = bp.LayoutMode
				log.Printf("DetermineLayout: screenWidth=%d, using LayoutMode=%s\n", screenWidth, lm.currentMode)
			}
			return lm.currentMode
		}
	}

	if len(lm.breakpoints) > 0 {
		lm.currentMode = lm.breakpoints[len(lm.breakpoints)-1].LayoutMode
		log.Printf("DetermineLayout: screenWidth=%d, using default LayoutMode=%s\n", screenWidth, lm.currentMode)
	}
	return lm.currentMode
}

func (lm *LayoutManager) GetCurrentLayoutMode() LayoutMode {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.currentMode
}

func (lm *LayoutManager) CalculatePositions(screenWidth, screenHeight int, elements []string) map[string]types.Position {
	lm.mutex.RLock()
	layoutMode := lm.currentMode
	lm.mutex.RUnlock()

	positions := make(map[string]types.Position)

	switch layoutMode {
	case LayoutHorizontal:
		positions = calculateHorizontal(screenWidth, screenHeight, elements)
	case LayoutVertical:
		positions = calculateVertical(screenWidth, screenHeight, elements)
	case LayoutGrid:
		positions = calculateGrid(screenWidth, screenHeight, elements)
	default:

		positions = calculateHorizontal(screenWidth, screenHeight, elements)
	}

	return positions
}

func calculateHorizontal(screenWidth, screenHeight int, elements []string) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	buttonWidth := 200
	buttonHeight := 50
	spacing := 50

	totalWidth := numElements*buttonWidth + (numElements-1)*spacing
	startX := (screenWidth - totalWidth) / 2
	yPos := screenHeight - buttonHeight - 50

	positions := make(map[string]types.Position)
	for i, elem := range elements {
		x := startX + i*(buttonWidth+spacing)
		positions[elem] = types.Position{
			X:      x,
			Y:      yPos,
			Width:  buttonWidth,
			Height: buttonHeight,
		}
	}

	return positions
}

func calculateVertical(screenWidth, screenHeight int, elements []string) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	buttonWidth := 150
	buttonHeight := 40
	spacing := 20

	totalHeight := numElements*buttonHeight + (numElements-1)*spacing
	startY := (screenHeight - totalHeight) / 2

	positions := make(map[string]types.Position)
	for i, elem := range elements {
		x := (screenWidth - buttonWidth) / 2
		y := startY + i*(buttonHeight+spacing)
		positions[elem] = types.Position{
			X:      x,
			Y:      y,
			Width:  buttonWidth,
			Height: buttonHeight,
		}
	}

	return positions
}

func calculateGrid(screenWidth, screenHeight int, elements []string) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	columns := 2
	rows := (numElements + 1) / 2

	buttonWidth := 180
	buttonHeight := 45
	spacingX := 30
	spacingY := 30

	totalWidth := columns*buttonWidth + (columns-1)*spacingX
	totalHeight := rows*buttonHeight + (rows-1)*spacingY

	startX := (screenWidth - totalWidth) / 2
	startY := (screenHeight - totalHeight) / 2

	positions := make(map[string]types.Position)
	for i, elem := range elements {
		row := i / columns
		col := i % columns
		x := startX + col*(buttonWidth+spacingX)
		y := startY + row*(buttonHeight+spacingY)
		positions[elem] = types.Position{
			X:      x,
			Y:      y,
			Width:  buttonWidth,
			Height: buttonHeight,
		}
	}

	return positions
}
