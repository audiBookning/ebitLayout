package responsive

import (
	"log"
	"sort"
	"sync"
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

	return &LayoutManager{
		breakpoints: breakpoints,
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

func (lm *LayoutManager) CalculatePositions(screenWidth, screenHeight int, elements []string) map[string]Position {
	layoutMode := lm.DetermineLayout(screenWidth)
	positions := make(map[string]Position)

	switch layoutMode {
	case LayoutHorizontal:
		positions = calculateHorizontal(screenWidth, screenHeight, elements)
	case LayoutVertical:
		positions = calculateVertical(screenWidth, screenHeight, elements)
	case LayoutGrid:
		positions = calculateGrid(screenWidth, screenHeight, elements)
	default:
		positions = calculateVertical(screenWidth, screenHeight, elements)
	}

	return positions
}
