package responsive

import (
	"log"
	"sort"
	"sync"

	"example.com/menu/cmd02/more05/types"
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

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
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
		positions = calculateHorizontal(screenWidth, screenHeight, elementIDs, alignment, elements)
	case LayoutVertical:
		positions = calculateVertical(screenWidth, screenHeight, elementIDs, alignment, elements)
	case LayoutGrid:
		positions = calculateGrid(screenWidth, screenHeight, elementIDs, alignment, elements)
	default:

		positions = calculateHorizontal(screenWidth, screenHeight, elementIDs, alignment, elements)
	}

	return positions
}

func getElementSizes(elements []types.Element) []types.Position {
	sizes := make([]types.Position, len(elements))
	for i, elem := range elements {
		width, height := elem.GetSize()
		sizes[i] = types.Position{
			Width:  width,
			Height: height,
		}
	}
	return sizes
}

func calculateHorizontal(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
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

func calculateVertical(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	spacing := 20

	totalHeight := 0
	for _, elem := range elements {
		_, height := elem.GetSize()
		totalHeight += height
	}
	totalHeight += (numElements - 1) * spacing

	startY := (screenHeight - totalHeight) / 2

	var startX int
	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth) / 2
	case AlignRight:
		startX = screenWidth
	}

	positions := make(map[string]types.Position)
	y := startY
	for i, elem := range elements {
		width, height := elem.GetSize()
		switch alignment {
		case AlignLeft:

		case AlignCenter:
			startX = (screenWidth - width) / 2
		case AlignRight:
			startX = screenWidth - width
		}
		positions[elementIDs[i]] = types.Position{
			X:      startX,
			Y:      y,
			Width:  width,
			Height: height,
		}
		y += height + spacing
	}

	return positions
}

func calculateGrid(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	columns := 2
	rows := (numElements + columns - 1) / columns

	spacingX := 30
	spacingY := 30

	colWidths := make([]int, columns)
	rowHeights := make([]int, rows)

	for i, elem := range elements {
		col := i % columns
		row := i / columns
		width, height := elem.GetSize()
		if width > colWidths[col] {
			colWidths[col] = width
		}
		if height > rowHeights[row] {
			rowHeights[row] = height
		}
	}

	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w
	}
	totalWidth += (columns - 1) * spacingX

	totalHeight := 0
	for _, h := range rowHeights {
		totalHeight += h
	}
	totalHeight += (rows - 1) * spacingY

	var startX, startY int

	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth - totalWidth) / 2
	case AlignRight:
		startX = screenWidth - totalWidth
	}

	startY = (screenHeight - totalHeight) / 2

	positions := make(map[string]types.Position)

	colOffsets := make([]int, columns)
	currentX := startX
	for c := 0; c < columns; c++ {
		colOffsets[c] = currentX
		currentX += colWidths[c] + spacingX
	}

	rowOffsets := make([]int, rows)
	currentY := startY
	for r := 0; r < rows; r++ {
		rowOffsets[r] = currentY
		currentY += rowHeights[r] + spacingY
	}

	for i, elem := range elements {
		col := i % columns
		row := i / columns
		width, height := elem.GetSize()
		x := colOffsets[col]
		y := rowOffsets[row]
		positions[elementIDs[i]] = types.Position{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}
	}

	return positions
}

func getElementsByIds(ids []string) []types.Element {

	return nil
}
