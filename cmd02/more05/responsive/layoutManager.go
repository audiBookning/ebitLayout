package responsive

import (
	"log"
	"sort"
	"sync"

	"example.com/menu/cmd02/more05/types"
)

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

func (lm *LayoutManager) CalculatePositions(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	lm.mutex.RLock()
	layoutMode := lm.currentMode
	lm.mutex.RUnlock()

	var positions map[string]types.Position

	switch layoutMode {
	case LayoutHorizontal:
		positions = lm.calcHorizontal(screenWidth, screenHeight, elementIDs, alignment, elements)
	case LayoutVertical:
		positions = lm.calcVertical(screenWidth, screenHeight, elementIDs, alignment, elements)
	case LayoutGrid:
		positions = lm.calcGrid(screenWidth, screenHeight, elementIDs, alignment, elements)
	default:

		positions = lm.calcHorizontal(screenWidth, screenHeight, elementIDs, alignment, elements)
	}

	return positions
}

func (lm *LayoutManager) calcHorizontal(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	spacing := 50

	totalWidth := 0
	for _, elem := range elements {
		width, _ := elem.GetSize()
		totalWidth += width
	}
	totalWidth += (numElements - 1) * spacing

	var startX int
	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth - totalWidth) / 2
	case AlignRight:
		startX = screenWidth - totalWidth
	}

	yPos := screenHeight - 50

	positions := make(map[string]types.Position)
	x := startX
	for i, elem := range elements {
		width, height := elem.GetSize()
		positions[elementIDs[i]] = types.Position{
			X:      x,
			Y:      yPos,
			Width:  width,
			Height: height,
		}
		x += width + spacing
	}

	return positions
}
