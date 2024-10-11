package responsive

import (
	"log"
	"sort"
	"sync"

	"example.com/menu/cmd02/more05/types"
)

// LayoutMode defines different layout strategies.
type LayoutMode string

const (
	LayoutHorizontal LayoutMode = "horizontal"
	LayoutVertical   LayoutMode = "vertical"
	LayoutGrid       LayoutMode = "grid"
)

// Breakpoint defines a screen width and the corresponding layout mode.
type Breakpoint struct {
	Width      int
	LayoutMode LayoutMode
}

// Alignment defines the horizontal alignment of UI elements.
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// LayoutManager manages responsive layouts based on breakpoints.
type LayoutManager struct {
	breakpoints []Breakpoint
	currentMode LayoutMode
	mutex       sync.RWMutex
}

// NewLayoutManager initializes a LayoutManager with given breakpoints.
// Breakpoints should be sorted in descending order of Width.
// If not sorted, they will be sorted automatically.
func NewLayoutManager(breakpoints []Breakpoint) *LayoutManager {
	// Sort breakpoints in descending order of Width
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

// DetermineLayout determines the current layout mode based on the screen width.
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
	// If no breakpoint matched, use the smallest layout mode
	if len(lm.breakpoints) > 0 {
		lm.currentMode = lm.breakpoints[len(lm.breakpoints)-1].LayoutMode
		log.Printf("DetermineLayout: screenWidth=%d, using default LayoutMode=%s\n", screenWidth, lm.currentMode)
	}
	return lm.currentMode
}

// GetCurrentLayoutMode returns the current layout mode.
func (lm *LayoutManager) GetCurrentLayoutMode() LayoutMode {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.currentMode
}

// CalculatePositions calculates the positions of UI elements based on the current layout and alignment.
func (lm *LayoutManager) CalculatePositions(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	lm.mutex.RLock()
	layoutMode := lm.currentMode
	lm.mutex.RUnlock()

	var positions map[string]types.Position

	// Calculate positions based on layoutMode
	switch layoutMode {
	case LayoutHorizontal:
		positions = calculateHorizontal(screenWidth, screenHeight, elementIDs, alignment, elements)
	case LayoutVertical:
		positions = calculateVertical(screenWidth, screenHeight, elementIDs, alignment, elements)
	case LayoutGrid:
		positions = calculateGrid(screenWidth, screenHeight, elementIDs, alignment, elements)
	default:
		// Fallback to horizontal layout if unknown layoutMode
		positions = calculateHorizontal(screenWidth, screenHeight, elementIDs, alignment, elements)
	}

	return positions
}

// Helper function to get element sizes
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

// Calculate positions for horizontal layout
func calculateHorizontal(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	spacing := 50 // You can also make spacing dynamic if needed

	// Calculate total width dynamically
	totalWidth := 0
	for _, elem := range elements {
		width, _ := elem.GetSize()
		totalWidth += width
	}
	totalWidth += (numElements - 1) * spacing

	// Determine startX based on alignment
	var startX int
	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth - totalWidth) / 2
	case AlignRight:
		startX = screenWidth - totalWidth
	}

	yPos := screenHeight - 50 // Adjust as needed

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

// Calculate positions for vertical layout
func calculateVertical(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	spacing := 20 // You can also make spacing dynamic if needed

	// Calculate total height dynamically
	totalHeight := 0
	for _, elem := range elements {
		_, height := elem.GetSize()
		totalHeight += height
	}
	totalHeight += (numElements - 1) * spacing

	startY := (screenHeight - totalHeight) / 2

	// Determine startX based on alignment
	var startX int
	switch alignment {
	case AlignLeft:
		startX = 0
	case AlignCenter:
		startX = (screenWidth) / 2 // Will adjust each element to be centered
	case AlignRight:
		startX = screenWidth // Will adjust each element to align to the right
	}

	positions := make(map[string]types.Position)
	y := startY
	for i, elem := range elements {
		width, height := elem.GetSize()
		switch alignment {
		case AlignLeft:
			// startX is already set to 0
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

// Calculate positions for grid layout
func calculateGrid(screenWidth, screenHeight int, elementIDs []string, alignment Alignment, elements []types.Element) map[string]types.Position {
	numElements := len(elements)
	if numElements == 0 {
		return nil
	}

	columns := 2
	rows := (numElements + columns - 1) / columns // Ceiling division

	spacingX := 30
	spacingY := 30

	// Calculate total grid width and height dynamically
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

	// Determine startX and startY based on alignment
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

// Placeholder function to retrieve elements by their IDs.
// You need to implement this based on your application's context.
func getElementsByIds(ids []string) []types.Element {
	// Example implementation. Replace with actual retrieval logic.
	// This might involve accessing the UI state or passing a reference.
	return nil
}
